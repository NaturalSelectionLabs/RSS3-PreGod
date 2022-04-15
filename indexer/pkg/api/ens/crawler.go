package ens

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
)

type ensCrawler struct {
	crawler.DefaultCrawler
}

func NewEnsCrawler() crawler.Crawler {
	return &ensCrawler{
		crawler.DefaultCrawler{
			Profiles: []model.Profile{},
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

		profile := model.Profile{
			ID:          param.Identity,
			Source:      constants.ProfileSourceIDENS.Int(),
			Name:        database.WrapNullString(ens.Domain),
			Bio:         database.WrapNullString(ens.Description),
			Avatars:     []string{ens.Text["avatar"]},
			Attachments: ens.Attachments,
		}

		c.Profiles = append(c.Profiles, profile)
	}

	return nil
}
