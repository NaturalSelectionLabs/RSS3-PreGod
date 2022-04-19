package gitcoin

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/moralis"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/xscan"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/zksync"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/datatype"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
)

type crawlerPropertyInf interface {
	run() error
	getConfig() crawlerConfig
}

type crawlerProperty struct {
	config    *crawlerConfig
	platform  GitcoinPlatform
	networkID constants.NetworkID
}

type zksyncCrawlerProperty struct {
	crawlerProperty
}

type xscanRunCrawlerProperty struct {
	crawlerProperty
}

var (
	zkCP = zksyncCrawlerProperty{
		crawlerProperty{
			config:    DefaultZksyncConfig,
			platform:  ZKSYNC,
			networkID: constants.NetworkIDEthereum,
		},
	}

	ethCP = xscanRunCrawlerProperty{
		crawlerProperty{
			config:    DefaultEthConfig,
			platform:  ETH,
			networkID: constants.NetworkIDEthereum,
		},
	}

	PolygonCP = xscanRunCrawlerProperty{
		crawlerProperty{
			config:    DefaultPolygonConfig,
			platform:  Polygon,
			networkID: constants.NetworkIDPolygon,
		},
	}

	crawlerPropertyMap = map[GitcoinPlatform]crawlerPropertyInf{
		ZKSYNC:  zkCP,
		ETH:     ethCP,
		Polygon: PolygonCP,
	}
)

func (property crawlerProperty) getConfig() crawlerConfig {
	return *property.config
}

func (property zksyncCrawlerProperty) run() error {
	config := property.config

	if config.NextRoundTime.After(time.Now()) {
		return nil
	}

	latestConfirmedBlockHeight, err := zksync.GetLatestBlockHeightWithConfirmations(config.Confirmations)
	if err != nil {
		logger.Errorf("zksync get latest block error: %v", err)

		return err
	}

	// scan the latest block content periodically
	endBlockHeight := config.FromHeight + config.Step
	if latestConfirmedBlockHeight < endBlockHeight {
		config.NextRoundTime = config.NextRoundTime.Add(config.SleepInterval)
		// use minStep when catching up with the latest block height
		config.Step = config.MinStep

		logger.Infof("zksync catch up with the latest block height, latestConfirmedBlockHeight[%d], endBlockHeight[%d]",
			latestConfirmedBlockHeight, endBlockHeight)

		return nil
	}

	logger.Infof("get zksync donations, from [%d] to [%d]", config.FromHeight, endBlockHeight)

	// get zksync donations
	donations, adminAddresses, err := GetZkSyncDonations(config.FromHeight, endBlockHeight)
	// GetZkSyncDonations(config.FromHeight, endBlockHeight)
	if err != nil {
		logger.Errorf("zksync get donations error: %v", err)

		return err
	}

	if len(donations) > 0 {
		err := setDB(donations, constants.NetworkIDEthereum, adminAddresses)
		if err != nil {
			logger.Errorf("set db error: %v", err)

			return err
		}
	}

	// set new fromHeight
	config.FromHeight = endBlockHeight
	logger.Infof("config.FromHeight: %d", config.FromHeight)

	return nil
}

func (property xscanRunCrawlerProperty) run() error {
	// donationPlatform := getDonationPlatform(networkId)
	p := property.config

	if p.NextRoundTime.After(time.Now()) {
		return nil
	}

	latestConfirmedBlockHeight, err := xscan.GetLatestBlockHeightWithConfirmations(property.networkID, p.Confirmations)
	if err != nil {
		logger.Errorf("[%s] get latest block error: %v", property.networkID.Symbol(), err)

		return err
	}

	endBlockHeight := p.FromHeight + p.Step
	if latestConfirmedBlockHeight < endBlockHeight {
		p.NextRoundTime = p.NextRoundTime.Add(p.SleepInterval)
		// use minStep when catching up with the latest block height
		p.Step = p.MinStep

		logger.Infof("gitcoin [%s] catch up with the latest block height", property.networkID.Symbol())

		return nil
	}

	logger.Infof("get [%s] donations, from [%d] to [%d]", property.networkID.Symbol(), p.FromHeight, endBlockHeight)

	donations, adminAddresses, err := GetEthDonations(p.FromHeight, endBlockHeight, property.platform)
	if err != nil {
		logger.Errorf("[%s] get donations error: %v", property.networkID.Symbol(), err)

		return err
	}

	if len(donations) > 0 {
		setDB(donations, property.networkID, adminAddresses)
	}

	// set new fromHeight
	p.FromHeight = endBlockHeight

	return nil
}

func setNote(
	author string,
	networkId constants.NetworkID,
	projectInfo *ProjectInfo,
	v *DonationInfo,
	tsp time.Time) (*model.Note, error) {
	if projectInfo == nil || v == nil {
		return nil, fmt.Errorf("invalid projectInfo or donationInfo")
	}

	note := model.Note{
		Identifier: rss3uri.NewNoteInstance(author, networkId.Symbol()).UriString(),
		Owner:      author,
		RelatedURLs: []string{
			moralis.GetTxHashURL(networkId.Symbol(), v.TxHash),
			"https://gitcoin.co/grants/2679/rss3-rss-with-human-curation", //TODO: read from db
		},
		Tags:    constants.ItemTagsDonationGitcoin.ToPqStringArray(),
		Authors: []string{author},
		Title:   "",
		Summary: "",
		Attachments: database.MustWrapJSON(datatype.Attachments{
			{
				Type:     "title",
				Content:  projectInfo.Title,
				MimeType: "text/plain",
			},
			{
				Type:     "description",
				Content:  projectInfo.Description,
				MimeType: "text/plain",
			},
			{
				Type:        "logo",
				Content:     projectInfo.Logo,
				MimeType:    "image/png",
				SizeInBytes: 0,
			},
		}),
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

	return &note, nil
}

func setDB(donations []DonationInfo,
	networkId constants.NetworkID,
	adminAddresses []string) error {
	items := make([]model.Note, 0)

	if len(donations) <= 0 {
		return nil
	}

	logger.Infof("len(adminAddresses): %d", len(adminAddresses))

	// get all project infos from db
	projects, err := GetProjectsInfo(adminAddresses)
	if err != nil {
		return fmt.Errorf("get projects error: %v", err)
	}

	logger.Infof("len(projects): %d", len(projects))

	if len(projects) <= 0 {
		return nil
	}

	logger.Infof("%d", len(donations))

	for _, v := range donations {
		author := rss3uri.NewAccountInstance(v.Donor, constants.PlatformSymbolEthereum).UriString()

		tsp, err := time.Parse(time.RFC3339, v.Timestamp)
		if err != nil {
			logger.Errorf("gitcoin parse time error: %v", err)

			tsp = time.Now()
		}

		// TODO: here will be add cache to reduce db interview time
		projectInfo, ok := projects[v.AdminAddress]
		if !ok {
			continue
		}

		note, err := setNote(author, networkId, &projectInfo, &v, tsp)
		if err != nil {
			logger.Errorf("gitcoin set note error: %v", err)

			continue
		}

		items = append(items, *note)
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

func GitCoinStart(platform GitcoinPlatform) error {
	property, ok := crawlerPropertyMap[platform]
	if !ok {
		return fmt.Errorf("invalid network id: %s", platform)
	}

	config := property.getConfig()
	signal.Notify(config.Interrupt, os.Interrupt)

	for {
		select {
		case <-config.Interrupt:
			logger.Infof("%s start gets interrupt signal", platform)

			return nil
		default:
			property.run()
			time.Sleep(500 * time.Millisecond)
		}
	}
}
