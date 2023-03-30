package model

type CrawlTask struct {
	PathToDiscover string
}

type CrawlStatus int

const (
	SuccessCrawlStatus CrawlStatus = iota
	FailedCrawlStatus
	AlreadyCrawledStatus
)

type CrawlResult struct {
	Status        CrawlStatus
	HTMLDoc       string
	CRC32Checksum uint32
}
