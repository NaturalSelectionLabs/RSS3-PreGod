package internal

import (
	"regexp"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/moralis"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
)

func GetOneTokenSymbolEmptyIdentifier(chainType moralis.ChainType) (string, error) {
	// var owner string
	var note model.Note

	internalDB := database.DB.
		Where("(metadata->>'token_symbol')=''").
		Where("(metadata->>'token_address')!='' ").
		Where("\"source\" in ('Ethereum ERC20')").
		Where("\"metadata_network\"=?", chainType.GetNetworkSymbol()).
		Where("tags[1] like ('Token')")

	if err := internalDB.First(&note).Error; err != nil {
		return "", err
	}

	return note.Owner, nil
}

func GetAllNotesAboutErc20ByIdentifier(identifier string, chainType moralis.ChainType) (map[string]model.Note, error) {
	var notes []model.Note

	var noteMap = map[string]model.Note{}

	internalDB := database.DB.
		Where("owner in (?)", identifier).
		Where("(metadata->>'token_symbol')=''").
		Where("(metadata->>'token_address')!='' ").
		Where("\"source\" in ('Ethereum ERC20')").
		Where("\"metadata_network\"=?", chainType.GetNetworkSymbol()).
		Where("tags[1] like ('Token')")

	if err := internalDB.Find(&notes).Error; err != nil {
		return nil, err
	}

	for i := range notes {
		noteMap[notes[i].Identifier] = notes[i]
	}

	return noteMap, nil
}

func ChangeNotesTokenSymbolMsg(notesMap map[string]model.Note, tokensMap moralis.Erc20TokensMap) ([]model.Note, error) {
	var notes = []model.Note{}

	for _, note := range notesMap {
		noteMetadata, unwrapErr := database.UnwrapJSON[map[string]interface{}](note.Metadata)
		if unwrapErr != nil {
			logger.Warnf("unwrap metadata err:%v", unwrapErr)

			continue
		}

		tokenAddress, ok := noteMetadata["token_address"].(string)
		if !ok {
			logger.Warnf("Identifier [%s] token_address not found", note.Identifier)

			continue
		}

		noteMetadata["token_symbol"] = tokensMap[tokenAddress].Symbol
		note.Metadata = database.MustWrapJSON(noteMetadata)
		notes = append(notes, note)

		logger.Debugf("token_symbol:%s", tokensMap[tokenAddress].Symbol)
	}

	logger.Debugf("len(notes):%d", len(notes))

	return notes, nil
}

func GetAccountByIdentifier(identifier string) (string, error) {
	compileRegex := regexp.MustCompile("rss3://account:(.*?)@")
	matchArr := compileRegex.FindStringSubmatch(identifier)

	if len(matchArr) > 0 {
		return matchArr[len(matchArr)-1], nil
	}

	return "", nil
}
