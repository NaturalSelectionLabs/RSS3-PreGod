package gitcoin

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/moralis"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/xscan"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/zksync"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/datatype"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
)

type crawlerPropertyInf interface {
	run() error
	start() error
	getConfig() crawlerConfig
	getPlatform() GitcoinPlatform
}

type crawlerProperty struct {
	config           *crawlerConfig
	platform         GitcoinPlatform
	networkID        constants.NetworkID
	platformID       constants.PlatformID
	metadataIdentity string
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
			config:           DefaultZksyncConfig,
			platform:         ZKSYNC,
			networkID:        constants.NetworkIDEthereum,
			platformID:       constants.PlatformID(1002),
			metadataIdentity: string("gitcoin-" + ZKSYNC),
		},
	}

	ethCP = xscanRunCrawlerProperty{
		crawlerProperty{
			config:           DefaultEthConfig,
			platform:         ETH,
			networkID:        constants.NetworkIDEthereum,
			platformID:       constants.PlatformID(1003),
			metadataIdentity: string("gitcoin-" + ETH),
		},
	}

	PolygonCP = xscanRunCrawlerProperty{
		crawlerProperty{
			config:           DefaultPolygonConfig,
			platform:         Polygon,
			networkID:        constants.NetworkIDPolygon,
			platformID:       constants.PlatformID(1004),
			metadataIdentity: string("gitcoin-" + Polygon),
		},
	}

	crawlerPropertyMap = map[GitcoinPlatform]crawlerPropertyInf{
		ZKSYNC:  &zkCP,
		ETH:     &ethCP,
		Polygon: &PolygonCP,
	}
)

func loopRun(property crawlerPropertyInf) error {
	config := property.getConfig()
	signal.Notify(config.Interrupt, os.Interrupt)

	for {
		property.run()
		time.Sleep(500 * time.Millisecond)
	}
}

func (property *crawlerProperty) getConfig() crawlerConfig {
	return *property.config
}

func (property *crawlerProperty) getPlatform() GitcoinPlatform {
	return property.platform
}

func (property *crawlerProperty) configCheck() error {
	if property.config.FromHeight < 0 {
		return fmt.Errorf("invalid from height: %d", property.config.FromHeight)
	}

	if property.config.Step <= 0 ||
		property.config.MinStep <= 0 {
		return fmt.Errorf("invalid step: %d, minStep: %d", property.config.Step, property.config.MinStep)
	}

	if property.config.Confirmations <= 0 {
		return fmt.Errorf("invalid confirmations: %d", property.config.Confirmations)
	}

	if property.config.SleepInterval <= 0 {
		return fmt.Errorf("invalid sleep interval: %d", property.config.SleepInterval)
	}

	return nil
}

func (property *zksyncCrawlerProperty) start() error {
	if err := UpdateZksToken(); err != nil {
		return fmt.Errorf("update zks token error: %v", err)
	}

	height, err := util.GetCrawlerMetadata(
		property.metadataIdentity, property.platformID)
	if err != nil {
		logger.Warnf("get last height error: %v", err)
	} else {
		property.config.FromHeight = height
	}

	if err := loopRun(property); err != nil {
		return fmt.Errorf("zksync run error: %s", err)
	}

	return nil
}

func (property *zksyncCrawlerProperty) run() error {
	if err := property.configCheck(); err != nil {
		return fmt.Errorf("zksync crawler run error: %s", err)
	}

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
	endBlockHeight := config.FromHeight + config.Step - 1
	if endBlockHeight <= 0 {
		logger.Fatalf("config.FromHeight [%d] + config.Step [%d] - 1 <= 0", config.FromHeight, config.Step)
	}

	if latestConfirmedBlockHeight < endBlockHeight {
		config.NextRoundTime = config.NextRoundTime.Add(config.SleepInterval)
		// use minStep when catching up with the latest block height
		config.Step = config.MinStep

		logger.Debugf("zksync catch up with the latest block height, latestConfirmedBlockHeight[%d], endBlockHeight[%d]",
			latestConfirmedBlockHeight, endBlockHeight)

		return nil
	}

	// get zksync donations
	donations, adminAddresses, err := GetZkSyncDonations(config.FromHeight, endBlockHeight)
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
	config.FromHeight = endBlockHeight + 1

	if err := util.SetCrawlerMetadata(
		property.metadataIdentity, config.FromHeight, property.platformID); err != nil {
		logger.Errorf("set crawler metadata error: %v", err)

		return err
	}

	return nil
}

func (property *xscanRunCrawlerProperty) start() error {
	if err := UpdateEthAndPolygonTokens(); err != nil {
		return fmt.Errorf("xscan run error: %v", err)
	}

	height, err := util.GetCrawlerMetadata(
		property.metadataIdentity, property.platformID)
	if err != nil {
		logger.Warnf("get last height error: %v", err)
	} else {
		property.config.FromHeight = height
	}

	if err := loopRun(property); err != nil {
		return fmt.Errorf("xscan run error: %s", err)
	}

	return nil
}

