package dao

import (
	"fmt"
	"strings"

	m "github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/lib/pq"
)

// BatchGetNodeList query data through database
func BatchGetNodeList(req m.BatchGetNodeListRequest) ([]model.Note, int64, error) {
	internalDB := database.DB
	ownerList := make([]string, 0)

	for _, instance := range req.InstanceList {
		owner := strings.ToLower(rss3uri.New(instance).String())
		unknownPlatform := fmt.Sprintf("%v@unknown", strings.Split(owner, "@")[0])
		ownerList = append(ownerList, owner, unknownPlatform)
	}

	if req.Tags != nil && len(req.Tags) != 0 {
		internalDB = internalDB.Where("tags && ?", pq.StringArray(req.Tags))
	}

	if req.ExcludeTags != nil && len(req.ExcludeTags) != 0 {
		internalDB = internalDB.Where("tags && ? = FALSE", pq.StringArray(req.ExcludeTags))
	}

	if req.ItemSources != nil && len(req.ItemSources) != 0 {
		internalDB = internalDB.Where("source IN ?", req.ItemSources)
	}

	if req.Networks != nil && len(req.Networks) != 0 {
		internalDB = internalDB.Where("metadata_network IN ?", req.Networks)
	}

	if len(req.LastIdentifier) > 0 {
		lastItem := model.Note{}
		if err := database.DB.Where(&model.Note{
			Identifier: strings.ToLower(req.LastIdentifier),
		}).First(&lastItem).Error; err != nil {
			return nil, 0, err
		}

		internalDB = internalDB.Where("date_created <= ?", lastItem.DateCreated).
			Where("identifier != ?", lastItem.Identifier)
	}

	// filter out user active transactions
	ethNotes := make([]model.Note, 0)
	ethMap := make(map[string]bool)

	if err := internalDB.
		Where("owner IN ?", ownerList).
		Where("source = ?", constants.NoteSourceNameEthereumNFT).
		Order("date_created DESC").
		Find(&ethNotes).Error; err != nil {
		return nil, 0, err
	}

	activeTxList := []string{}

	for _, ethNote := range ethNotes {
		h := strings.Split(ethNote.MetadataProof, "-")[0]
		if ok := ethMap[h]; !ok {
			activeTxList = append(activeTxList, h)
			ethMap[h] = true
		}
	}

	internalDB = internalDB.
		Where("owner IN ?", ownerList).
		Where("(source = ? AND (metadata ->> 'transaction_hash' IN ? OR tags && ?)) OR source != ?",
			constants.NoteSourceNameEthereumNFT,
			activeTxList,
			pq.StringArray([]string{"POAP"}),
			constants.NoteSourceNameEthereumNFT).
		Order("date_created DESC").
		Order("contract_address DESC").
		Order("log_index DESC").
		Order("token_id DESC")

	var count int64
	if err := internalDB.Model(&model.Note{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	noteList := make([]model.Note, 0)
	if err := internalDB.Limit(req.Limit).Find(&noteList).Error; err != nil {
		return nil, 0, err
	}

	return noteList, count, nil
}
