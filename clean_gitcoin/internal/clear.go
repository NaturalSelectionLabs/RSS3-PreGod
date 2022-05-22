package clear

import (
	"strings"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	for i, _ := range notes {
		notes[i].Tags = constants.ItemTagsDonationGitcoin.ToPqStringArray()
	}

	logger.Infof("note[0].Tags:%v", notes[0].Tags)
	return notes
}

func CreateNotes(db *gorm.DB, notes []model.Note, updateAll bool) ([]model.Note, error) {
	for i := range notes {
		notes[i].Identifier = strings.ToLower(notes[i].Identifier)
		notes[i].Owner = strings.ToLower(notes[i].Owner)

		if notes[i].Metadata == nil {
			notes[i].Metadata = []byte("{}")
		}

	}

	if err := db.Clauses(NewCreateClauses()...).Create(&notes).Error; err != nil {
		return nil, err
	}

	return notes, nil
}

func NewCreateClauses() []clause.Expression {
	clauses := []clause.Expression{
		// clause.Returning{}
	}

	clauses = append(clauses, clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{"tags"}),
		UpdateAll: true,
	})

	return clauses
}
