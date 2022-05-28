package indexer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/cache"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"golang.org/x/sync/errgroup"
)

func GetItems(requestURL string, instance rss3uri.Instance, latest bool) error {
	lockerKey := fmt.Sprintf("hub %s", requestURL)

	if _, err := cache.GetRaw(context.Background(), lockerKey); err != nil && errors.Is(err, redis.Nil) {
		if err = cache.SetRaw(context.Background(), lockerKey, time.Now().String(), time.Second*10); err != nil {
			return err
		}

		if latest {
			return getItems(instance)
		}

		go func() {
			if err := getItems(instance); err != nil {
				logger.Error(err)
			}
		}()
	}

	return nil
}

func getItems(instance rss3uri.Instance) error {
	eg := errgroup.Group{}

	for _, networkID := range constants.GetNetworkList(constants.PlatformIDEthereum) {
		networkID := networkID
		client := resty.New()

		eg.Go(func() error {
			return getItem(client, model.Account{
				Identity:        strings.ToLower(instance.GetIdentity()),
				Platform:        constants.PlatformSymbol(instance.GetSuffix()).ID().Int(),
				ProfileID:       strings.ToLower(instance.GetIdentity()),
				ProfilePlatform: constants.PlatformSymbol(instance.GetSuffix()).ID().Int(),
				Source:          int(constants.NetworkIDCrossbell),
			}, networkID)
		})
	}

	return eg.Wait()
}

func getItem(client *resty.Client, account model.Account, networkID constants.NetworkID) error {
	// Get the timestamp of the latest data to avoid duplicate pulls whenever possible.
	var timestamp time.Time

	if err := database.DB.
		Raw(
			`WITH "timestamp" AS (SELECT date_created AS "timestamp"
                     FROM note3
                     WHERE owner = ?
                     ORDER BY date_created DESC
                     LIMIT 1)
SELECT COALESCE("timestamp".timestamp, TIMESTAMPTZ 'epoch') AS "timestamp"
FROM "timestamp";`,
			rss3uri.New(rss3uri.NewAccountInstance(account.Identity, constants.PlatformID(account.ProfilePlatform).Symbol())).String(),
		).
		Scan(&timestamp).
		Error; err != nil {
		return err
	}

	request := client.NewRequest()
	params := map[string]string{
		"proof":             strings.ToLower(account.Identity),
		"platform_id":       strconv.Itoa(account.Platform),
		"network_id":        strconv.Itoa(int(networkID)),
		"profile_source_id": strconv.Itoa(account.Source),
		"owner_id":          strings.ToLower(account.ProfileID),
		"owner_platform_id": strconv.Itoa(account.ProfilePlatform),
		"timestamp":         big.NewInt(timestamp.Unix()).String(),
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

	if response.StatusCode() != http.StatusOK || result.Error.Code != 0 {
		logger.Error(response.StatusCode(), result.Error.Code, result.Error.Msg)

		return nil
	}

	return nil
}
