package indexer

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/api"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/go-resty/resty/v2"
	"golang.org/x/sync/errgroup"
)

func GetItems(accounts []model.Account) error {
	eg := errgroup.Group{}
	client := resty.New()

	for _, account := range accounts {
		account := account
		eg.Go(func() error {
			start := time.Now()
			request := client.NewRequest()
			params := map[string]string{
				"proof":       account.ID,
				"platform_id": strconv.Itoa(account.Platform),
				"network_id":  strconv.Itoa(int(constants.NetworkSymbol(constants.PlatformID(account.Platform).Symbol()).ID())),
			}
			result := Response{}
			response, err := request.
				SetQueryParams(params).
				SetResult(&result).
				Get(fmt.Sprintf("%s/item", config.Config.Hub.IndexerEndpoint))
			if err != nil {
				logger.Error(err)

				return api.ErrorIndexer
			}

			if response.StatusCode() != http.StatusOK || result.Error.Code != 0 {
				return api.ErrorIndexer
			}

			logger.Info(account.ID, account.Platform, start, time.Now().Sub(start))
			return nil
		})
	}

	return eg.Wait()
}
