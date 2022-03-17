package processor

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/jike"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/misskey"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/moralis"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/twitter"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	jsoniter "github.com/json-iterator/go"
)

var jsoni = jsoniter.ConfigCompatibleWithStandardLibrary

func Dispatch(args ...string) (crawler.CrawlerResult, error) {
	var result crawler.CrawlerResult

	param := new(crawler.WorkParam)

	// unmarshal the first argument from string to WorkParam
	err := jsoni.UnmarshalFromString(args[0], &param)

	if err != nil {
		return result, err
	}

	switch param.NetworkID {
	case constants.NetworkIDEthereumMainnet,
		constants.NetworkIDBNBChain,
		constants.NetworkIDAvalanche,
		constants.NetworkIDFantom,
		constants.NetworkIDPolygon:
		return moralis.Crawl(param, &result)
	case constants.NetworkIDMisskey:
		return misskey.Crawl(param, &result)
	case constants.NetworkIDJike:
		return jike.Crawl(param, &result)
	case constants.NetworkIDTwitter:
		return twitter.Crawl(param, &result)
	default:
		return result, fmt.Errorf("unknown network id")
	}
}
