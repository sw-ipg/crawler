package pipeline

func NewPipeline() CrawlPipe {
	return UrlDeduplicationPipeFunc(PageDownloaderPipeFunc(ChecksumCalculationPipeFunc(DocDeduplicationPipeFunc(ExtractUrlsPipeFunc(EndCrawlPipeFunc())))))
}
