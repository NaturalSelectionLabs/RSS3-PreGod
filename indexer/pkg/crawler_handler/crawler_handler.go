package crawler_handler

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/jike"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/misskey"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/moralis"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/twitter"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
)

type CrawlerHandlerResultInf interface {
}

type CrawlerHandlerResultBase struct {
	CrawlerHandlerResultInf
	Error util.ErrorBase
}

type CrawlerHandlerInf interface {
	Excute() CrawlerHandlerResultInf
}

type CrawlerHandlerBase struct {
	WorkParam crawler.WorkParam
	CrawlerHandlerInf
}

func MakeCrawlers(network constants.NetworkID) crawler.Crawler {
	switch network {
	case constants.NetworkIDEthereumMainnet,
		constants.NetworkIDBNBChain,
		constants.NetworkIDAvalanche,
		constants.NetworkIDFantom,
		constants.NetworkIDPolygon:
		return moralis.NewMoralisCrawler()
	case constants.NetworkIDJike:
		return jike.NewJikeCrawler()
	case constants.NetworkIDTwitter:
		return twitter.NewTwitterCrawler()
	case constants.NetworkIDMisskey:
		return misskey.NewMisskeyCrawler()
	default:
		return nil
	}
}
