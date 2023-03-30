package model

type CrawlTaskKey struct {
	Hostname string `json:"domain"`
}

type CrawlTaskValue struct {
	Path string `json:"path"`
}

type DocKey struct {
	FileName string `json:"fileName"`
}

type DocValue = string
