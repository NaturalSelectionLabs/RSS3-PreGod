package gitcoin

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/moralis"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/types"
	"github.com/valyala/fastjson"
)

const grantUrl = "https://gitcoin.co/grants/grants.json"
const grantsApi = "https://gitcoin.co/api/v0.1/grants/"
const donationSentTopic = "0x3bb7428b25f9bdad9bd2faa4c6a7a9e5d5882657e96c1d24cc41c1d6c1910a98"
const bulkCheckoutAddress = "0x7d655c57f71464B6f83811C55D84009Cd9f5221C"

type (
	GrantInfo    = types.GrantInfo
	ProjectInfo  = types.ProjectInfo
	DonationInfo = types.DonationInfo
)

// GetGrants returns all grant projects.
func GetGrants() (content []byte, err error) {
	content, err = util.Get(grantUrl, nil)

	return
}

func GetProject(adminAddress string) (content []byte, err error) {
	url := grantsApi + "?admin_address=" + adminAddress
	content, err = util.Get(url, nil)

	return
}

func GetGrantsInfo() ([]GrantInfo, error) {
	content, err := GetGrants()
	if err != nil {
		return nil, err
	}

	var parser fastjson.Parser
	parsedJson, parseErr := parser.Parse(string(content))

	if parseErr != nil {
		return nil, nil
	}

	grantArrs := parsedJson.GetArray()
	grants := make([]GrantInfo, len(grantArrs))

	for _, grant := range grantArrs {
		projects := grant.GetArray()

		item := GrantInfo{Title: projects[0].String(), AdminAddress: projects[1].String()}
		grants = append(grants, item)
	}

	return grants, nil
}

func GetProjectsInfo(adminAddress string, title string) (ProjectInfo, error) {
	var project ProjectInfo

	content, err := GetProject(adminAddress)
	if err != nil {
		return project, err
	}

	var parser fastjson.Parser
	parsedJson, parseErr := parser.Parse(string(content))

	if parseErr != nil {
		return project, nil
	}

	if "[]" == string(content) {
		// project is inactive
		project.Active = false
		project.AdminAddress = adminAddress
		project.Title = title
	} else {
		project.Active = true
		project.AdminAddress = adminAddress
		project.Title = title
		project.Id = parsedJson.GetInt64("id")
		project.Slug = string(parsedJson.GetStringBytes("slug"))
		project.Description = string(parsedJson.GetStringBytes("description"))
		project.ReferUrl = string(parsedJson.GetStringBytes("reference_url"))
		project.Logo = string(parsedJson.GetStringBytes("logo"))
		project.TokenAddress = string(parsedJson.GetStringBytes("token_address"))
		project.TokenSymbol = string(parsedJson.GetStringBytes("token_symbol"))
		project.ContractAddress = string(parsedJson.GetStringBytes("contract_address"))
	}

	return project, nil
}

func GetDonations(fromBlock int64, toBlock int64) ([]DonationInfo, error) {
	chainType := "eth"
	apiKey := "" // TODO, read api key from config
	logs, err := moralis.GetLogs(fromBlock, toBlock, bulkCheckoutAddress, donationSentTopic, chainType, apiKey)

	if err != nil {
		return nil, err
	}

	donations := make([]DonationInfo, len(logs.Result))

	for _, item := range logs.Result {
		donor := "0x" + item.Topic3[26:]
		donation := DonationInfo{
			Donor: donor,
		}

		donations = append(donations, donation)
	}

	return donations, nil
}