package pipeline

import (
	"crawler/app/model"
	"fmt"
	"io"
	"log"
	"net/http"
)

func PageDownloaderPipeFunc(next CrawlPipe) CrawlPipeFunc {
	return func(ctx Context, task model.CrawlTask) (model.CrawlResult, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, task.PathToDiscover, nil)
		if err != nil {
			return failCrawlResult(), fmt.Errorf("cannot create request: %w", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return failCrawlResult(), fmt.Errorf("cannot do req: %w", err)
		}

		defer func() {
			if err = resp.Body.Close(); err != nil {
				log.Printf("cannot close resp body: %s", err)
			}
		}()

		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return failCrawlResult(), fmt.Errorf("cannot read resp: %s", err)
		}

		ctx.htmlDoc = string(bytes)
		return next.Do(ctx, task)
	}
}
