package processor

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/jike"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/misskey"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/moralis"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/twitter"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/backends/result"
	machineryConfig "github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
)

type Processor struct {
	server  *machinery.Server
	workers []*machinery.Worker
}

var (
	processor Processor
)

func Setup() error {
	cnf := &machineryConfig.Config{
		Broker:          config.Config.RabbitMQ.Addr,
		DefaultQueue:    "indexer_queue",
		ResultBackend:   config.Config.RabbitMQ.Addr,
		ResultsExpireIn: 3600,
		AMQP: &machineryConfig.AMQPConfig{
			Exchange:      "indexer_exchange",
			ExchangeType:  "direct",
			BindingKey:    "indexer_task",
			PrefetchCount: 3,
		},
	}

	server, err := machinery.NewServer(cnf)
	if err != nil {
		return err
	}

	processor.server = server

	// Register tasks
	tasks := make(map[string]interface{})

	for _, networkId := range constants.NetworkIDMap {
		var crawler interface{}

		switch networkId {
		case constants.NetworkIDEthereumMainnet,
			constants.NetworkIDBNBChain,
			constants.NetworkIDAvalanche,
			constants.NetworkIDFantom,
			constants.NetworkIDPolygon:
			crawler = moralis.Crawl
		case constants.NetworkIDMisskey:
			crawler = misskey.Crawl
		case constants.NetworkIDJike:
			crawler = jike.Crawl
		case constants.NetworkIDTwitter:
			crawler = twitter.Crawl
		default:
			crawler = nil
		}

		if crawler != nil {
			tasks[string(networkId)] = crawler
		}
	}

	return server.RegisterTasks(tasks)
}

func NewWorker(queueName string, consumerName string, concurrency int) error {
	worker := processor.server.NewWorker(consumerName, concurrency)
	worker.Queue = queueName
	processor.workers = append(processor.workers, worker)

	return worker.Launch()
}

func SendTask(task tasks.Signature) (*result.AsyncResult, error) {
	asyncResult, err := processor.server.SendTask(&task)
	if err != nil {
		return nil, err
	}

	return asyncResult, nil
}
