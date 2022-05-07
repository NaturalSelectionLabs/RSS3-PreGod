package task

import (
	"context"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/ethclient"
)

func SubscribeEns() {
	client, err := ethclient.Dial()
	if err != nil {
		logger.Errorf("SubscribeEns: ethclient Dial error, %v", err)
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{
			common.HexToAddress("0x283Af0B28c62C092C9727F1Ee09c02CA627EB7F5")
		}
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		logger.Errorf("SubscribeEns: ethclient SubscribeFilterLogs error, %v", err)
	}

	for {
		select {
		case err := <-sub.Err():
			logger.Errorf("SubscribeEns: ethclient subscribe error, %v", err)
		case vLog := <-logs:
			fmt.Println(vLog) // pointer to event log
		}
	}
}
