package crawler_handler

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
)

type GetItemsHandler struct {
	CrawlerHandlerBase
}

type GetItemsResult struct {
	CrawlerHandlerResultBase
	Result *crawler.DefaultCrawler
}

func NewGetItemsHandler(workParam crawler.WorkParam) *GetItemsHandler {
	return &GetItemsHandler{
		CrawlerHandlerBase{
			WorkParam: workParam,
		},
	}
}

func NewGetItemsResult() *GetItemsResult {
	return &GetItemsResult{
		CrawlerHandlerResultBase{
			Error: util.GetErrorBase(util.ErrorCodeSuccess),
		},
		nil,
	}
}

func (pt *GetItemsHandler) Excute() *GetItemsResult {
	var err error

	var c crawler.Crawler

	var r *crawler.DefaultCrawler

	result := NewGetItemsResult()

	instance := rss3uri.NewAccountInstance(pt.WorkParam.Identity, pt.WorkParam.PlatformID.Symbol())

	c = MakeCrawlers(pt.WorkParam.NetworkID)
	if c == nil {
		result.Error = util.GetErrorBase(util.ErrorCodeNotSupportedNetwork)

		logger.Errorf("unsupported network id[%d]", pt.WorkParam.NetworkID)

		return result
	}

	err = c.Work(pt.WorkParam)

	if err != nil {
		result.Error = util.GetErrorBase(util.ErrorCodeNotSupportedNetwork)

		logger.Errorf("crawler fails while working: %s", err)

		return result
	}

	r = c.GetResult()
	if r.Items != nil {
		for _, item := range r.Items {
			db.InsertItem(item)
		}
	}

	if r.Assets != nil {
		db.SetAssets(instance, r.Assets, pt.WorkParam.NetworkID)
	}

	if r.Notes != nil {
		db.AppendNotes(instance, r.Notes)
	}

	result.Result = r

	return result
}
