package jike

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
)

type jikeCrawler struct {
	crawler.DefaultCrawler
}

func NewJikeCrawler() crawler.Crawler {
	return &jikeCrawler{
		crawler.DefaultCrawler{
			Assets: []model.Asset{},
			Notes:  []model.Note{},
		},
	}
}

func (mc *jikeCrawler) Work(param crawler.WorkParam) error {
	timeline, err := GetUserTimeline(param.Identity)

	if err != nil {
		return err
	}

	owner := rss3uri.NewAccountInstance(param.OwnerID, param.OwnerPlatformID.Symbol()).UriString()
	author := rss3uri.NewAccountInstance(param.Identity, constants.PlatformSymbolJike).UriString()

	for _, item := range timeline {
		note := model.Note{
			Identifier:      rss3uri.NewNoteInstance(item.Id, constants.NetworkSymbolJike).UriString(),
			Owner:           owner,
			RelatedURLs:     []string{item.Link},
			Tags:            constants.ItemTagsJikePost.ToPqStringArray(),
			Authors:         []string{owner, author},
			Summary:         item.Summary,
			Attachments:     database.MustWrapJSON(item.Attachments),
			Source:          constants.NoteSourceNameJikePost.String(),
			MetadataNetwork: constants.NetworkSymbolJike.String(),
			MetadataProof:   item.Id,
			Metadata:        database.MustWrapJSON(map[string]interface{}{
				// "from": item.Author,
			}),
			DateCreated: item.Timestamp,
			DateUpdated: item.Timestamp,
		}

		mc.Notes = append(mc.Notes, note)
	}

	return nil
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
