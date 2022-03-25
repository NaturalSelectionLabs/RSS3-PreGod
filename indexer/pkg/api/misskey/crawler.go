package misskey

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
)

type misskeyCrawler struct {
	crawler.DefaultCrawler
}

func NewMisskeyCrawler() crawler.Crawler {
	return &misskeyCrawler{
		crawler.DefaultCrawler{
			Items: []*model.Item{},
			Notes: []*model.ItemId{},
		},
	}
}

func (mc *misskeyCrawler) Work(param crawler.WorkParam) error {
	noteList, err := GetUserNoteList(param.Identity, param.Limit, param.TimeStamp)

	if err != nil {
		logger.Errorf("%v : unable to retrieve misskey note list for %s", err, param.Identity)

		return err
	}

	for _, note := range noteList {
		ni := model.NewItem(
			param.NetworkID,
			note.Link,
			model.Metadata{
				"network": constants.NetworkSymbolMisskey,
				"from":    note.Author,
			},
			constants.ItemTagsMisskeyNote,
			[]string{note.Author},
			"",
			note.Summary,
			note.Attachments,
			note.CreatedAt,
		)
		mc.Items = append(mc.Items, ni)

		mc.Notes = append(mc.Notes, &model.ItemId{
			NetworkID: param.NetworkID,
			Proof:     note.Link,
		})
	}

	return nil
}

func (mc *misskeyCrawler) GetUserBio(Identity string) (string, error) {
	accountInfo, err := formatUserAccount(Identity)
	if err != nil {
		return "", err
	}

	userShow, err := GetUserShow(accountInfo)

	if err != nil {
		return "", err
	}

	userBios := userShow.Bios
	userBioJson, err := crawler.GetUserBioJson(userBios)

	if err != nil {
		return "", err
	}

	return userBioJson, nil
}
