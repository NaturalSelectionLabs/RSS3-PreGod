package misskey

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
)

type misskeyCrawler struct {
	crawler.DefaultCrawler
}

func NewMisskeyCrawler() crawler.Crawler {
	return &misskeyCrawler{
		crawler.DefaultCrawler{
			Assets: []model.Asset{},
			Notes:  []model.Note{},
		},
	}
}

func (mc *misskeyCrawler) Work(param crawler.WorkParam) error {
	noteList, err := GetUserNoteList(param.Identity, param.Limit, param.Timestamp)

	if err != nil {
		logger.Errorf("%v : unable to retrieve misskey note list for %s", err, param.Identity)

		return err
	}

	for _, item := range noteList {
		note := model.Note{
			Identifier:      rss3uri.NewNoteInstance(item.Id, constants.NetworkSymbolMisskey).UriString(),
			Owner:           item.Author,
			RelatedURLs:     []string{item.Link},
			Tags:            constants.ItemTagsMisskeyNote.ToPqStringArray(),
			Authors:         []string{item.Author},
			Summary:         item.Summary,
			Attachments:     database.MustWrapJSON(item.Attachments),
			Source:          constants.NoteSourceNameMisskeyNote.String(),
			MetadataNetwork: constants.NetworkSymbolMisskey.String(),
			MetadataProof:   item.Id,
			Metadata: database.MustWrapJSON(map[string]interface{}{
				"from": item.Author,
			}),
			DateCreated: item.CreatedAt,
			DateUpdated: item.CreatedAt, // TODO: check if updatedAt is available
		}

		mc.Notes = append(mc.Notes, note)
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
