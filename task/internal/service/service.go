package service

import (
	"embed"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Service struct {
	EthClient *ethclient.Client
	ABI       abi.ABI
}

var (
	//go:embed event.abi
	abiFileSystem embed.FS
)

func NewService() *Service {
	var err error

	var s = &Service{}

	// get ethclient
	s.EthClient, err = ethclient.Dial(config.Config.Indexer.Gateway.Endpoint)
	if err != nil {
		logger.Errorf("task: ethclient Dial error, %v", err)

		return nil
	}

	// get abi
	abiFile, err := abiFileSystem.Open("event.abi")

	if err != nil {
		logger.Errorf("task: open abi file error, %v", err)

		return nil
	}

	s.ABI, err = abi.JSON(abiFile)
	if err != nil {
		logger.Errorf("task: abi file parse error, %v", err)

		return nil
	}

	return s
}
