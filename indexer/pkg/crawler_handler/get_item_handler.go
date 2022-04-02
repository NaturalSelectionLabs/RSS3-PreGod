package crawler_handler

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
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

func (pt *GetItemsHandler) Excute() (*GetItemsResult, error) {
	var err error

	var c crawler.Crawler

	var r *crawler.DefaultCrawler

	result := NewGetItemsResult()

	c = MakeCrawlers(pt.WorkParam.NetworkID)
	if c == nil {
		result.Error = util.GetErrorBase(util.ErrorCodeNotSupportedNetwork)

		return result, fmt.Errorf("unsupported network id[%d]", pt.WorkParam.NetworkID)
	}

	err = c.Work(pt.WorkParam)

	if err != nil {
		result.Error = util.GetErrorBase(util.ErrorCodeNotSupportedNetwork)

		return result, fmt.Errorf("crawler fails while working: %s", err)
	}

	r = c.GetResult()

	tx := database.DB.Begin()
	defer tx.Rollback()

	if r.Assets != nil {
		if _, err := database.CreateAssets(tx, r.Assets, true); err != nil {
			return result, err
		}
	}

	if r.Notes != nil {
		if _, err := database.CreateNotes(tx, r.Notes, true); err != nil {
			return result, err
		}
	}

	if err = tx.Commit().Error; err != nil {
		return result, err
	}

	result.Result = r

	return result, nil
}
