package pipeline

import (
	"crawler/app/model"
	"hash/crc32"
)

func ChecksumCalculationPipeFunc(next CrawlPipe) CrawlPipeFunc {
	return func(ctx Context, task model.CrawlTask) (model.CrawlResult, error) {
		ctx.crc32Checksum = crc32.ChecksumIEEE([]byte(ctx.htmlDoc))
		return next.Do(ctx, task)
	}
}
