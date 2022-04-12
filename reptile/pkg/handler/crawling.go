package handler

import (
	"fmt"
	"strconv"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/gitcoin"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/httpx"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/valyala/fastjson"
)

const gitcoinPrefix = "https://gitcoin.co/grants/v1/api/grant/"

func GetResult(pos int) (*gitcoin.ProjectInfo, error) {
	info := new(gitcoin.ProjectInfo)
	url := gitcoinPrefix + strconv.Itoa(pos)
	// logger.Infof("url: %s", url)
	result, err := httpx.Get(url, nil)

	if err != nil {
		return nil, fmt.Errorf("get result false:%s", err)
	}

	if result == nil {
		return nil, nil
	}

	var parser fastjson.Parser
	parsedJson, parseErr := parser.Parse(string(result))

	if parseErr != nil {
		return nil, fmt.Errorf("get result false:%s", parseErr)
	}

	if parsedJson.GetInt("status") != 200 {
		logger.Errorf("get result [%d] not find: %s", pos,
			string(parsedJson.GetStringBytes("message")))

		return nil, nil
	}

	grantsObjectJson := parsedJson.Get("grants")

	info.Active = grantsObjectJson.GetBool("active")
	info.Id = grantsObjectJson.GetInt64("id")
	info.Title = string(grantsObjectJson.GetStringBytes("title"))
	info.Slug = string(grantsObjectJson.GetStringBytes("slug"))
	info.Description = string(grantsObjectJson.GetStringBytes("description"))
	info.ReferUrl = string(grantsObjectJson.GetStringBytes("reference_url"))
	info.Logo = string(grantsObjectJson.GetStringBytes("logo_url"))
	info.AdminAddress = string(grantsObjectJson.GetStringBytes("admin_address"))
	info.TokenAddress = string(grantsObjectJson.GetStringBytes("token_address"))
	info.TokenSymbol = string(grantsObjectJson.GetStringBytes("token_symbol"))
	info.ContractAddress = string(grantsObjectJson.GetStringBytes("contract_address"))

	// logger.Infof("%v", info)

	return info, nil
}
