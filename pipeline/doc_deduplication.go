package pipeline

import (
	"crawler/app/model"
	"fmt"
)

func DocDeduplicationPipeFunc(next CrawlPipe) CrawlPipeFunc {
	return func(ctx Context, task model.CrawlTask) (model.CrawlResult, error) {
		checksums, err := ctx.indexStorage.GetChecksumsForPath(ctx, task.PathToDiscover)
		if err != nil {
			return failCrawlResult(), fmt.Errorf("cannot get checksums for path: %w", err)
		}

		for _, c := range checksums {
			if c == ctx.crc32Checksum {
				return alreadyCrawledStatus(), nil
			}
		}

		return next.Do(ctx, task)
	}
}
