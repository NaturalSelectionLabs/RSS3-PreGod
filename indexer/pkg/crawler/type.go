package crawler

import (
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	jsoniter "github.com/json-iterator/go"
)

type Crawler interface {
	Work(param WorkParam) error
	// GetResult return &{Assets, Notes, Items}
	GetResult() *DefaultCrawler
	// GetBio
	// Since some apps have multiple bios,
	// they need to be converted into json and then collectively transmitted
	GetUserBio(Identity string) (string, error)
}

type DefaultCrawler struct {
	Assets   []model.Asset
	Notes    []model.Note
	Profiles []model.Profile
}

// CrawlerResult inherits the function by default

func (cr *DefaultCrawler) Work(param WorkParam) error {
	return nil
}

func (cr *DefaultCrawler) GetResult() *DefaultCrawler {
	return cr
}

func (cr *DefaultCrawler) GetUserBio(Identity string) (string, error) {
	return "", nil
}

type WorkParam struct {
	Identity   string
	NetworkID  constants.NetworkID
	PlatformID constants.PlatformID // optional
	Limit      int                  // optional, aka Count, limit the number of items to be crawled

	Timestamp time.Time // optional, if provided, only index items newer than this time
}

type userBios struct {
	Bios []string `json:"bios"`
}

func GetUserBioJson(bios []string) (string, error) {
	jsoni := jsoniter.ConfigCompatibleWithStandardLibrary

	userbios := userBios{Bios: bios}
	userBioJson, err := jsoni.MarshalToString(userbios)

	if err != nil {
		return "", err
	}

	return userBioJson, nil
}
