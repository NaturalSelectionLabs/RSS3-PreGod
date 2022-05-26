package main

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/poap"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler_handler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/ethereum/go-ethereum/common"
)

func MakeCrawlers[T constants.NetworkID | constants.PlatformID](network T) crawler.Crawler {
	switch any(network).(type) {
	case constants.NetworkID:
		switch constants.NetworkID(network) {
		case constants.NetworkIDGnosisMainnet:
			return poap.NewPoapCrawler()
		default:
			return nil
		}

	default:
		return nil
	}
}

func Excute(pt *crawler_handler.GetItemsHandler) (*crawler_handler.GetItemsResult, error) {
	var c crawler.Crawler

	var r *crawler.DefaultCrawler

	result := crawler_handler.NewGetItemsResult()

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

	tx := database.DB.Begin()
	defer tx.Rollback()

	if r.Notes != nil && len(r.Notes) > 0 {
		if dbNotes, err := database.CreateNotes(tx, r.Notes, true); err != nil {
			return result, err
		} else {
			r.Notes = dbNotes
		}
	}

	if r.Erc20Notes != nil && len(r.Erc20Notes) > 0 {
		if dbNotes, err := database.CreateNotesDoNothing(tx, r.Erc20Notes); err != nil {
			return result, err
		} else {
			r.Erc20Notes = dbNotes
		}
	}

	if err := tx.Commit().Error; err != nil {
		return result, err
	}

	result.Result = r

	return result, nil
}

func getOwnerFeed(instance rss3uri.Instance, owner string) {
	networkIDs := constants.GetEthereumPlatformNetworks()
	for _, networkID := range networkIDs {
		getItemHandler := crawler_handler.NewGetItemsHandler(crawler.WorkParam{
			Identity:   instance.GetIdentity(),
			PlatformID: constants.PlatformIDEthereum,
			NetworkID:  networkID,
			OwnerID:    owner,
		})

		_, err := Excute(getItemHandler)
		if err != nil {
			logger.Errorf("SubscribeEns: get item error, %v", err)

			continue
		}
	}
}

func main() {
	if err := database.Setup(); err != nil {
		logger.Errorf("subscribe.script: database.Setup err: %v", err)

		return
	}

	var db = database.DB

	var total int64

	for {
		domains := make([]model.Domains, 0)
		page := 0

		// get ens data
		if err := db.
			Where(&model.Domains{Type: "ens"}).
			Where("block_timestamp < ?", "2022-05-19").
			Order("block_timestamp DESC").
			Limit(1000).
			Offset(page * 1000).
			Find(&domains).Error; err != nil {
			logger.Errorf("subscribe.script: database get err: %v", err)

			return
		}

		if len(domains) == 0 {
			return
		}

		for _, ens := range domains {
			total += 1

			logger.Infof("total: ==== %v\n", total)

			address := common.BytesToAddress(ens.AddressOwner).String()
			instance, err := rss3uri.NewInstance("account", address, "ethereum")

			if err != nil {
				logger.Infof("get instance error: %v", err)

				continue
			}

			// get cache feed
			var count int64
			if err := db.Where("owner = ?", instance.UriString()).Model(&model.Note{}).Count(&count).Error; err == nil && count > 0 {
				continue
			}

			// get latest feed
			go func() {
				getOwnerFeed(instance, address)
			}()
		}

		page += 1
	}
}
