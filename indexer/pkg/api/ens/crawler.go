package ens

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
)

type ensCrawler struct {
	crawler.DefaultCrawler
}

func NewEnsCrawler() crawler.Crawler {
	return &ensCrawler{
		crawler.DefaultCrawler{
			Items:  []*model.Item{},
			Assets: []*model.ItemId{},
		},
	}
}

func (c *ensCrawler) Work(param crawler.WorkParam) error {
	ensList := GetENSList(param.Identity)
	for _, ens := range ensList {

		metadata := make(map[string]interface{}, len(ens.text))
		for k, v := range ens.text {
			metadata[k] = v
		}

		item := model.NewItem(
			param.NetworkID,
			ens.domain,
			metadata,
			constants.ItemTagENS,
			[]string{param.Identity},
			"",
			"",
			[]model.Attachment{},
			ens.createdAt,
		)

		c.Items = append(c.Items, item)

		c.Assets = append(c.Assets, &model.ItemId{
			NetworkID: param.NetworkID,
			Proof:     ens.txHash,
		})

	}

	return nil
}
