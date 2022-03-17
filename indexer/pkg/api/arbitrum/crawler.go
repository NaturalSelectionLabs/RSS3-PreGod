package arbitrum

import (
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
)

//nolint:funlen // disable line length check
func Crawl(param *crawler.WorkParam, result *crawler.CrawlerResult) (crawler.CrawlerResult, error) {
	nftTransfers, err := GetNFTTransfers(param.Identity)
	if err != nil {
		return *result, err
	}

	assets, err := GetNFTs(param.Identity)
	if err != nil {
		return *result, err
	}

	networkId := constants.NetworkSymbolArbitrum.GetID()

	// parse notes
	for _, v := range nftTransfers {
		result.Notes = append(result.Notes, &model.ItemId{
			NetworkID: networkId,
			Proof:     v.Hash,
		})
	}

	// parse assets
	for _, v := range assets {
		hasProof := false

		for _, nftTransfer := range nftTransfers {
			if nftTransfer.EqualsToToken(v) {
				hasProof = true

				result.Assets = append(result.Assets, &model.ItemId{
					NetworkID: networkId,
					Proof:     nftTransfer.Hash,
				})
			}
		}

		if !hasProof {
			logger.Error("Asset doesn't has proof")
		}
	}

	for _, v := range nftTransfers {
		tsp, err := time.Parse(time.RFC3339, v.TimeStamp)
		if err != nil {
			tsp = time.Now()
		}

		author, err := rss3uri.NewInstance("account", v.From, string(constants.NetworkSymbolArbitrum))
		if err != nil {
			// TODO
			logger.Error(err)
		}

		hasObject := false
		attachments := make([]model.Attachment, 0)

		for _, asset := range assets {
			if v.EqualsToToken(asset) && asset.MetaData != "" {
				hasObject = true
			}
		}

		if !hasObject {
			// TODO: get object
			logger.Errorf("Asset doesn't has the metadata.")
		}

		item := model.NewItem(
			networkId,
			v.Hash,
			model.Metadata{
				"network": constants.NetworkSymbolArbitrum,
				"from":    v.From,
				"to":      v.To,
			},
			constants.ItemTagsNFT,
			[]string{author.String()},
			"",
			"",
			attachments,
			tsp,
		)
		result.Items = append(result.Items, item)
	}

	return *result, nil
}
