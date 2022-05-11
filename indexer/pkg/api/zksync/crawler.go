package zksync

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
)

func Start() error {
	if err := zkCP.start(); err != nil {
		return fmt.Errorf("start crawler error: %v", err)
	}

	return nil
}

func (crawler *crawler) start() error {
	if err := UpdateZksToken(); err != nil {
		return fmt.Errorf("update zks token error: %v", err)
	}

	height, err := util.GetCrawlerMetadata(
		crawler.metadataIdentity, crawler.platformID)
	if err != nil {
		logger.Warnf("get last height error: %v", err)
	} else {
		crawler.config.FromHeight = height
	}

	crawler.loopRun()

	return nil
}

func (crawler *crawler) loopRun() {
	for {
		crawler.run()

		// Since the interval time of each time may change dynamically,
		// it is necessary to read the interval time of the next round from config.SleepInterval
		config := crawler.getConfig()

		sleepInterval := config.SleepInterval
		if crawler.getConfig().SleepInterval <= 0 {
			sleepInterval = DeafultGetNextBlockDuration
		}

		time.Sleep(sleepInterval)

		config.SleepInterval = DeafultGetNextBlockDuration
	}
}

func (crawler *crawler) run() error {
	if err := crawler.checkConfig(); err != nil {
		return fmt.Errorf("zksync crawler run error: %s", err)
	}

	config := crawler.config

	latestConfirmedBlockHeight, err := GetLatestBlockHeightWithConfirmations(config.Confirmations)
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
		config.SleepInterval = GetLatestNextBlockDuration

		// use minStep when catching up with the latest block height
		config.Step = config.MinStep

		return nil
	}

	// get and format zksync metadata
	zksyncInfos, err := crawler.formatZkSyncMetadata(config.FromHeight, endBlockHeight)
	if err != nil {
		logger.Errorf("zksync get info error: %v", err)

		return err
	}

	if len(zksyncInfos) > 0 {
		err := crawler.setDB(zksyncInfos, constants.NetworkIDEthereum)
		if err != nil {
			logger.Errorf("set db error: %v", err)

			return err
		}
	}

	logger.Infof("zksync: from [%d] to [%d], the latest confirmed block height [%d]",
		config.FromHeight, endBlockHeight, latestConfirmedBlockHeight)

	// set new fromHeight
	config.FromHeight = endBlockHeight + 1

	if err := util.SetCrawlerMetadata(crawler.metadataIdentity, config.FromHeight,
		crawler.platformID); err != nil {
		logger.Errorf("set crawler metadata error: %v", err)

		return err
	}

	return nil
}

func (crawler *crawler) getConfig() *crawlerConfig {
	return crawler.config
}

func (crawler *crawler) checkConfig() error {
	if crawler.config.FromHeight < 0 {
		return fmt.Errorf("invalid from height: %d", crawler.config.FromHeight)
	}

	if crawler.config.Step <= 0 ||
		crawler.config.MinStep <= 0 {
		return fmt.Errorf("invalid step: %d, minStep: %d", crawler.config.Step, crawler.config.MinStep)
	}

	if crawler.config.Confirmations <= 0 {
		return fmt.Errorf("invalid confirmations: %d", crawler.config.Confirmations)
	}

	if crawler.config.SleepInterval <= 0 {
		return fmt.Errorf("invalid sleep interval: %d", crawler.config.SleepInterval)
	}

	return nil
}

