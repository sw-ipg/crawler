package pipeline

import (
	"crawler/app/model"
	"database/sql"
	"fmt"
	"time"
)

func UrlDeduplicationPipeFunc(next CrawlPipe) CrawlPipeFunc {
	return func(ctx Context, task model.CrawlTask) (model.CrawlResult, error) {
		lastDate, err := ctx.indexStorage.GetLastDateOfDoc(ctx, task.PathToDiscover)
		if err != nil && err != sql.ErrNoRows {
			return failCrawlResult(), fmt.Errorf("cannot get last date of doc: %w", err)
		}

		if time.Since(lastDate) < 1*time.Hour {
			return alreadyCrawledStatus(), nil
		}

		return next.Do(ctx, task)
	}
}
