package clear

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
)

func GetDataFromDB(limit int, offset int) ([]model.Note, error) {
	var notes []model.Note

	internalDB := database.DB.
		Where("attachments != '[]'").
		Where("tags && '{\"Token\"}'").
		Order("date_created DESC").
		Limit(limit).
		Offset(offset)

	var count int64
	if err := internalDB.Model(&model.Note{}).Count(&count).Error; err != nil {
		return nil, err
	}

	return notes, nil
}

func ClearGitCoinData(notes []model.Note) {
	// get projects
	for _, note := range notes {
		note.Tags = constants.ItemTagsToken.ToPqStringArray()
	}
}
