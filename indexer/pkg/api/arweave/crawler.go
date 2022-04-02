package arweave

import (
	"errors"
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

type crawlConfig struct {
	fromHeight    int64
	confirmations int64
	step          int64
	minStep       int64
	sleepInterval time.Duration
}

type crawler struct {
	identity  ArAccount
	interrupt chan os.Signal
	complete  chan error
	cfg       *crawlConfig
}

func NewCrawler(identity ArAccount, crawlCfg *crawlConfig) *crawler {
	return &crawler{
		identity,
		make(chan os.Signal, 1),
		make(chan error),
		crawlCfg,
	}
}

func (ar *crawler) run() error {
	startBlockHeight := ar.cfg.fromHeight
	step := ar.cfg.step
	endBlockHeight := startBlockHeight + step
	tempDelay := ar.cfg.sleepInterval

	latestConfirmedBlockHeight, err := GetLatestBlockHeightWithConfirmations(ar.cfg.confirmations)
	if err != nil {
		logger.Errorf("get latest block error: %v", err)

		return err
	}

	for {
		// handle interrupt
		if ar.gotInterrupt() {
			return ErrInterrupt
		}

		if latestConfirmedBlockHeight <= endBlockHeight {
			for {
				time.Sleep(tempDelay)

				latestConfirmedBlockHeight, err = GetLatestBlockHeightWithConfirmations(ar.cfg.confirmations)
				if err != nil {
					logger.Errorf("get latest block error: %v", err)

					return err
				}

				if latestConfirmedBlockHeight > endBlockHeight {
					break
				}
			}

			// use minStep if we are at the end of the chain
			step = ar.cfg.minStep
		}

		log.Println("Getting articles from", startBlockHeight, "to", endBlockHeight,
			"with step", step, "and temp delay", tempDelay,
			"and latest confirmed block height", latestConfirmedBlockHeight,
		)

		ar.parseMirrorArticles(startBlockHeight, endBlockHeight, ar.identity)

		startBlockHeight = endBlockHeight
		endBlockHeight = startBlockHeight + step
	}
}

//TODO: I think it will be the same as other crawler formats in the future,
// and it will return to an abstract and unified crawler
func (ar *crawler) parseMirrorArticles(from, to int64, owner ArAccount) error {
	articles, err := GetMirrorContents(from, to, owner)
	if err != nil {
		logger.Errorf("GetMirrorContents error: [%v]", err)

		return err
	}

	logger.Info("Got articles:", len(articles))

	items := make([]*model.Item, 0)

	for _, article := range articles {
		attachment := model.Attachment{
			Type:     "body",
			Content:  article.Content,
			MimeType: "text/markdown",
		}

		tsp := time.Unix(article.Timestamp, 0)

		author, err := rss3uri.NewInstance("account", article.Author, string(constants.PlatformSymbolEthereum))
		if err != nil {
			//TODO: may send to a error queue or whatever in the future
			logger.Errorf("arweave NewInstance error: [%v]", err)

			tsp = time.Now()
		}

		logger.Infof("author: [%v]", author)

		ni := model.NewItem(
			constants.NetworkIDArweaveMainnet,
			article.TxHash,
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
		notes := []*model.ObjectId{{
			NetworkID: constants.NetworkIDArweaveMainnet,
			Proof:     article.TxHash,
		}}
		instance := rss3uri.NewAccountInstance(article.Author, constants.PlatformSymbolArweave)
		db.AppendNotes(instance, notes)
	}

	db.InsertItems(items, constants.NetworkIDArweaveMainnet)

	return nil
}

func (ar *crawler) Start() error {
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

func (ar *crawler) gotInterrupt() bool {
	select {
	case <-ar.interrupt:
		signal.Stop(ar.interrupt)

		return true
	default:
		return false
	}
}
