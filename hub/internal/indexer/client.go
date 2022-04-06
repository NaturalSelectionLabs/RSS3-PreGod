package indexer

import (
	"errors"
	"fmt"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/go-resty/resty/v2"
	"golang.org/x/sync/errgroup"
	"net/http"
)

func GetItems(profiles []model.Profile) error {
	eg := errgroup.Group{}
	client := resty.New()

	for _, profile := range profiles {
		eg.Go(func() error {
			fmt.Println(profile)
			request := client.NewRequest()
			request.SetQueryParams(map[string]string{
				"proof":       "KallyDev",
				"platform_id": "6",
				"network_id":  "12",
			})
			//request.SetResult()
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
