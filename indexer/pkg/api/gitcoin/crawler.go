package gitcoin

import (
	"strings"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/zksync"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
)

type gitcoinCrawler struct {
	crawler.CrawlerResult
	zksTokensCache       map[int64]zksync.Token
	inactiveAdminsCache  map[string]bool
	hostingProjectsCache map[string]ProjectInfo
}

func NewGitcoinCrawler() gitcoinCrawler {
	return gitcoinCrawler{
		crawler.CrawlerResult{
			Assets: []*model.ItemId{},
			Notes:  []*model.ItemId{},
			Items:  []*model.Item{},
		},
		make(map[int64]zksync.Token),
		make(map[string]bool),
		make(map[string]ProjectInfo),
	}
}

func NewCrawler() crawler.Crawler {
	return &gitcoinCrawler{
		crawler.CrawlerResult{
			Assets: []*model.ItemId{},
			Notes:  []*model.ItemId{},
			Items:  []*model.Item{},
		},
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

func (gc *gitcoinCrawler) Work(owner string, network constants.NetworkID) error {
	tokens, err := zksync.GetTokens()

	if err != nil {
		return err
	}

	for _, token := range tokens {
		gc.zksTokensCache[token.Id] = token
	}

	return nil
}

func (gc *gitcoinCrawler) GetResult() *crawler.CrawlerResult {
	return &crawler.CrawlerResult{
		Assets: gc.Assets,
		Notes:  gc.Notes,
		Items:  gc.Items,
	}
}
