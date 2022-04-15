package crawler_handler

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
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

	// the error here does not affect the execution of the crawler
	if dbQcmErr != nil {
		pt.WorkParam.BlockHeight = metadata.LastBlock
		pt.WorkParam.Timestamp = metadata.UpdatedAt
	}

	if err := c.Work(pt.WorkParam); err != nil {
		result.Error = util.GetErrorBase(util.ErrorCodeNotSupportedNetwork)

		return result, fmt.Errorf("crawler fails while working: %s", err)
	}

	r = c.GetResult()

	tx := database.DB.Begin()
	defer tx.Rollback()

	if r.Assets != nil && len(r.Assets) > 0 {
		if dbAssets, err := database.CreateAssets(tx, r.Assets, true); err != nil {
			return result, err
		} else {
			r.Assets = dbAssets
		}
	}

	if r.Notes != nil && len(r.Notes) > 0 {
		if dbNotes, err := database.CreateNotes(tx, r.Notes, true); err != nil {
			return result, err
		} else {
			r.Notes = dbNotes
		}
	}

	if r.Profiles != nil && len(r.Profiles) > 0 {
		if dbProfiles, err := database.CreateProfiles(tx, r.Profiles, true); err != nil {
			return result, err
		} else {
			r.Profiles = dbProfiles
		}
	}

	// TODO: stores the crawler last worked metadata
	if _, err := database.CreateCrawlerMetadata(tx, &model.CrawlerMetadata{
		AccountInstance: pt.WorkParam.OwnerID,
		PlatformID:      pt.WorkParam.PlatformID,
	}, true); err != nil {
		return result, err
	}

	if err := tx.Commit().Error; err != nil {
		return result, err
	}

	result.Result = r

	return result, nil
}
