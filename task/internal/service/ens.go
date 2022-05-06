package service

import (
	"context"
	"fmt"
	"math/big"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
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

func (s *Service) SubscribeEns() {
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

				// save ens data into db 
				// todo 

				// trigger task: get owner feed 
				
			
			}
		}
	}
}
