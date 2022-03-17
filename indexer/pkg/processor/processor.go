package processor

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
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
		Broker:          config.Config.Broker.Addr,
		DefaultQueue:    "indexer_queue",
		ResultBackend:   config.Config.Broker.Addr,
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
	tasks := map[string]interface{}{
		"dispatch": Dispatch,
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
