package main

import (
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler_handler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/ethereum/go-ethereum/common"
)

var success int64
var db = database.DB
var total int64

func getOwnerFeed(instance rss3uri.Instance, owner string) {
	var networkIDs = constants.GetEthereumPlatformNetworks()

	var err error

	for _, networkID := range networkIDs {
		getItemHandler := crawler_handler.NewGetItemsHandler(crawler.WorkParam{
			Identity:   instance.GetIdentity(),
			PlatformID: constants.PlatformIDEthereum,
			NetworkID:  networkID,
			OwnerID:    owner,
		})

		_, err = getItemHandler.Excute()
		if err != nil {
			logger.Errorf("subscribe.script:: get item error, %v", err)

			continue
		}
	}

	if err == nil {
		success += 1

		logger.Infof("subscribe.script: success load user feed, count = %v", success)
	}
}

func main() {
	if err := database.Setup(); err != nil {
		logger.Errorf("subscribe.script: database.Setup err: %v", err)

		return
	}

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
				logger.Infof("note already exists, owner = %v", instance.UriString())

				continue
			}

			// get latest feed
			getOwnerFeed(instance, address)

			time.Sleep(10 * time.Second)
		}

		page += 1
	}
}
