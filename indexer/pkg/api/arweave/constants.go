package arweave

import "time"

type ArAccount string

const (
	MirrorUploader ArAccount = "Ky1c1Kkt-jZ9sY1hvLF5nCf6WWdBhIU5Un_BMYh-t3c"
)

const (
	DefaultCrawlStep     = 100
	DefaultFromHeight    = 571518 // MirrorUploader account was created at block #559678
	DefaultConfirmations = 10
	DefaultCrawlMinStep  = 10
)

var DefaultCrawlConfig = &crawlConfig{
	DefaultFromHeight,
	DefaultConfirmations,
	DefaultCrawlStep,
	DefaultCrawlMinStep,
	60 * time.Second,
}
