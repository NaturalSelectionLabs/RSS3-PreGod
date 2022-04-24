package handler

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/gitcoin"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
)

// var reptileDB *gorm.DB

const GitcoinProjectIdentity = "gitcoin-project"
const GitcoinProjectPlatformID = 1000
const GitcoinProjectIdentity2 = "gitcoin-project2"
const GitcoinProjectPlatformID2 = 1001

func SetLastPostion(pos int, identity string, platformID constants.PlatformID) error {
	if _, err := database.CreateCrawlerMetadata(database.DB, &model.CrawlerMetadata{
		AccountInstance: identity,
		PlatformID:      platformID,
		LastBlock:       pos,
	}, true); err != nil {
		return fmt.Errorf("set last position error: %s", err)
	}

	return nil
}

func GetLastPostion(identity string, platformID constants.PlatformID) int {
	metadata, dbQcmErr := database.QueryCrawlerMetadata(database.DB, identity, platformID)
	if dbQcmErr != nil {
		logger.Errorf("query crawler metadata error: %s", dbQcmErr)

		return 0
	}

	return metadata.LastBlock
}

func SetResultInDB(project *gitcoin.ProjectInfo) error {
	if project == nil {
		return fmt.Errorf("project is nil")
	}

	if err := database.DB.Clauses(database.NewCreateClauses(true)...).Create(project).Error; err != nil {
		return fmt.Errorf("set result in db error: %s ", err)
	}

	return nil
}

func SetResultsInDB(projects []gitcoin.ProjectInfo) error {
	if err := database.DB.Clauses(database.NewCreateClauses(true)...).Create(&projects).Error; err != nil {
		return err
	}

	return nil
}

func GetResultMax() int {
	max := 0

	rows, err := database.DB.Table("reptile-gitcoin.data").Select("max(id)").Rows()
	if err != nil {
		logger.Errorf("get result max error: %s", err)
	}

	if rows.Next() {
		err := rows.Scan(&max)
		if err != nil {
			logger.Errorf("get result max error: %s", err)
		}
	}

	return max
}

func GetResultFromDB(pos int, endpos int) ([]gitcoin.ProjectInfo, error) {
	var projects []gitcoin.ProjectInfo

	internalDB := database.DB.
		Where("id > ?", pos).
		Where("id <= ?", endpos)

	if err := internalDB.Find(&projects).Error; err != nil {
		return nil, err
	}

	return projects, nil
}

func UpdateResultsInDB(project *gitcoin.ProjectInfo, adminAddress string) error {
	if err := database.DB.Model(project).Where("id=?", project.Id).Update("admin_address", adminAddress).Error; err != nil {
		logger.Errorf("update result in db error: %s", err)
	}

	return nil
}
