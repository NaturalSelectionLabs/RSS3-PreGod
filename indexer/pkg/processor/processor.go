package processor

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/jike"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/moralis"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
)

type processor struct {
	lowQ, highQ chan *crawler.WorkParam
}

type Processor interface {
	processTask(t *crawler.WorkParam) error

	ListenAndServe()
}

func NewProcessor(lq, hq chan *crawler.WorkParam) Processor {
	return &processor{lq, hq}
}

func makeCrawlers(network constants.NetworkID) crawler.Crawler {
	switch network {
	case constants.NetworkIDEthereumMainnet,
		constants.NetworkIDBNBChain,
		constants.NetworkIDAvalanche,
		constants.NetworkIDFantom,
		constants.NetworkIDPolygon:
		return moralis.NewMoralisCrawler()
	case constants.NetworkIDJike:
		return jike.NewJikeCrawler()
	default:
		return nil
	}
}

func (w *processor) processTask(t *crawler.WorkParam) error {
	var err error

	var c crawler.Crawler

	var r *crawler.CrawlerResult

	var processorParam = crawler.WorkParam{
		Identity:  t.Identity,
		NetworkID: t.NetworkID,
	}

	instance := rss3uri.NewAccountInstance(t.Identity, t.PlatformID.Symbol())

	c = makeCrawlers(t.NetworkID)
	if c == nil {
		err = fmt.Errorf("unsupported network id: %d", t.NetworkID)

		goto RETURN
	}

	err = c.Work(processorParam)

	if err != nil {
		err = fmt.Errorf("crawler fails while working: %s", err)

		goto RETURN
	}

	r = c.GetResult()
	if r.Items != nil {
		for _, item := range r.Items {
			db.InsertItem(item)
		}
	}

	if r.Assets != nil {
		db.SetAssets(instance, r.Assets, t.NetworkID)
	}

	if r.Notes != nil {
		db.AppendNotes(instance, r.Notes)
	}

RETURN:
	if err != nil {
		logger.Error(err)

		return err
	} else {
		return nil
	}
}

func (w *processor) ListenAndServe() {
	select {
	case t := <-w.highQ:
		w.processTask(t)
	default:
		select {
		case t := <-w.highQ:
			w.processTask(t)
		case t := <-w.lowQ:
			w.processTask(t)
		}
	}
}
