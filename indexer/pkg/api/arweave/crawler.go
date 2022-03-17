package arweave

import (
	"strconv"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
)

func Crawl(param *crawler.WorkParam, result *crawler.CrawlerResult) (crawler.CrawlerResult, error) {
	startBlockHeight := int64(1)
	step := param.Step
	tempDelay := param.SleepInterval

	// TODO: why is this loop never ending?
	for {
		latestBlockHeight, err := GetLatestBlockHeight()
		if err != nil {
			return *result, err
		}

		endBlockHeight := startBlockHeight + step
		if latestBlockHeight <= endBlockHeight {
			time.Sleep(tempDelay)

			latestBlockHeight, err = GetLatestBlockHeight()
			if err != nil {
				return *result, err
			}

			endBlockHeight = latestBlockHeight
			step = 10
		}

		err = getArticles(result, startBlockHeight, endBlockHeight, param.Identity)

		if err != nil {
			return *result, err
		}
	}
}

func getArticles(result *crawler.CrawlerResult, from int64, to int64, owner string) error {
	articles, err := GetArticles(from, to, owner)
	if err != nil {
		return err
	}

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

		result.Items = append(result.Items, ni)
	}

	return nil
}
