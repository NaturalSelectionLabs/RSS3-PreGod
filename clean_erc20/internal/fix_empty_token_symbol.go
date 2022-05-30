package internal

import (
	"regexp"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/moralis"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
)

type NodeUnit struct {
	Note         model.Note
	TokenAddress string
}

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

func GetNotesAboutErc20ByIdentifier(chainType moralis.ChainType, pagSize int) (map[string]NodeUnit, error) {
	var notes []model.Note

	var noteMap = map[string]NodeUnit{}

	internalDB := database.DB.
		Where("(metadata->>'token_symbol')=''").
		Where("(metadata->>'token_address')!='' ").
		Where("\"source\" in ('Ethereum ERC20')").
		Where("\"metadata_network\"=?", chainType.GetNetworkSymbol()).
		Where("tags[1] like ('Token')").
		Order("date_created ASC").
		Limit(pagSize)

	if err := internalDB.Find(&notes).Error; err != nil {
		return nil, err
	}

	for i := range notes {
		noteMetadata, unwrapErr := database.UnwrapJSON[map[string]interface{}](notes[i].Metadata)
		if unwrapErr != nil {
			logger.Warnf("unwrap metadata err:%v", unwrapErr)

			continue
		}

		tokenAddress, ok := noteMetadata["token_address"].(string)
		if !ok {
			logger.Warnf("Identifier [%s] token_address not found", notes[i].Identifier)

			continue
		}

		nodeUnit := NodeUnit{
			Note:         notes[i],
			TokenAddress: tokenAddress,
		}

		noteMap[notes[i].Identifier] = nodeUnit
	}

	return noteMap, nil
}

func ChangeNotesTokenSymbolMsg(notesMap map[string]NodeUnit, tokensMap moralis.Erc20TokensMap) ([]model.Note, error) {
	var notes = []model.Note{}

	for _, noteUnit := range notesMap {
		noteMetadata, unwrapErr := database.UnwrapJSON[map[string]interface{}](noteUnit.Note.Metadata)
		if unwrapErr != nil {
			logger.Warnf("unwrap metadata err:%v", unwrapErr)

			continue
		}

		tokenSymbol, ok := tokensMap[noteUnit.TokenAddress]
		if !ok {
			logger.Warnf("token symbol not found for token address:%s", noteUnit.TokenAddress)

			continue
		}

		noteMetadata["token_symbol"] = tokenSymbol.Symbol
		noteUnit.Note.Metadata = database.MustWrapJSON(noteMetadata)
		notes = append(notes, copyNote(noteUnit.Note))

		// logger.Debugf("tokenAddress:%s, token_symbol:%s", noteUnit.TokenAddress, tokenSymbol)
	}

	// logger.Debugf("len(notes):%d", len(notes))

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

func copyNote(note model.Note) model.Note {
	newNote := model.Note{
		Identifier:          note.Identifier,
		TransactionHash:     note.TransactionHash,
		TransactionLogIndex: note.TransactionLogIndex,
		Owner:               note.Owner,
		ProfileSourceID:     note.ProfileSourceID,
		RelatedURLs:         note.RelatedURLs,
		Tags:                note.Tags,
		Authors:             note.Authors,
		Title:               note.Title,
		Summary:             note.Summary,
		Attachments:         note.Attachments,
		Source:              note.Source,
		MetadataNetwork:     note.MetadataNetwork,
		MetadataProof:       note.MetadataProof,
		Metadata:            note.Metadata,
		DateCreated:         note.DateCreated,
		DateUpdated:         note.DateUpdated,
	}

	return newNote
}
