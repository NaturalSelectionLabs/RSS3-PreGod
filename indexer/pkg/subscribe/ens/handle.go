package ens

import (
	"embed"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

type Ens struct {
	EthClient *ethclient.Client
	ABI       abi.ABI
	Database  *gorm.DB
}

var (
	//go:embed event.abi
	abiFileSystem embed.FS
)

func Run() {
	var err error

	var s = &Ens{
		Database: database.DB,
	}

	// get ethclient
	s.EthClient, err = ethclient.Dial(config.Config.Indexer.Gateway.Endpoint)
	if err != nil {
		logger.Errorf("task: ethclient Dial error, %v", err)

		return 
	}

	// get abi
	abiFile, err := abiFileSystem.Open("event.abi")

	if err != nil {
		logger.Errorf("task: open abi file error, %v", err)

		return 
	}

	s.ABI, err = abi.JSON(abiFile)
	if err != nil {
		logger.Errorf("task: abi file parse error, %v", err)

		return 
	}

	s.SubscribeEns()

}
