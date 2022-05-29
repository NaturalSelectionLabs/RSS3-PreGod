package internal

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/valyala/fastjson"
)

func GetDataFromDB(limit int) ([]model.Note, error) {
	var notes []model.Note

	internalDB := database.DB.
		Where("identifier not like ('rss3://note:%@ethereum') ").
		Where("related_urls[1] like ('https://etherscan.io/tx/%')").
		Where("\"source\" in ('Ethereum ERC20')").
		Order("date_created DESC").
		Limit(limit)

		// internalDB := database.DB.
		// Where("\"identifier\" in ('rss3://note:0x5c170dfde06db67469eb32c7bbe40d5bfe987766279bde14d6906dd231b65825-0@bnb')")

		// logger.Debugf("Limit:%d, Offset:%d", limit, offset)

	if err := internalDB.Find(&notes).Error; err != nil {
		return nil, err
	}

	return notes, nil
}

func ReplaceEndpoint(notes []model.Note) {
	// get projects
	for i := range notes {
		var parser fastjson.Parser
		parsedJson, err := parser.Parse(string(notes[i].Metadata))

		if err != nil {
			logger.Errorf("parse metadata err:%v", err)

			continue
		}

		transactionHash := string(parsedJson.GetStringBytes("transaction_hash"))
		network := constants.NetworkSymbol(parsedJson.GetStringBytes("network"))

		notes[i].RelatedURLs = []string{
			GetTxHashURL(network, transactionHash),
		}
		// logger.Infof("note[i]:%v", notes[i])
	}
}

func GetTxHashURL(
	network constants.NetworkSymbol,
	transactionHash string,
) string {
	switch network {
	case constants.NetworkSymbolEthereum:
		return "https://etherscan.io/tx/" + (transactionHash)

	case constants.NetworkSymbolPolygon:
		return "https://polygonscan.com/tx/" + (transactionHash)

	case constants.NetworkSymbolBNBChain:
		return "https://bscscan.com/tx/" + (transactionHash)

	case constants.NetworkSymbolAvalanche:
		return "https://avascan.info/blockchain/c/tx/" + (transactionHash)
	case constants.NetworkSymbolFantom:
		return "https://ftmscan.com/tx/" + (transactionHash)
	case constants.NetworkSymbolZkSync:
		return "https://zkscan.io/explorer/transactions/" + (transactionHash)
	default:
		return ""
	}
}