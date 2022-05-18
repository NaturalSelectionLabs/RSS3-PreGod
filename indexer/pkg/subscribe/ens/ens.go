package ens

import (
	"context"
	"math/big"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler_handler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type NameRegisteredData struct {
	Name    string
	Cost    *big.Int
	Expires *big.Int
}

var TopicHashNameRegistered = common.HexToHash("0xca6abbe9d7f11422cb6ca7629fbf6fe9efb1c621f71ce8f02b9f2a230097404f")

func (s *Ens) SubscribeEns() {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{
			common.HexToAddress("0x283Af0B28c62C092C9727F1Ee09c02CA627EB7F5"),
		},
	}

	logs := make(chan types.Log)
	sub, err := s.EthClient.SubscribeFilterLogs(context.Background(), query, logs)

	if err != nil {
		logger.Errorf("SubscribeEns: ethclient SubscribeFilterLogs error, %v", err)

		return
	}

	for {
		select {
		case err := <-sub.Err():
			logger.Errorf("SubscribeEns: ethclient subscribe error, %v", err)
		case vLog := <-logs:
			if vLog.Topics[0] == TopicHashNameRegistered {
				var data = NameRegisteredData{}

				// parse contract log
				if err := s.ABI.UnpackIntoInterface(&data, "NameRegistered", vLog.Data); err != nil {
					logger.Errorf("SubscribeEns: parse data into NameRegistered error, %v", err)

					continue
				}

				// get owner by topics
				owner := common.HexToAddress(vLog.Topics[2].Hex())

				// get block details
				block, err := s.EthClient.BlockByNumber(context.Background(), big.NewInt(int64(vLog.BlockNumber)))
				if err != nil {
					logger.Errorf("SubscribeEns: get block error, %v", err)

					continue
				}

				// save ens data into db
				ens := &model.Domains{
					TransactionHash: vLog.TxHash.Bytes(),
					Type:            "ens",
					Name:            data.Name,
					AddressOwner:    owner.Bytes(),
					ExpiredAt:       time.Unix(data.Expires.Int64(), 0),
					Source:          "subscribe",
					BlockTimestamp:  time.Unix(int64(block.Time()), 0),
				}
				err = s.CreateEns(ens)
				if err != nil {
					logger.Errorf("SubscribeEns: db insert error, %v", err)

					continue
				}

				// trigger task: get owner feed
				go func() {
					s.GetOwnerFeed(owner.String())
				}()
			}
		}
	}
}

func (s *Ens) GetOwnerFeed(owner string) {
	instance, err := rss3uri.NewInstance("account", owner, "ethereum")
	if err != nil {
		return
	}

	networkIDs := constants.GetEthereumPlatformNetworks()
	for _, networkID := range networkIDs {
		getItemHandler := crawler_handler.NewGetItemsHandler(crawler.WorkParam{
			Identity:   instance.GetIdentity(),
			PlatformID: constants.PlatformIDEthereum,
			NetworkID:  networkID,
			OwnerID:    owner,
		})

		_, err = getItemHandler.Excute()
		if err != nil {
			logger.Errorf("SubscribeEns: get item error, %v", err)

			continue
		}
	}
}

func (s *Ens) CreateEns(ens *model.Domains) error {
	if err := s.Database.Create(ens).Error; err != nil {
		return err
	}
	return nil
}
