package jike

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
)

type jikeCrawler struct {
	crawler.CrawlerResult
}

func NewJikeCrawler() crawler.Crawler {
	return &jikeCrawler{
		crawler.CrawlerResult{
			Items: []*model.Item{},
			Notes: []*model.ItemId{},
		},
	}
}

func (mc *jikeCrawler) Work(param crawler.WorkParam) error {
	timeline, err := GetUserTimeline(param.Identity)

	if err != nil {
		return *result, err
	}

	for _, item := range timeline {
		ni := model.NewItem(
			param.NetworkID,
			item.Link,
			model.Metadata{
				"network": constants.NetworkSymbolJike,
				"from":    item.Author,
			},
			constants.ItemTagsJikePost,
			[]string{item.Author},
			"",
			item.Summary,
			item.Attachments,
			item.Timestamp,
		)
		mc.Items = append(mc.Items, ni)

		mc.Notes = append(mc.Notes, &model.ItemId{
			NetworkID: param.NetworkID,
			Proof:     item.Link,
		})
	}

	return nil
}

func (mc *jikeCrawler) GetResult() *crawler.CrawlerResult {
	return &crawler.CrawlerResult{
		Notes: mc.Notes,
		Items: mc.Items,
	}
}

func (tc *jikeCrawler) GetUserBio(Identity string) (string, error) {
	if err := Login(); err != nil {
		return "", err
	}

	userProfile, err := GetUserProfile(Identity)

	if err != nil {
		return "", err
	}

	userBios := []string{userProfile.Bio}
	userBioJson, err := crawler.GetUserBioJson(userBios)

	if err != nil {
		return "", err
	}

	return userBioJson, nil
}
