package dao

import (
	"strings"

	m "github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
)

// BatchGetNodeList query data through database
func BatchGetNodeList(req m.BatchGetNodeListRequest) ([]model.Note, int64, error) {
	internalDB := database.DB
	ownerList := []string{}

	for _, instance := range req.InstanceList {
		ownerList = append(ownerList, strings.ToLower(rss3uri.New(instance).String()))
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

	internalDB = internalDB.Where("owner IN ?", ownerList).Order("date_created DESC").Order("identifier DESC")

	var count int64
	if err := internalDB.Model(&model.Note{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	noteList := []model.Note{}
	if err := internalDB.Limit(req.Limit).Find(&noteList).Error; err != nil {
		return nil, 0, err
	}

	return noteList, count, nil
}