func (property *xscanRunCrawlerProperty) run() error {
	// donationPlatform := getDonationPlatform(networkId)
	config := property.config

	if config.NextRoundTime.After(time.Now()) {
		return nil
	}

	latestConfirmedBlockHeight, err := xscan.GetLatestBlockHeightWithConfirmations(property.networkID, config.Confirmations)
	if err != nil {
		logger.Errorf("[%s] get latest block error: %v", property.networkID.Symbol(), err)

		return err
	}

	endBlockHeight := config.FromHeight + config.Step - 1
	if latestConfirmedBlockHeight < endBlockHeight {
		config.NextRoundTime = config.NextRoundTime.Add(config.SleepInterval)
		// use minStep when catching up with the latest block height
		config.Step = config.MinStep

		logger.Infof("gitcoin [%s] catch up with the latest block height", property.networkID.Symbol())

		return nil
	}

	logger.Infof("get [%s] donations, from [%d] to [%d]", property.networkID.Symbol(), config.FromHeight, endBlockHeight)

	donations, adminAddresses, err := GetEthDonations(config.FromHeight, endBlockHeight, property.platform)
	if err != nil {
		logger.Errorf("[%s] get donations error: %v", property.networkID.Symbol(), err)

		return err
	}

	logger.Infof("len(donations): %d, len(adminAddresses): %d", len(donations), len(adminAddresses))

	if len(donations) > 0 {
		setDB(donations, property.networkID, adminAddresses)
	}

	// set new fromHeight
	config.FromHeight = endBlockHeight + 1

	if err := util.SetCrawlerMetadata(
		property.metadataIdentity, config.FromHeight, property.platformID); err != nil {
		logger.Errorf("set crawler metadata error: %v", err)

		return err
	}

	return nil
}

// Since the txhash transaction that pulls eth and polygon may have batch transactions in the same txhash,
// it is necessary to use an array suffix to mark different transactions of
// the same txhash when storing in the database as the primary key
type noteInstanceBuilder struct {
	countMap map[string]int
}

func getNewNoteInstanceBuilder() *noteInstanceBuilder {
	return &noteInstanceBuilder{
		countMap: map[string]int{},
	}
}

func setNoteInstance(
	niBuilder *noteInstanceBuilder,
	txHash string) (string, error) {
	if niBuilder == nil {
		return "", fmt.Errorf("note instance builder is nil")
	}

	if txHash == "" {
		return "", fmt.Errorf("tx hash is empty")
	}

	hashCount, ok := niBuilder.countMap[txHash]
	if !ok {
		niBuilder.countMap[txHash] = 0

		return txHash + "-0", nil
	}

	hashCount += 1

	niBuilder.countMap[txHash] = hashCount

	return txHash + "-" + strconv.Itoa(hashCount), nil
}

func setNote(
	donationInfo *DonationInfo,
	networkId constants.NetworkID,
	projectInfo *ProjectInfo,
	v *DonationInfo,
	tsp time.Time,
	niBuilder *noteInstanceBuilder) (*model.Note, error) {
	if projectInfo == nil || v == nil {
		return nil, fmt.Errorf("invalid projectInfo or donationInfo")
	}

	if niBuilder == nil {
		return nil, fmt.Errorf("note instance builder is nil")
	}

	author := rss3uri.NewAccountInstance(donationInfo.Donor, constants.PlatformSymbolEthereum).UriString()
	summary := util.EllipsisContent(projectInfo.Description, 400)
	instanceKey, err := setNoteInstance(niBuilder, donationInfo.TxHash)

	if err != nil {
		return nil, fmt.Errorf("set note instance error: %s", err)
	}

	note := model.Note{
		Identifier: rss3uri.NewNoteInstance(instanceKey, networkId.Symbol()).UriString(),
		Owner:      author,
		RelatedURLs: []string{
			moralis.GetTxHashURL(networkId.Symbol(), v.TxHash),
			"https://gitcoin.co/grants/2679/rss3-rss-with-human-curation", //TODO: read from db
		},
		Tags:    constants.ItemTagsDonationGitcoin.ToPqStringArray(),
		Authors: []string{author},
		Title:   projectInfo.Title,
		Summary: summary,
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
	// logger.Infof("note:%v", note)
	// time.Sleep(2 * time.Second)

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

	niBuilder := getNewNoteInstanceBuilder()

	for _, v := range donations {
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

		note, err := setNote(&v, networkId, &projectInfo, &v, tsp, niBuilder)
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

	if err := property.start(); err != nil {
		return fmt.Errorf("start crawler error: %v", err)
	}

	return nil
}
