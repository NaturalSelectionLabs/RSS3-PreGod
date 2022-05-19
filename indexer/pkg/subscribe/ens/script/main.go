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
)

func getOwnerFeed(owner string) {
	instance, err := rss3uri.NewInstance("account", owner, "ethereum")
	if err != nil {
		return
	}

	networkIDs := constants.GetEthereumPlatformNetworks()
	for _, networkID := range networkIDs {
		getItemHandler := crawler_handler.NewGetItemsHandler(crawler.WorkParam{
			Identity:   instance.GetIdentity(),
			PlatformID: constants.PlatformIDEthereum,
			NetworkID:  networkID,
			OwnerID:    owner,
		})

		_, err = getItemHandler.Excute()
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

	db := database.DB

	for {
		domains := make([]model.Domains, 0)
		page := 0

		if err := db.Where(&model.Domains{Type: "ens"}).Limit(1000).Offset(page * 1000).Find(&domains).Error; err != nil {
			logger.Errorf("subscribe.script: database get err: %v", err)

			return
		}

		if len(domains) == 0 {
			return
		}

		for _, ens := range domains {
			go func() {
				getOwnerFeed(string(ens.AddressOwner))
			}()

			time.Sleep(30 * time.Second)
		}

		page += 1
	}
}
