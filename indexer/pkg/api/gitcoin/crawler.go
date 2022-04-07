package gitcoin

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/moralis"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/xscan"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/zksync"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/datatype"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/httpx"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/valyala/fastjson"
)

type crawler struct {
	eth     crawlerConfig
	polygon crawlerConfig
	zk      crawlerConfig

	ZksTokensCache       map[int64]zksync.Token
	inactiveAdminsCache  map[string]bool
	hostingProjectsCache map[string]ProjectInfo
}

func NewCrawler(ethParam, polygonParam, zkParam crawlerConfig) *crawler {
	return &crawler{
		ethParam,
		polygonParam,
		zkParam,
		make(map[int64]zksync.Token),
		make(map[string]bool),
		make(map[string]ProjectInfo),
	}
}

func (gc *crawler) InitZksTokenCache() error {
	tokens, err := zksync.GetTokens()
	if err != nil {
		logger.Errorf("zksync get tokens error: %v", err)

		return err
	}

	for _, token := range tokens {
		gc.ZksTokensCache[token.Id] = token
	}

	return nil
}

func (gc *crawler) InitGrants() error {
	grants, err := GetGrantsInfo() // get grant project list metadata
	if err != nil {
		return err
	}

	for _, item := range grants {
		if item.AdminAddress != "0x0" {
			gc.updateHostingProject(item.AdminAddress) // get grant project detailed info

			time.Sleep(10 * time.Second)
		}
	}

	return nil
}

// UpdateZksToken update Token by tokenId
func (gc *crawler) UpdateZksToken() error {
	tokens, err := zksync.GetTokens()
	if err != nil {
		logger.Errorf("zksync get tokens error: %v", err)

		return err
	}

	for _, token := range tokens {
		gc.ZksTokensCache[token.Id] = token
	}

	return nil
}

// GetZksToken returns Token by tokenId
func (gc *crawler) GetZksToken(id int64) zksync.Token {
	return gc.ZksTokensCache[id]
}

// inactiveAdminAddress checks an admin address is active or not
func (gc *crawler) inactiveAdminAddress(adminAddress string) bool {
	adminAddress = strings.ToLower(adminAddress)

	return gc.inactiveAdminsCache[adminAddress]
}

// addInactiveAdminAddress adds an admin address
func (gc *crawler) addInactiveAdminAddress(adminAddress string) {
	adminAddress = strings.ToLower(adminAddress)
	gc.inactiveAdminsCache[adminAddress] = true
}

func (gc *crawler) hostingProject(adminAddress string) (ProjectInfo, bool) {
	p, ok := gc.hostingProjectsCache[adminAddress]

	return p, ok
}

func (gc *crawler) needUpdateProject(adminAddress string) bool {
	if len(gc.hostingProjectsCache) == 0 {
		return true
	}

	p, ok := gc.hostingProject(adminAddress)

	return ok && !p.Active
}

func (gc *crawler) updateHostingProject(adminAddress string) (inactive bool, err error) {

	project, err := GetProjectsInfo(adminAddress, "")
	if err != nil {
		logger.Errorf("zksync get projects info error: %v", err)

		return
	}

	if !project.Active {
		gc.addInactiveAdminAddress(adminAddress)
	}

	gc.hostingProjectsCache[adminAddress] = project // TODO: add to db
	inactive = !project.Active

	return
}

