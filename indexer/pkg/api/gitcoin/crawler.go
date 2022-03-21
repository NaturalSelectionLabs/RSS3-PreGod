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
	Complete      chan error
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
	startBlockHeight := gc.zk.FromHeight
	step := gc.zk.Step
	tempDelay := gc.zk.SleepInterval

	// token cache
	tokens, err := zksync.GetTokens()
	if err != nil {
		return err
	}

	for _, token := range tokens {
		gc.zksTokensCache[token.Id] = token
	}

	latestBlockHeight, err := zksync.GetLatestBlockHeight()
	if err != nil {
		return err
	}

	latestConfirmedBlockHeight := latestBlockHeight - gc.zk.Confirmations

	// scan the latest block content periodically
	for {
		endBlockHeight := startBlockHeight + step
		if latestConfirmedBlockHeight <= endBlockHeight {
			time.Sleep(tempDelay)

			latestBlockHeight, err = zksync.GetLatestBlockHeight()
			if err != nil {
				return err
			}

			endBlockHeight = latestBlockHeight - gc.zk.Confirmations
			step = gc.zk.MinStep
		}

		// get zksync donations
		donations, err := gc.GetZkSyncDonations(startBlockHeight, endBlockHeight)
		if err != nil {
			return err
		}

		setDB(donations, constants.NetworkIDZksync)
	}
}

func (gc *gitcoinCrawler) xscanWork(networkId constants.NetworkID) error {
	startBlockHeight := int64(1)

	var p Param
	if networkId == constants.NetworkIDEthereumMainnet {
		p = gc.eth
	} else if networkId == constants.NetworkIDPolygon {
		p = gc.polygon
	}

	step := p.Step
	minStep := p.MinStep
	sleepInterval := p.SleepInterval

	for {
		latestBlockHeight, err := xscan.GetLatestBlockHeight(networkId)
		if err != nil {
			return err
		}

		endBlockHeight := startBlockHeight + step
		if latestBlockHeight <= endBlockHeight {
			time.Sleep(sleepInterval)

			latestBlockHeight, err = xscan.GetLatestBlockHeight(networkId)
			if err != nil {
				return err
			}

			endBlockHeight = latestBlockHeight
			step = minStep
		}

		var chainType ChainType
		if networkId == constants.NetworkIDEthereumMainnet {
			chainType = ETH
		} else if networkId == constants.NetworkIDPolygon {
			chainType = Polygon
		}

		donations, err := GetEthDonations(startBlockHeight, endBlockHeight, chainType)
		if err != nil {
			return err
		}

		setDB(donations, networkId)
	}
}

func setDB(donations []DonationInfo, networkId constants.NetworkID) {
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
		db.InsertItem(item)
	}
}

func (gc *gitcoinCrawler) ZkStart() error {
	signal.Notify(gc.zk.Interrupt, os.Interrupt)

	go func() {
		gc.zk.Complete <- gc.zksyncRun()
	}()

	select {
	case err := <-gc.zk.Complete:
		return err
	default:
		return nil
	}
}

func (gc *gitcoinCrawler) EthStart() error {
	signal.Notify(gc.eth.Interrupt, os.Interrupt)

	go func() {
		gc.eth.Complete <- gc.xscanWork(constants.NetworkIDEthereumMainnet)
	}()

	select {
	case err := <-gc.eth.Complete:
		return err
	default:
		return nil
	}
}

func (gc *gitcoinCrawler) PolygonStart() error {
	signal.Notify(gc.polygon.Interrupt, os.Interrupt)

	go func() {
		gc.polygon.Complete <- gc.xscanWork(constants.NetworkIDPolygon)
	}()

	select {
	case err := <-gc.polygon.Complete:
		return err
	default:
		return nil
	}
}
