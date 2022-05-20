package clear

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
)

func GetDataFromDB(limit int, offset int) ([]model.Note, error) {
	var notes []model.Note

	internalDB := database.DB.
		Where("attachments != '[]'").
		Where("tags && '{\"Token\"}'").
		Order("date_created DESC").
		Limit(limit).
		Offset(offset)

	if err := internalDB.Find(&notes).Error; err != nil {
		return nil, err
	}

	return notes, nil
}

func ClearGitCoinData(notes []model.Note) []model.Note {
	// get projects
	for i, note := range notes {
		logger.Infof("note.Tags:%v", note.Tags)
		note.Tags = constants.ItemTagsDonationGitcoin.ToPqStringArray()
		logger.Infof("note.Tags:%v， %d", note.Tags, i)
	}

	logger.Infof("note[0].Tags:%v", notes[0].Tags)
	return notes
}
