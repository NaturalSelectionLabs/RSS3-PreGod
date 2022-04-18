package es

import (
	"context"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/elastic/go-elasticsearch/v7"
)

var (
	client *elasticsearch.Client
)

func init() {
	if err := Setup(); err != nil {
		panic(err)
	}
}

func Setup() (err error) {
	client, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{config.Config.Elasticsearch.Address},
		Username:  config.Config.Elasticsearch.Username,
		Password:  config.Config.Elasticsearch.Password,
	})
	if err != nil {
		return err
	}

	_, err = client.Ping()

	return err
}

func Ping(ctx context.Context) error {
	_, err := client.Ping()

	return err
}
