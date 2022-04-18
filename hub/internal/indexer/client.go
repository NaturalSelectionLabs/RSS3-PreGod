package indexer

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/go-resty/resty/v2"
	"golang.org/x/sync/errgroup"
)

func GetItems(instance rss3uri.Instance, accounts []model.Account) error {
	eg := errgroup.Group{}

	// Add self
	accounts = append(accounts, model.Account{
		Identity:        strings.ToLower(instance.GetIdentity()),
		Platform:        constants.PlatformSymbol(instance.GetSuffix()).ID().Int(),
		ProfileID:       strings.ToLower(instance.GetIdentity()),
		ProfilePlatform: constants.PlatformSymbol(instance.GetSuffix()).ID().Int(),
		Source:          int(constants.NetworkIDCrossbell),
	})

	for _, account := range accounts {
		account := account

		for _, networkID := range constants.GetNetworkList(constants.PlatformID(account.Platform)) {
			networkID := networkID
			client := resty.New()

			eg.Go(func() error {
				start := time.Now()
				request := client.NewRequest()
				params := map[string]string{
					"proof":             strings.ToLower(account.Identity),
					"platform_id":       strconv.Itoa(account.Platform),
					"network_id":        strconv.Itoa(int(networkID)),
					"profile_source_id": strconv.Itoa(account.Source),
					"owner_id":          strings.ToLower(account.ProfileID),
					"owner_platform_id": strconv.Itoa(account.ProfilePlatform),
				}
				result := Response{}
				response, err := request.
					SetQueryParams(params).
					SetResult(&result).
					Get(fmt.Sprintf("%s/item", config.Config.Hub.IndexerEndpoint))
				if err != nil {
					logger.Error(err)

					return nil
				}

				// TODO Remove it
				logger.Info(params["proof"], "\x20", params["platform_id"], "\x20", params["network_id"], "\x20", start, "\x20", time.Now().Sub(start))

				if response.StatusCode() != http.StatusOK || result.Error.Code != 0 {
					logger.Error(response.StatusCode(), result.Error.Code, result.Error.Msg)

					return nil
				}

				return nil
			})
		}
	}

	return eg.Wait()
}
