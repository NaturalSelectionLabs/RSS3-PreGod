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
)

type Param struct {
	FromHeight    int64
	Step          int64
	MinStep       int64
	Confirmations int64
	SleepInterval time.Duration
	Interrupt     chan os.Signal
}

func NewParam(from, step, minStep, confirmations, sleepInterval int64) Param {
	return Param{
		FromHeight:    from,
		Step:          step,
		MinStep:       minStep,
		Confirmations: confirmations,
		SleepInterval: time.Duration(sleepInterval),
		Interrupt:     make(chan os.Signal, 1),
	}
}

type gitcoinCrawler struct {
	eth     Param
	polygon Param
	zk      Param

	zksTokensCache       map[int64]zksync.Token
	inactiveAdminsCache  map[string]bool
	hostingProjectsCache map[string]ProjectInfo
}

func NewGitcoinCrawler(ethParam, polygonParam, zkParam Param) *gitcoinCrawler {
	return &gitcoinCrawler{
		ethParam,
		polygonParam,
		zkParam,
		make(map[int64]zksync.Token),
		make(map[string]bool),
		make(map[string]ProjectInfo),
	}
}

// UpdateZksToken update Token by tokenId
func (gc *gitcoinCrawler) UpdateZksToken() error {
	tokens, err := zksync.GetTokens()
	if err != nil {
		return err
	}

	for _, token := range tokens {
		gc.zksTokensCache[token.Id] = token
	}

	return nil
}

// GetZksToken returns Token by tokenId
func (gc *gitcoinCrawler) GetZksToken(id int64) zksync.Token {
	return gc.zksTokensCache[id]
}

// inactiveAdminAddress checks an admin address is active or not
func (gc *gitcoinCrawler) inactiveAdminAddress(adminAddress string) bool {
	adminAddress = strings.ToLower(adminAddress)

	return gc.inactiveAdminsCache[adminAddress]
}

// addInactiveAdminAddress adds an admin address
func (gc *gitcoinCrawler) addInactiveAdminAddress(adminAddress string) {
	adminAddress = strings.ToLower(adminAddress)
	gc.inactiveAdminsCache[adminAddress] = true
}

func (gc *gitcoinCrawler) hostingProject(adminAddress string) (ProjectInfo, bool) {
	p, ok := gc.hostingProjectsCache[adminAddress]

	return p, ok
}

func (gc *gitcoinCrawler) needUpdateProject(adminAddress string) bool {
	p, ok := gc.hostingProject(adminAddress)

	return !(ok && p.Active)
}

func (gc *gitcoinCrawler) updateHostingProject(adminAddress string) (inactive bool, err error) {
	project, err := GetProjectsInfo(adminAddress, "")
	if err != nil {
		return
	}

	if !project.Active {
		gc.addInactiveAdminAddress(adminAddress)
	}

	gc.hostingProjectsCache[adminAddress] = project
	inactive = !project.Active

	return
}

func (gc *gitcoinCrawler) zksyncRun() error {
	// token cache
	if len(gc.zksTokensCache) == 0 {
		tokens, err := zksync.GetTokens()
		if err != nil {
			return err
		}

		for _, token := range tokens {
			gc.zksTokensCache[token.Id] = token
		}
	}

	latestBlockHeight, err := zksync.GetLatestBlockHeight()
	if err != nil {
		return err
	}

	latestConfirmedBlockHeight := latestBlockHeight - gc.zk.Confirmations

	// scan the latest block content periodically
	endBlockHeight := gc.zk.FromHeight + gc.zk.Step
	if latestConfirmedBlockHeight <= endBlockHeight {
		time.Sleep(gc.zk.SleepInterval)

		latestBlockHeight, err = zksync.GetLatestBlockHeight()
		if err != nil {
			return err
		}

		endBlockHeight = latestBlockHeight - gc.zk.Confirmations
		gc.zk.Step = gc.zk.MinStep
	}

	// get zksync donations
	donations, err := gc.GetZkSyncDonations(gc.zk.FromHeight, endBlockHeight)
	if err != nil {
		return err
	}

	setDB(donations, constants.NetworkIDZksync)

	// set new from height
	gc.zk.FromHeight = endBlockHeight

	return nil
}

func (gc *gitcoinCrawler) xscanRun(networkId constants.NetworkID) error {
	var p *Param
	if networkId == constants.NetworkIDEthereumMainnet {
		p = &gc.eth
	} else if networkId == constants.NetworkIDPolygon {
		p = &gc.polygon
	}

	latestBlockHeight, err := xscan.GetLatestBlockHeight(networkId)
	if err != nil {
		return err
	}

	endBlockHeight := p.FromHeight + p.Step
	if latestBlockHeight <= endBlockHeight {
		time.Sleep(p.SleepInterval)

		latestBlockHeight, err = xscan.GetLatestBlockHeight(networkId)
		if err != nil {
			return err
		}

		endBlockHeight = latestBlockHeight
		p.Step = p.MinStep
	}

	var chainType ChainType
	if networkId == constants.NetworkIDEthereumMainnet {
		chainType = ETH
	} else if networkId == constants.NetworkIDPolygon {
		chainType = Polygon
	}

	donations, err := GetEthDonations(p.FromHeight, endBlockHeight, chainType)
	if err != nil {
		return err
	}

	setDB(donations, networkId)

	// set new from height
	p.FromHeight = endBlockHeight

	return nil
}

func setDB(donations []DonationInfo, networkId constants.NetworkID) {
	items := []*model.Item{}

	for _, v := range donations {
		tsp, err := time.Parse(time.RFC3339, v.Timestamp)
		if err != nil {
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
			nil,
			nil,
			"",
			"",
			[]model.Attachment{},
			tsp,
		)

		items = append(items, item)
	}

	db.InsertItems(items, networkId)
}

func (gc *gitcoinCrawler) ZkStart() error {
	signal.Notify(gc.zk.Interrupt, os.Interrupt)

	for {
		select {
		case <-gc.zk.Interrupt:
			return nil
		default:
			gc.zksyncRun()
		}
	}
}

func (gc *gitcoinCrawler) EthStart() error {
	signal.Notify(gc.eth.Interrupt, os.Interrupt)

	for {
		select {
		case <-gc.eth.Interrupt:
			return nil
		default:
			gc.xscanRun(constants.NetworkIDEthereumMainnet)
		}
	}
}

func (gc *gitcoinCrawler) PolygonStart() error {
	signal.Notify(gc.polygon.Interrupt, os.Interrupt)

	for {
		select {
		case <-gc.polygon.Interrupt:
			return nil
		default:
			gc.xscanRun(constants.NetworkIDPolygon)
		}
	}
}
