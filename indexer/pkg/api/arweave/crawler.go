package arweave

import (
	"errors"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
)

var ErrTimeout = errors.New("received timeout")
var ErrInterrupt = errors.New("received interrupt")

type arCrawler struct {
	fromHeight    int64
	confirmations int64
	step          int64
	minStep       int64
	sleepInterval time.Duration
	identity      string
	interrupt     chan os.Signal
	complete      chan error
}

func NewArCrawler(fromHeight, step, minStep, confirmatioins, sleepInterval int64, identity string) *arCrawler {
	return &arCrawler{
		fromHeight,
		confirmatioins,
		step,
		minStep,
		time.Duration(sleepInterval),
		identity,
		make(chan os.Signal, 1),
		make(chan error),
	}
}

func (ar *arCrawler) run() error {
	startBlockHeight := ar.fromHeight
	step := ar.step
	tempDelay := ar.sleepInterval

	// get latest block height
	latestBlockHeight, err := GetLatestBlockHeight()
	if err != nil {
		return err
	}

	latestConfirmedBlockHeight := latestBlockHeight - ar.confirmations

	for {
		// handle interrupt
		if ar.gotInterrupt() {
			return ErrInterrupt
		}

		// get articles
		endBlockHeight := startBlockHeight + step
		if latestConfirmedBlockHeight <= endBlockHeight {
			time.Sleep(tempDelay)

			latestBlockHeight, err = GetLatestBlockHeight()
			if err != nil {
				return err
			}

			latestConfirmedBlockHeight = latestBlockHeight - ar.confirmations
			step = 10
		}

		ar.getArticles(startBlockHeight, latestConfirmedBlockHeight, ar.identity)
	}
}

func (ar *arCrawler) getArticles(from, to int64, owner string) error {
	articles, err := GetArticles(from, to, owner)
	if err != nil {
		return err
	}

	items := make([]*model.Item, 0)

	for _, article := range articles {
		attachment := model.Attachment{
			Type:     "body",
			Content:  article.Content,
			MimeType: "text/markdown",
		}

		tsp, err := time.Parse(time.RFC3339, strconv.FormatInt(article.TimeStamp, 10))
		if err != nil {
			logger.Error(err)

			tsp = time.Now()
		}

		ni := model.NewItem(
			constants.NetworkSymbolArweaveMainnet.GetID(),
			article.Digest,
			model.Metadata{
				"network": constants.NetworkSymbolArweaveMainnet,
				"proof":   article.Digest,
			},
			constants.ItemTagsMirrorEntry,
			[]string{article.Author},
			article.Title,
			article.Content, // TODO: According to RIP4, if the body is too long, then only record part of the body, followed by ... at the end
			[]model.Attachment{attachment},
			tsp,
		)

		items = append(items, ni)
	}

	setDB(items)

	return nil
}

func setDB(items []*model.Item) {
	for _, item := range items {
		db.InsertItem(item)
	}
}

func (ar *arCrawler) Start() error {
	signal.Notify(ar.interrupt, os.Interrupt)

	go func() {
		ar.complete <- ar.run()
	}()

	select {
	case err := <-ar.complete:
		return err
	default:
		return nil
	}
}

func (ar *arCrawler) gotInterrupt() bool {
	select {
	case <-ar.interrupt:
		signal.Stop(ar.interrupt)

		return true
	default:
		return false
	}
}
