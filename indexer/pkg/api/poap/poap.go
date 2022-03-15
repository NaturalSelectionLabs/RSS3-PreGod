package poap

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/httpx"
	jsoniter "github.com/json-iterator/go"
)

const endpoint = "https://api.poap.xyz"

var jsoni = jsoniter.ConfigCompatibleWithStandardLibrary

func GetActions(user string) ([]PoapResponse, error) {
	if len(user) == 0 {
		var err = fmt.Errorf("user address is empty")

		return []PoapResponse{}, err
	}

	url := fmt.Sprintf("%s/actions/scan/%s",
		endpoint, user)
	response, err := httpx.Get(url, nil)

	if err != nil {
		return []PoapResponse{}, err
	}

	res := new([]PoapResponse)

	err = jsoni.Unmarshal(response, &res)
	if err != nil {
		return []PoapResponse{}, err
	}

	return *res, nil
}
