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
	internalDB := database.DB.Where("owner = ?", strings.ToLower(rss3uri.New(req.InstanceList[0]).String()))
	for _, instance := range req.InstanceList[1:] {
		internalDB = internalDB.Or("owner = ?", strings.ToLower(rss3uri.New(instance).String()))
	}
	internalDB = internalDB.Order("date_created DESC").Order("identifier DESC")

	var count int64
	if err := internalDB.Model(&model.Note{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	noteList := []model.Note{}
	if err := internalDB.Limit(req.Limit).Offset((req.Page - 1) * req.Limit).Find(&noteList).Error; err != nil {
		return nil, 0, err
	}

	return noteList, count, nil
}
