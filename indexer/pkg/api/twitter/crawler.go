package twitter

import (
	"fmt"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
)

type twitterCrawler struct {
	crawler.DefaultCrawler
}

func NewTwitterCrawler() crawler.Crawler {
	return &twitterCrawler{
		crawler.DefaultCrawler{
			Assets: []model.Asset{},
			Notes:  []model.Note{},
		},
	}
}

func (tc *twitterCrawler) Work(param crawler.WorkParam) error {
	if param.NetworkID != constants.NetworkIDTwitter {
		return fmt.Errorf("network is not twitter")
	}

	contentInfos, err := GetTimeline(param.Identity, uint32(param.Limit))
	if err != nil {
		logger.Error(err)

		return err
	}

	for _, item := range contentInfos {
		tsp, err := item.GetTsp()
		if err != nil {
			// TODO: log error
			logger.Error(tsp, err)
			tsp = time.Now()
		}

		author := rss3uri.NewAccountInstance(item.ScreenName, constants.PlatformSymbolTwitter).UriString()

		note := model.Note{
			Identifier:      rss3uri.NewNoteInstance(item.Hash, constants.NetworkSymbolTwitter).UriString(),
			Owner:           author,
			RelatedURLs:     []string{item.Link},
			Tags:            constants.ItemTagsTweet.ToPqStringArray(),
			Authors:         []string{author},
			Summary:         item.PreContent,
			Attachments:     database.MustWrapJSON(item.Attachments),
			Source:          constants.NoteSourceNameTwitterTweet.String(),
			MetadataNetwork: constants.NetworkSymbolTwitter.String(),
			MetadataProof:   item.Hash,
			DateCreated:     tsp,
			DateUpdated:     tsp, // TODO: does twitter support updating tweets?
		}

		tc.Notes = append(tc.Notes, note)
	}

	return nil
}

func (tc *twitterCrawler) GetUserBio(Identity string) (string, error) {
	userShow, err := GetUserShow(Identity)

	if err != nil {
		return "", err
	}

	userBios := []string{userShow.Description, userShow.Entities}
	userBioJson, err := crawler.GetUserBioJson(userBios)

	if err != nil {
		return "", err
	}

	return userBioJson, nil
}
