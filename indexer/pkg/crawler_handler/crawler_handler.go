package crawler_handler

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/moralis"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/poap"
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

func MakeCrawlers[T constants.NetworkID | constants.PlatformID](network T) crawler.Crawler {
	switch any(network).(type) {
	case constants.NetworkID:
		switch constants.NetworkID(network) {
		case constants.NetworkIDEthereum,
			constants.NetworkIDBNBChain,
			constants.NetworkIDAvalanche,
			constants.NetworkIDFantom,
			constants.NetworkIDPolygon:
			return moralis.NewMoralisCrawler()
		case constants.NetworkIDGnosisMainnet:
			return poap.NewPoapCrawler()
		default:
			return nil
		}

	case constants.PlatformID:
		switch constants.PlatformID(network) {
		case constants.PlatformIDEthereum:
			return moralis.NewMoralisCrawler()
		default:
			return nil
		}

	default:
		return nil
	}
}
