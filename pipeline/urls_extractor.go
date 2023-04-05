package pipeline

import (
	"crawler/app/model"
	"crawler/ds"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ExtractUrlsPipeFunc(next CrawlPipe) CrawlPipeFunc {
	return func(ctx Context, task model.CrawlTask) (model.CrawlResult, error) {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(ctx.htmlDoc))
		if err != nil {
			return failCrawlResult(), fmt.Errorf("cannot get document from html: %w", err)
		}

		taskPath, _ := url.Parse(task.PathToDiscover)
		doc.Find("a").Each(func(i int, selection *goquery.Selection) {
			href, ok := selection.Attr("href")
			if !ok {
				return
			}

			var hostname string
			if u, err := url.Parse(href); err != nil {
				log.Printf("ERROR: invalid URL: %s", err)
				return
			} else {
				hostname = u.Hostname()
			}

			if hostname != taskPath.Hostname() { // only internal links
				return
			}

			if err = ctx.workItemsWriter.WriteJSON(
				ctx,
				ds.NewKeyValuePair(
					model.CrawlTaskKey{
						Hostname: hostname,
					},
					model.CrawlTaskValue{
						Path: href,
					},
				),
			); err != nil {
				log.Printf("ERROR: cannot write to kafka: %s", err)
				return
			}
		})

		return next.Do(ctx, task)
	}
}
