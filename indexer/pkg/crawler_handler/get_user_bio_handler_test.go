package crawler_handler_test

import (
	"strings"
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler_handler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
)

func TestGetBio(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		identity   string
		platformID constants.PlatformID
		success    bool
		contains   string
	}{
		{
			name:       "misskey-ok",
			identity:   "song@misskey.io",
			platformID: constants.PlatformIDMisskey,
			success:    true,
			contains:   "song.cheers.bio",
		},
		{
			name:       "twitter-ok",
			identity:   "diygod",
			platformID: constants.PlatformIDTwitter,
			success:    true,
			contains:   "diygod.me",
		},
		{
			name:       "jike-ok",
			identity:   "C05E4867-4251-4F11-9096-C1D720B41710",
			platformID: constants.PlatformIDJike,
			success:    true,
			contains:   "henryqw.eth",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			result, err := request(tt.identity, tt.platformID)
			if tt.success {
				if result.Error.ErrorCode != 0 {
					t.Errorf("Expected success, but got error: %+v", result.Error)
				}

				if !strings.Contains(result.UserBio, tt.contains) {
					t.Errorf("Expected bio to contain %s, but got %s", tt.contains, result.UserBio)
				}
			} else {
				if result.Error.ErrorCode == 0 {
					t.Errorf("Expected not success, but got response: %+v, err: %s", result, err)
				}
			}
		})
	}
}

func request(identity string, platformID constants.PlatformID) (*crawler_handler.GetUserBioResult, error) {
	getuserBioHandler := crawler_handler.NewGetUserBioHandler(
		crawler.WorkParam{
			Identity:   identity,
			PlatformID: platformID,
		})

	handlerResult, err := getuserBioHandler.Excute()
	if err != nil {
		return nil, err
	}

	return handlerResult, nil
}