// GetProjectsInfo returns project info from gitcoin
func GetProjectsInfo(adminAddress string, title string) (ProjectInfo, error) {
	var project ProjectInfo

	headers := make(map[string]string)
	httpx.SetCommonHeader(headers)

	url := grantsApi + "?admin_address=" + adminAddress

	maxRetries := 3
	content, err := httpx.Get(url, headers)

	for i := 1; i <= maxRetries; i++ {
		if err == nil {
			break
		}

		content, err = httpx.Get(url, headers)

		logger.Warnf("GetProjectsInfo error [%v], times: [%d]", err, i)
		time.Sleep(10 * time.Second)
	}

	if err != nil {
		logger.Errorf("gitcoin get project info error: [%v]", err)

		return project, err
	}

	// check reCAPTCHA
	if strings.Contains(string(content), "Hold up, the bots want to know if you're one of them") {
		err = fmt.Errorf("gitcoin get project info error, reCAPTCHA")
		return project, err
	}

	logger.Infof("GetProjectsInfo success, url: [%s]", url)

	var parser fastjson.Parser
	parsedJson, parseErr := parser.Parse(string(content))

	if parseErr != nil {
		logger.Errorf("gitcoin parse json error: [%v]", parseErr)

		return project, parseErr
	}

	if "[]" == string(content) {
		// project is inactive
		project.Active = false
		project.AdminAddress = adminAddress
		project.Title = title
	} else {
		// project is active
		project.Active = true
		project.AdminAddress = adminAddress
		project.Title = string(parsedJson.GetStringBytes("title"))
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

func (gc *crawler) zksyncRun() error {
	if gc.zk.NextRoundTime.After(time.Now()) {
		return nil
	}

	latestConfirmedBlockHeight, err := zksync.GetLatestBlockHeightWithConfirmations(gc.zk.Confirmations)
	if err != nil {
		logger.Errorf("zksync get latest block error: %v", err)

		return err
	}

	// scan the latest block content periodically
	endBlockHeight := gc.zk.FromHeight + gc.zk.Step
	if latestConfirmedBlockHeight < endBlockHeight {
		gc.zk.NextRoundTime = gc.zk.NextRoundTime.Add(gc.zk.SleepInterval)
		// use minStep when catching up with the latest block height
		gc.zk.Step = gc.zk.MinStep

		logger.Info("zksync catch up with the latest block height")

		return nil
	}

	logger.Infof("get zksync donations, from [%d] to [%d]", gc.zk.FromHeight, endBlockHeight)

	// get zksync donations
	donations, err := gc.GetZkSyncDonations(gc.zk.FromHeight, endBlockHeight)
	if err != nil {
		logger.Errorf("zksync get donations error: %v", err)

		return err
	}

	if len(donations) > 0 {
		setDB(donations, constants.NetworkIDZksync)
	}

	// set new fromHeight
	gc.zk.FromHeight = endBlockHeight

	return nil
}

func (gc *crawler) getConfig(networkId constants.NetworkID) *crawlerConfig {

	if networkId == constants.NetworkIDEthereum {
		return &gc.eth
	}
	if networkId == constants.NetworkIDPolygon {
		return &gc.polygon
	}
	logger.Errorf("unsupported network")

	return nil
}

func getDonationPlatform(networkId constants.NetworkID) GitcoinPlatform {
	if networkId == constants.NetworkIDEthereum {
		return ETH
	}
	if networkId == constants.NetworkIDPolygon {
		return Polygon
	}
	logger.Errorf("unsupported network")
	return ""
}

func (gc *crawler) xscanRun(networkId constants.NetworkID) error {
	donationPlatform := getDonationPlatform(networkId)
	p := gc.getConfig(networkId)

	if p.NextRoundTime.After(time.Now()) {
		return nil
	}

	latestConfirmedBlockHeight, err := xscan.GetLatestBlockHeightWithConfirmations(networkId, p.Confirmations)
	if err != nil {
		logger.Errorf("[%s] get latest block error: %v", networkId.Symbol(), err)

		return err
	}

	endBlockHeight := p.FromHeight + p.Step
	if latestConfirmedBlockHeight < endBlockHeight {
		p.NextRoundTime = p.NextRoundTime.Add(p.SleepInterval)
		// use minStep when catching up with the latest block height
		p.Step = p.MinStep

		logger.Infof("gitcoin [%s] catch up with the latest block height", networkId.Symbol())

		return nil
	}

	logger.Infof("get [%s] donations, from [%d] to [%d]", networkId.Symbol(), p.FromHeight, endBlockHeight)

	donations, err := GetEthDonations(p.FromHeight, endBlockHeight, donationPlatform)
	if err != nil {
		logger.Errorf("[%s] get donations error: %v", networkId.Symbol(), err)

		return err
	}

	if len(donations) > 0 {
		setDB(donations, networkId)
	}

	// set new fromHeight
	p.FromHeight = endBlockHeight

	return nil
}

func setDB(donations []DonationInfo, networkId constants.NetworkID) error {
	//logger.Infof("set db, network: [%s]", networkId.Symbol())
	items := make([]model.Note, 0)

	for _, v := range donations {
		author := rss3uri.NewAccountInstance(v.Donor, constants.PlatformSymbolEthereum).UriString()

		tsp, err := time.Parse(time.RFC3339, v.Timestamp)
		if err != nil {
			logger.Errorf("gitcoin parse time error: %v", err)

			tsp = time.Now()
		}
		// TODO: read from db to get project info
		// if not in db, ok is false
		ok := true
		if !ok {
			GetProjectsInfo(v.AdminAddress, "")
		}
		attachment := datatype.Attachments{
			{
				Type:     "title",
				Content:  "", //TODO: Read from db
				MimeType: "text/plain",
			},
			{
				Type:     "description",
				Content:  "", //TODO: Read from db
				MimeType: "text/plain",
			},
			{
				Type:        "logo",
				Content:     "", //TODO: Read from db
				MimeType:    "", // TODO
				SizeInBytes: 0,  //TODO
			},
		}
		note := model.Note{
			Identifier: rss3uri.NewNoteInstance(author, networkId.Symbol()).UriString(),
			Owner:      author,
			RelatedURLs: []string{
				moralis.GetTxHashURL(networkId.Symbol(), v.TxHash),
				"https://gitcoin.co/grants/2679/rss3-rss-with-human-curation", //TODO: read from db
			},
			Tags:            constants.ItemTagsDonationGitcoin.ToPqStringArray(),
			Authors:         []string{author},
			Title:           "",
			Summary:         "",
			Attachments:     database.MustWrapJSON(attachment),
			Source:          constants.NoteSourceNameGitcoinContribution.String(),
			MetadataNetwork: constants.NetworkSymbolEthereum.String(),
			MetadataProof:   v.TxHash,
			Metadata: database.MustWrapJSON(map[string]interface{}{
				"from": v.Donor,
				"to":   v.GetTxTo(),

				"destination":  v.AdminAddress,
				"value_amount": v.FormatedAmount.String(),
				"value_symbol": v.Symbol,
				"approach":     v.Approach,
			}),
			DateCreated: tsp,
			DateUpdated: tsp,
		}

		items = append(items, note)
	}

	// TODO: make insert db a general method @Zerber
	tx := database.DB.Begin()
	defer tx.Rollback()

	if items != nil && len(items) > 0 {
		if _, dbErr := database.CreateNotes(tx, items, true); dbErr != nil {
			return dbErr
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (gc *crawler) ZkStart() error {
	signal.Notify(gc.zk.Interrupt, os.Interrupt)

	for {
		select {
		case <-gc.zk.Interrupt:
			logger.Info("ZkStart gets interrupt signal")

			return nil
		default:
			gc.zksyncRun()
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func (gc *crawler) EthStart() error {
	signal.Notify(gc.eth.Interrupt, os.Interrupt)

	for {
		select {
		case <-gc.eth.Interrupt:
			logger.Info("EthStart gets interrupt signal")

			return nil
		default:
			gc.xscanRun(constants.NetworkIDEthereum)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func (gc *crawler) PolygonStart() error {
	signal.Notify(gc.polygon.Interrupt, os.Interrupt)

	for {
		select {
		case <-gc.polygon.Interrupt:
			logger.Info("PolygonStart gets interrupt signal")

			return nil
		default:
			gc.xscanRun(constants.NetworkIDPolygon)
			time.Sleep(500 * time.Millisecond)
		}
	}
}
