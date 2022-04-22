package crossbell

import (
	"context"
	"database/sql"
	"errors"
	"sync/atomic"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type Crawler interface {
	Initialize() error
	Run() error
}

var (
	ErrInvalidConfig        = errors.New("invalid config")
	ErrNotConnectedDatabase = errors.New("not connected to database")
)

var _ Crawler = &crawler{}

type crawler struct {
	config         *Config
	ethereumClient *ethclient.Client
	db             *gorm.DB
	headerCh       chan *types.Header

	latestBlockNumber   int64
	internalBlockNumber int64
}

func (c *crawler) Initialize() (err error) {
	if c.config == nil {
		return ErrInvalidConfig
	}

	c.ethereumClient, err = ethclient.Dial(c.config.RPC)
	if err != nil {
		return err
	}

	logger.Info("connected to crossbell rpc")

	c.headerCh = make(chan *types.Header)

	if database.DB == nil {
		return ErrNotConnectedDatabase
	}

	c.db = database.DB

	return nil
}

func (c *crawler) Run() error {
	if err := c.Initialize(); err != nil {
		return err
	}

	subscription, err := c.ethereumClient.SubscribeNewHead(context.Background(), c.headerCh)
	if err != nil {
		return err
	}

	defer subscription.Unsubscribe()

	logger.Info("subscribe new head success")

	eg := errgroup.Group{}

	eg.Go(c.runHandler)
	eg.Go(c.runSubscriber)

	return eg.Wait()
}

func (c *crawler) runHandler() error {
	if err := c.db.
		Table("crawler_metadata").
		Select("last_block").
		Where(&model.CrawlerMetadata{
			AccountInstance: ContractAddressProfile,
			PlatformID:      0,
		}).
		Scan(&c.internalBlockNumber).Error; err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	// TODO

	return nil
}

func (c *crawler) runSubscriber() error {
	blockNumber, err := c.ethereumClient.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	c.latestBlockNumber = int64(blockNumber)

	for header := range c.headerCh {
		atomic.StoreInt64(&c.latestBlockNumber, header.Number.Int64())
	}

	return nil
}

func New(config *Config) Crawler {
	return &crawler{
		config: config,
	}
}
