package app

import (
	"context"
	"crawler/app/model"
	"crawler/communication"
	"crawler/config"
	"crawler/ds"
	"crawler/pipeline"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/pressly/goose"

	"github.com/jmoiron/sqlx"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/stdlib"
)

type App struct {
	settings config.Settings

	urlsReader   communication.KafkaTopicReader[model.CrawlTaskKey, model.CrawlTaskValue]
	urlsWriter   communication.KafkaTopicWriter[model.CrawlTaskKey, model.CrawlTaskValue]
	docsWriter   communication.KafkaTopicWriter[model.DocKey, model.DocValue]
	indexStorage communication.IndexMetadataStorage

	pipeline pipeline.CrawlPipe
}

func NewApp(s config.Settings) *App {
	return &App{
		settings:     s,
		urlsReader:   communication.NewKafkaReader[model.CrawlTaskKey, model.CrawlTaskValue](s.KafkaBrokersList, s.KafkaUrlsTopic, s.KafkaConsumerGroupId, log.Default()),
		urlsWriter:   communication.NewKafkaWriter[model.CrawlTaskKey, model.CrawlTaskValue](s.KafkaBrokersList, s.KafkaUrlsTopic, log.Default()),
		docsWriter:   communication.NewKafkaWriter[model.DocKey, model.DocValue](s.KafkaBrokersList, s.KafkaResultsTopic, log.Default()),
		indexStorage: communication.NewIndexMetadataStorage(initPostgres(s.PostgresConnectionString)),
		pipeline:     pipeline.NewPipeline(),
	}
}

func (a *App) RunCrawlPipeline(ctx context.Context) {
	tasksChan := a.urlsReader.ReadJSON(ctx)
	pipesCtx := pipeline.NewContext(ctx, a.urlsWriter, a.indexStorage)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case t, ok := <-tasksChan:
				if !ok {
					log.Printf("INFO: tasks chan has been closed, crawl pipeline shutdown")
					return
				}

				log.Printf("INFO: new task received: %+v", t)

				result, err := a.pipeline.Do(
					pipesCtx,
					model.CrawlTask{
						PathToDiscover: t.Value.Path,
					},
				)
				if err != nil {
					log.Printf("ERROR: cannot do pipeline job: %s", err)
					continue
				}

				if result.Status != model.SuccessCrawlStatus {
					continue
				}

				docKey := model.DocKey{
					FileName: generateFileName(t.Key.Hostname),
				}

				docKeyBytes, err := json.Marshal(docKey)
				if err != nil {
					log.Printf("ERROR: cannot marshal docKey bytes: %s", err)
					continue
				}

				if err = a.docsWriter.WriteBytes(
					ctx,
					ds.NewKeyValuePair(
						docKeyBytes,
						[]byte(result.HTMLDoc),
					),
				); err != nil {
					log.Printf("ERROR: unable to write json to kafka: %s", err)
					continue
				}

				if err = a.indexStorage.AddDoc(ctx, model.IndexDoc{
					Domain:          t.Key.Hostname,
					Path:            t.Value.Path,
					CRC32Checksum:   result.CRC32Checksum,
					Date:            time.Now(),
					StorageFileName: docKey.FileName,
				}); err != nil {
					log.Printf("ERROR: cannot write doc to index storage: %s", err)
					continue
				}
			}
		}
	}()
}

func (a *App) AddUrl(ctx context.Context, rawUrl string) error {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return fmt.Errorf("cannot parse url: %w", err)
	}

	return a.urlsWriter.WriteJSON(ctx, ds.NewKeyValuePair(model.CrawlTaskKey{Hostname: parsedUrl.Hostname()}, model.CrawlTaskValue{Path: rawUrl}))
}

func generateFileName(domain string) string {
	return fmt.Sprintf("%s-%s.html", domain, uuid.New().String())
}

func initPostgres(connString string) *sqlx.DB {
	conn := sqlx.MustOpen("pgx", connString)
	if err := conn.Ping(); err != nil {
		panic(fmt.Errorf("cannot ping postgres: %w", err))
	}

	if err := goose.Up(conn.DB, "./migration"); err != nil {
		panic(fmt.Errorf("cannot up migration: %w", err))
	}

	return conn
}
