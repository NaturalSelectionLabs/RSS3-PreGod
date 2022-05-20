package internal

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
)

func GetDataFromDB() {
	var notes []model.Note
	internalDB := database.DB.
		Where("attachments != '[]'and tags && '{\"Token\"}' order by date_created desc limit ? offset ?")
}

func ClearGitCoinData() {

}

func SaveDataInDB() {

}
