package crawler_handler

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
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
	var c crawler.Crawler

	var r *crawler.DefaultCrawler

	result := NewGetItemsResult()

	c = MakeCrawlers(pt.WorkParam.NetworkID)
	if c == nil {
		result.Error = util.GetErrorBase(util.ErrorCodeNotSupportedNetwork)

		return result, fmt.Errorf("unsupported network id[%d]", pt.WorkParam.NetworkID)
	}

	metadata, dbQcmErr := database.QueryCrawlerMetadata(database.DB, pt.WorkParam.Identity, pt.WorkParam.PlatformID)

	// Historical legacy, the code here is no longer needed, LastBlock = 0
	// the error here does not affect the execution of the crawler
	if dbQcmErr != nil && metadata != nil {
		pt.WorkParam.BlockHeight = metadata.LastBlock
		pt.WorkParam.Timestamp = metadata.UpdatedAt
	}

	if err := c.Work(pt.WorkParam); err != nil {
		result.Error = util.GetErrorBase(util.ErrorCodeNotSupportedNetwork)

		return result, fmt.Errorf("crawler fails while working: %s", err)
	}

	r = c.GetResult()

	db := database.DB

	if r.Notes != nil && len(r.Notes) > 0 {
		if dbNotes, err := database.CreateNotes(db, r.Notes, true); err != nil {
			return result, err
		} else {
			r.Notes = dbNotes
		}
	}

	if r.Erc20Notes != nil && len(r.Erc20Notes) > 0 {
		if dbNotes, err := database.CreateNotesDoNothing(db, r.Erc20Notes); err != nil {
			return result, err
		} else {
			r.Erc20Notes = dbNotes
		}
	}

	go func() {
		if r.Assets != nil && len(r.Assets) > 0 {
			if _, err := database.CreateAssets(db, r.Assets, true); err != nil {
				logger.Errorf("database.CreateAssets error: %v", err)
			}
		}
	}()

	result.Result = r

	return result, nil
}
