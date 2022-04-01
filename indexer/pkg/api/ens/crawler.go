package ens

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
)

type ensCrawler struct {
	crawler.DefaultCrawler
}

func NewEnsCrawler() crawler.Crawler {
	return &ensCrawler{
		crawler.DefaultCrawler{
			Profiles: []*model.Profile{},
		},
	}
}

func (c *ensCrawler) Work(param crawler.WorkParam) error {
	ensList, err := GetENSList(param.Identity)

	if err != nil {
		return err
	}

	for _, ens := range ensList {
		metadata := make(map[string]interface{}, len(ens.Text))
		for k, v := range ens.Text {
			metadata[k] = v
		}

		profile := model.NewProfile(
			param.NetworkID,
			ens.TxHash,
			ens.Text,
			ens.Domain,
			ens.Description,
			[]string{ens.Text["avatar"]},
			[]model.Attachment{},
			[]string{param.Identity},
		)

		c.Profiles = append(c.Profiles, profile)
	}

	return nil
}
