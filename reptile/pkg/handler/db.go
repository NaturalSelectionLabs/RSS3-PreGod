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
const GitcoinProjectNetworkID = 1000
const GitcoinProjectIdentity2 = "gitcoin-project2"
const GitcoinProjectNetworkID2 = 1001

func SetLastPostion(pos int, identity string, networkID constants.NetworkID) error {
	if _, err := database.CreateCrawlerMetadata(database.DB, &model.CrawlerMetadata{
		AccountInstance: identity,
		NetworkId:       networkID,
		LastBlock:       pos,
	}, true); err != nil {
		return fmt.Errorf("set last position error: %s", err)
	}

	return nil
}

func GetLastPostion(identity string, networkID constants.NetworkID) int {
	metadata, dbQcmErr := database.QueryCrawlerMetadata(database.DB, identity, networkID)
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

func GetResultTotal() int64 {
	var count int64

	if err := database.DB.Model(&gitcoin.ProjectInfo{}).Distinct("id").Count(&count).Error; err != nil {
		logger.Errorf("get result total error: %s", err)
	}

	return count
}

func GetResultFromDB(pos int, endpos int) ([]gitcoin.ProjectInfo, error) {
	var projects []gitcoin.ProjectInfo

	internalDB := database.DB.
		Where("id > ?", pos).
		Where("id <= ?", pos+endpos)

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
