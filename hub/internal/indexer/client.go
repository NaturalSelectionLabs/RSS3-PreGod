package indexer

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/go-resty/resty/v2"
	"golang.org/x/sync/errgroup"
)

func GetItems(accounts []model.Account) error {
	eg := errgroup.Group{}
	client := resty.New()

	for _, account := range accounts {
		account := account
		eg.Go(func() error {
			request := client.NewRequest()
			params := map[string]string{
				"proof":       account.ID,
				"platform_id": strconv.Itoa(account.Platform),
				"network_id":  strconv.Itoa(int(constants.NetworkSymbol(constants.PlatformID(account.Platform).Symbol()).ID())),
			}
			request.SetQueryParams(params)
			// TODO request.SetResult()
			response, err := request.Get(EndpointItem)
			if err != nil {
				return err
			}

			if response.StatusCode() != http.StatusOK /* TODO || response struct code != 0 */ {
				return errors.New("indexer error")
			}

			return nil
		})
	}

	return eg.Wait()
}
