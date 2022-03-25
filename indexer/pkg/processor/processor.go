package processor

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/jike"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/misskey"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/moralis"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/twitter"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
)

type ProcessTaskHandler interface {
	Fun() error
}

type ProcessTaskParam struct {
	ProcessTaskHandler
	TaskType  ProcessTaskType
	WorkParam crawler.WorkParam
}

type ProcessTaskResult struct {
	TaskType   ProcessTaskType
	TaskResult ProcessTaskErrorCode
}

type Processor struct {
	// Emergency use, highest priority, such as user data not found
	UrgentQ chan ProcessTaskHandler
	// General use, such as access to authenticate user information
	HighQ chan ProcessTaskHandler
	// Unaffected condition use, such as polling query data
	LowQ chan ProcessTaskHandler
}

var GlobalProcessor *Processor

func Setup() error {
	GlobalProcessor = NewProcessor()
	go GlobalProcessor.ListenAndServe()

	return nil
}

func NewProcessor() *Processor {
	processor := new(Processor)

	processor.UrgentQ = make(chan ProcessTaskHandler)
	processor.HighQ = make(chan ProcessTaskHandler)
	processor.LowQ = make(chan ProcessTaskHandler)

	logger.Infof("NewProcessor init:%v", processor)

	return processor
}

func MakeCrawler(network constants.NetworkID) crawler.Crawler {
	switch network {
	case constants.NetworkIDEthereumMainnet,
		constants.NetworkIDBNBChain,
		constants.NetworkIDAvalanche,
		constants.NetworkIDFantom,
		constants.NetworkIDPolygon:
		return moralis.NewMoralisCrawler()
	case constants.NetworkIDJike:
		return jike.NewJikeCrawler()
	case constants.NetworkIDTwitter:
		return twitter.NewTwitterCrawler()
	case constants.NetworkIDMisskey:
		return misskey.NewMisskeyCrawler()
	default:
		return nil
	}
}

func (w *Processor) ListenAndServe() {
	for {
		select {
		case t := <-w.UrgentQ:
			t.Fun()
		case t := <-w.HighQ:
			t.Fun()
		case t := <-w.LowQ:
			t.Fun()
		}
	}
}
