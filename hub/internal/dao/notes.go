package dao

import (
	"strings"

	m "github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/lib/pq"
)

// BatchGetNodeList query data through database
func BatchGetNodeList(req m.BatchGetNodeListRequest) ([]model.Note, int64, error) {
	internalDB := database.DB
	ownerList := make([]string, 0)

	for _, instance := range req.InstanceList {
		ownerList = append(ownerList, strings.ToLower(rss3uri.New(instance).String()))
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
			Where(
				"(transaction_hash = ? and transaction_log_index < ?) or (transaction_hash < ?)",
				lastItem.TransactionHash, lastItem.TransactionLogIndex, lastItem.TransactionHash,
			).
			Where("identifier != ?", lastItem.Identifier)
	}

	internalDB = internalDB.
		Where("owner IN ?", ownerList).
		Order("date_created DESC").
		Order("transaction_hash DESC").
		Order("transaction_log_index DESC")

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