func (crawler *crawler) formatZkSyncMetadata(fromBlock int64, toBlock int64) ([]*ZkSyncInfo, error) {
	zksyncInfos := []*ZkSyncInfo{}

	for i := fromBlock; i <= toBlock; i++ {
		trxs, err := GetTxsByBlock(i)
		if err != nil {
			logger.Errorf("get txs by block error: [%v]", err)

			return nil, err
		}

		for _, tx := range trxs {
			// to address empty
			to := strings.ToLower(tx.Op.To)
			if to == "" ||
				to == "0x0" ||
				to == "0x0000000000000000000000000000000000000000" {
				continue
			}

			tokenId := tx.Op.TokenId
			token := ZksTokensCache[tokenId]
			formatedAmount := big.NewInt(1)
			formatedAmount.SetString(tx.Op.Amount, 10)

			zksyncInfos = append(zksyncInfos, &ZkSyncInfo{
				From:           tx.Op.From,
				To:             to,
				TokenAddress:   token.Address,
				Amount:         tx.Op.Amount,
				Symbol:         token.Symbol,
				FormatedAmount: formatedAmount,
				Decimals:       token.Decimals,
				Timestamp:      tx.CreatedAt,
				TxHash:         tx.TxHash,
				Type:           tx.Op.Type,
			})
		}
	}

	return zksyncInfos, nil
}

func (crawler *crawler) setDB(zksyncInfo []*ZkSyncInfo, networkID constants.NetworkID) error {
	items := make([]model.Note, 0)

	if len(zksyncInfo) <= 0 {
		return nil
	}

	niBuilder := crawler.getNewNoteInstanceBuilder()

	for _, zksync := range zksyncInfo {
		tsp, err := time.Parse(time.RFC3339, zksync.Timestamp)
		if err != nil {
			tsp = time.Now()
		}

		note, err := crawler.setNote(zksync, networkID, tsp, niBuilder)
		if err != nil {
			continue
		}

		items = append(items, *note)
	}

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

func (crawler *crawler) setNote(zksyncInfo *ZkSyncInfo, networkID constants.NetworkID, tsp time.Time,
	niBuilder *noteInstanceBuilder) (*model.Note, error) {
	author := rss3uri.NewAccountInstance(zksyncInfo.From, constants.PlatformSymbolEthereum).UriString()
	instanceKey, err := crawler.setNoteInstance(niBuilder, zksyncInfo.TxHash)

	if err != nil {
		return nil, fmt.Errorf("set note instance error: %s", err)
	}

	note := model.Note{
		Identifier: rss3uri.NewNoteInstance(instanceKey, networkID.Symbol()).UriString(),
		RelatedURLs: []string{
			fmt.Sprintf("https://zkscan.io/explorer/transactions/%v", zksyncInfo.TxHash),
		},
		Owner:           author,
		Tags:            constants.ItemTagsToken.ToPqStringArray(),
		Authors:         []string{author},
		ContractAddress: zksyncInfo.TokenAddress,
		Source:          constants.NoteSourceNameEthereumERC20.String(),
		MetadataNetwork: constants.NetworkSymbolZkSync.String(),
		MetadataProof:   instanceKey,
		Metadata: database.MustWrapJSON(map[string]interface{}{
			"from":          zksyncInfo.From,
			"to":            zksyncInfo.To,
			"token_address": zksyncInfo.TokenAddress,
			"value_amount":  zksyncInfo.FormatedAmount.String(),
			"value_symbol":  zksyncInfo.Symbol,
			"tx_hash":       zksyncInfo.TxHash,
			"type":          zksyncInfo.Type,
			"decimals":      zksyncInfo.Decimals,
		}),
		DateCreated: tsp,
		DateUpdated: tsp,
	}

	if strings.Contains(strings.ToUpper(zksyncInfo.Type), "NFT") {
		note.Tags = constants.ItemTagsNFT.ToPqStringArray()
	}

	return &note, nil
}

func (crawler *crawler) getNewNoteInstanceBuilder() *noteInstanceBuilder {
	return &noteInstanceBuilder{
		countMap: map[string]int{},
	}
}

func (crawler *crawler) setNoteInstance(niBuilder *noteInstanceBuilder, txHash string) (string, error) {
	if niBuilder == nil || len(txHash) == 0 {
		return "", fmt.Errorf("instance build error")
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
