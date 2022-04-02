package gitcoin

import (
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/xscan"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/zksync"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
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

	zksTokensCache       map[int64]zksync.Token
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

func (gc *crawler) initGrants() error {
	grants, err := GetGrantsInfo()
	if err != nil {
		return err
	}

	for _, item := range grants {
		if item.AdminAddress != "0x0" {
			gc.updateHostingProject(item.AdminAddress)

			time.Sleep(5 * time.Second)
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
		gc.zksTokensCache[token.Id] = token
	}

	return nil
}

// GetZksToken returns Token by tokenId
func (gc *crawler) GetZksToken(id int64) zksync.Token {
	return gc.zksTokensCache[id]
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

	gc.hostingProjectsCache[adminAddress] = project
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
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		logger.Errorf("gitcoin get project info error: [%v]", err)

		return project, err
	}

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

func (gc *crawler) zksyncRun() error {
	logger.Info("Starting run zksync")

	// token cache
	if len(gc.zksTokensCache) == 0 {
		tokens, err := zksync.GetTokens()
		if err != nil {
			logger.Errorf("zksync get tokens error: %v", err)

			return err
		}

		for _, token := range tokens {
			gc.zksTokensCache[token.Id] = token
		}
	}

	// get grants
	if len(gc.hostingProjectsCache) == 0 {
		gc.initGrants()
	}

	latestConfirmedBlockHeight, err := zksync.GetLatestBlockHeightWithConfirmations(gc.zk.Confirmations)
	if err != nil {
		logger.Errorf("zksync get latest block error: %v", err)

		return err
	}

	// scan the latest block content periodically
	endBlockHeight := gc.zk.FromHeight + gc.zk.Step
	if latestConfirmedBlockHeight < endBlockHeight {
		time.Sleep(gc.zk.SleepInterval)

		latestConfirmedBlockHeight, err = zksync.GetLatestBlockHeightWithConfirmations(gc.zk.Confirmations)
		if err != nil {
			logger.Errorf("zksync get latest block error: %v", err)

			return err
		}

		if latestConfirmedBlockHeight < endBlockHeight {
			return nil
		}

		// use minStep when catching up with the latest block height
		gc.zk.Step = gc.zk.MinStep
	}

	// get zksync donations
	donations, err := gc.GetZkSyncDonations(gc.zk.FromHeight, endBlockHeight)
	if err != nil {
		logger.Errorf("zksync get donations error: %v", err)

		return err
	}

	if len(donations) > 0 {
		setDB(donations, constants.NetworkIDZksync)
	}

	// set new from height
	gc.zk.FromHeight = endBlockHeight

	return nil
}

func (gc *crawler) xscanRun(networkId constants.NetworkID) error {
	var p *crawlerConfig
	if networkId == constants.NetworkIDEthereum {
		p = &gc.eth
	} else if networkId == constants.NetworkIDPolygon {
		p = &gc.polygon
	}

	latestConfirmedBlockHeight, err := xscan.GetLatestBlockHeightWithConfirmations(networkId, p.Confirmations)
	if err != nil {
		logger.Errorf("xscan get latest block error: %v", err)

		return err
	}

	endBlockHeight := p.FromHeight + p.Step
	if latestConfirmedBlockHeight < endBlockHeight {
		for {
			time.Sleep(p.SleepInterval)

			latestConfirmedBlockHeight, err = xscan.GetLatestBlockHeightWithConfirmations(networkId, p.Confirmations)
			if err != nil {
				logger.Errorf("xscan get latest block error: %v", err)

				return err
			}

			if latestConfirmedBlockHeight > endBlockHeight {
				break
			}
		}

		// use minStep when catching up with the latest block height
		p.Step = p.MinStep
	}

	var chainType ChainType
	if networkId == constants.NetworkIDEthereum {
		chainType = ETH
	} else if networkId == constants.NetworkIDPolygon {
		chainType = Polygon
	}

	donations, err := GetEthDonations(p.FromHeight, endBlockHeight, chainType)
	if err != nil {
		logger.Errorf("[%s] get donations error: %v", networkId.Symbol(), err)

		return err
	}

	setDB(donations, networkId)

	// set new from height
	p.FromHeight = endBlockHeight

	return nil
}

func setDB(donations []DonationInfo, networkId constants.NetworkID) {
	logger.Infof("set db, network: [%s]", networkId.Symbol())

	items := make([]*model.Item, 0)

	for _, v := range donations {
		instance := rss3uri.NewAccountInstance(v.Donor, constants.PlatformSymbolEthereum)
		author, err := rss3uri.NewInstance("account", v.Donor, string(constants.PlatformSymbolEthereum))

		if err != nil {
			logger.Errorf("gitcoin setDB get new instance error: %v", err)

			return
		}

		tsp, err := time.Parse(time.RFC3339, v.Timestamp)
		if err != nil {
			logger.Errorf("gitcoin parse time error: %v", err)

			tsp = time.Now()
		}

		item := model.NewItem(
			networkId,
			v.TxHash,
			model.Metadata{
				"Donor":            v.Donor,
				"AdminAddress":     v.AdminAddress,
				"TokenAddress":     v.TokenAddress,
				"Symbol":           v.Symbol,
				"Amount":           v.FormatedAmount,
				"DonationApproach": v.Approach,
			},
			constants.ItemTagsDonationGitcoin,
			[]string{author.UriString()},
			"",
			"",
			[]model.Attachment{},
			tsp,
		)
		items = append(items, item)

		// append notes
		notes := []*model.ObjectId{{
			NetworkID: networkId,
			Proof:     v.TxHash,
		}}
		db.AppendNotes(instance, notes)
	}

	db.InsertItems(items, networkId)
}

func (gc *crawler) ZkStart() error {
	logger.Info("Start crawling gitcoin zksync")
	signal.Notify(gc.zk.Interrupt, os.Interrupt)

	for {
		select {
		case <-gc.zk.Interrupt:
			return nil
		default:
			gc.zksyncRun()
			time.Sleep(5 * time.Second)
		}
	}
}

func (gc *crawler) EthStart() error {
	logger.Info("Start crawling gitcoin eth")
	signal.Notify(gc.eth.Interrupt, os.Interrupt)

	for {
		select {
		case <-gc.eth.Interrupt:
			return nil
		default:
			gc.xscanRun(constants.NetworkIDEthereum)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func (gc *crawler) PolygonStart() error {
	logger.Info("Start crawling gitcoin polygon")
	signal.Notify(gc.polygon.Interrupt, os.Interrupt)

	for {
		select {
		case <-gc.polygon.Interrupt:
			return nil
		default:
			gc.xscanRun(constants.NetworkIDPolygon)
			time.Sleep(500 * time.Millisecond)
		}
	}
}
