package crossbell

import (
	"errors"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Crawler interface {
	Initialize() error
	Run() error
}

var (
	ErrInvalidConfig = errors.New("invalid config")
)

var _ Crawler = &crawler{}

type crawler struct {
	config         *Config
	ethereumClient *ethclient.Client
}

func (c *crawler) Initialize() (err error) {
	if c.config == nil {
		return ErrInvalidConfig
	}

	c.ethereumClient, err = ethclient.Dial(c.config.RPC)
	if err != nil {
		return err
	}

	return nil
}

func (c *crawler) Run() error {
	if err := c.Initialize(); err != nil {
		return err
	}

	// TODO

	return nil
}

func New(config *Config) Crawler {
	return &crawler{
		config: config,
	}
}
