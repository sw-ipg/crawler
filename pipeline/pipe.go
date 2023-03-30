package pipeline

import (
	"context"
	"crawler/app/model"
	"crawler/communication"
)

type Context struct {
	context.Context

	workItemsWriter communication.KafkaTopicWriter[model.CrawlTaskKey, model.CrawlTaskValue]
	indexStorage    communication.IndexMetadataStorage

	htmlDoc       string
	crc32Checksum uint32
}

func NewContext(
	ctx context.Context,
	workItemsWriter communication.KafkaTopicWriter[model.CrawlTaskKey, model.CrawlTaskValue],
	indexStorage communication.IndexMetadataStorage,
) Context {
	return Context{
		Context:         ctx,
		workItemsWriter: workItemsWriter,
		indexStorage:    indexStorage,
	}
}

func failCrawlResult() model.CrawlResult {
	return model.CrawlResult{
		Status: model.FailedCrawlStatus,
	}
}

func alreadyCrawledStatus() model.CrawlResult {
	return model.CrawlResult{
		Status: model.AlreadyCrawledStatus,
	}
}

type CrawlPipe interface {
	Do(ctx Context, task model.CrawlTask) (model.CrawlResult, error)
}

type CrawlPipeFunc func(ctx Context, task model.CrawlTask) (model.CrawlResult, error)

func (f CrawlPipeFunc) Do(ctx Context, task model.CrawlTask) (model.CrawlResult, error) {
	return f(ctx, task)
}

func EndCrawlPipeFunc() CrawlPipeFunc {
	return func(ctx Context, task model.CrawlTask) (model.CrawlResult, error) {
		return model.CrawlResult{
			Status:        model.SuccessCrawlStatus,
			HTMLDoc:       ctx.htmlDoc,
			CRC32Checksum: ctx.crc32Checksum,
		}, nil
	}
}
