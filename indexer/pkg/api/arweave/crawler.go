package arweave

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
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

func NewArCrawler(fromHeight, step, minStep, confirmations, sleepInterval int64, identity string) *arCrawler {
	return &arCrawler{
		fromHeight,
		confirmations,
		step,
		minStep,
		time.Duration(sleepInterval),
		identity,
		make(chan os.Signal, 1),
		make(chan error),
	}
}

func (ar *arCrawler) run() error {
	step := ar.step
	startBlockHeight := ar.fromHeight
	endBlockHeight := startBlockHeight + step
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
		startBlockHeight = startBlockHeight + step
		endBlockHeight = endBlockHeight + step

		if latestConfirmedBlockHeight <= endBlockHeight {
			time.Sleep(tempDelay)

			latestBlockHeight, err = GetLatestBlockHeight()

			if err != nil {
				return err
			}

			latestConfirmedBlockHeight = latestBlockHeight - ar.confirmations
			step = 10
		}

		log.Println("Getting articles from", startBlockHeight, "to", endBlockHeight,
			"with step", step, "and temp delay", tempDelay,
			"and latest confirmed block height", latestConfirmedBlockHeight,
		)
		ar.getArticles(startBlockHeight, endBlockHeight, ar.identity)
	}
}

func (ar *arCrawler) getArticles(from, to int64, owner string) error {
	logger.Infof("Getting articles from %d to %d", from, to)

	articles, err := GetArticles(from, to, owner)
	if err != nil {
		return err
	}

	logger.Info("Got articles:", len(articles))

	items := make([]*model.Item, len(articles))

	for _, article := range articles {
		attachment := model.Attachment{
			Type:     "body",
			Content:  article.Content,
			MimeType: "text/markdown",
		}

		tsp := time.Unix(article.TimeStamp, 0)

		author, err := rss3uri.NewInstance("account", article.Author, string(constants.PlatformSymbolEthereum))
		if err != nil {
			return fmt.Errorf("poap [%s] get new instance error:", err)
		}

		ni := model.NewItem(
			constants.NetworkIDArweaveMainnet,
			article.Digest,
			model.Metadata{
				"network": constants.NetworkSymbolArweaveMainnet,
				"proof":   article.Digest,
			},
			constants.ItemTagsMirrorEntry,
			[]string{author.String()},
			article.Title,
			article.Content, // TODO: According to RIP4, if the body is too long, then only record part of the body, followed by ... at the end
			[]model.Attachment{attachment},
			tsp,
		)

		items = append(items, ni)
		notes := []*model.ItemId{{
			NetworkID: constants.NetworkIDArweaveMainnet,
			Proof:     "", // TODO: @atlas
		}}
		db.AppendNotes(author, notes)
	}

	if len(items) > 0 {
		db.InsertItems(items, constants.NetworkSymbolArweaveMainnet.GetID())
	}

	return nil
}

func (ar *arCrawler) Start() error {
	signal.Notify(ar.interrupt, os.Interrupt)

	log.Println("Starting Arweave crawler...")

	go func() {
		ar.complete <- ar.run()
	}()

	for {
		select {
		case err := <-ar.complete:
			return err
		}
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
