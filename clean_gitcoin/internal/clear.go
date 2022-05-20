package clear

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
)

func GetDataFromDB(limit int, offset int) ([]model.Note, int64, error) {
	var notes []model.Note
	internalDB := database.DB.
		Where("attachments != '[]'and tags && '{\"Token\"}' order by date_created desc limit ? offset ?", limit, offset)

	var count int64
	if err := internalDB.Model(&model.Note{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return noteList, count, nil
}

func ClearGitCoinData(notes []model.Note) {
	for _, note := range notes {
		var relatedURLs []string = []string {
			moralis.GetTxHashURL(
		}
	}
}

func SaveDataInDB() {

}
