package handler

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/gitcoin"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"

	"gorm.io/gorm"
)

// var reptileDB *gorm.DB

const identity = "gitcoin-project"
const networkID = 1000

func SetLastPostion(pos int) error {
	if _, err := database.CreateCrawlerMetadata(database.DB, &model.CrawlerMetadata{
		AccountInstance: identity,
		NetworkId:       networkID,
		LastBlock:       pos,
	}, true); err != nil {
		return fmt.Errorf("set last position error: %s", err)
	}

	return nil
}

func GetLastPostion() int {
	metadata, dbQcmErr := database.QueryCrawlerMetadata(database.DB, identity, networkID)
	if dbQcmErr != nil {
		logger.Errorf("query crawler metadata error: %s", dbQcmErr)

		return 0
	}

	return metadata.LastBlock
}

func createReptileGitcoinData(db *gorm.DB, project *gitcoin.ProjectInfo) error {
	if err := db.Clauses(database.NewCreateClauses(true)...).Create(project).Error; err != nil {
		return err
	}

	return nil
}

func SetResultInDB(project *gitcoin.ProjectInfo) error {
	if project == nil {
		return fmt.Errorf("project is nil")
	}

	if err := createReptileGitcoinData(database.DB, project); err != nil {
		return fmt.Errorf("set result in db error: %s ", err)
	}

	return nil
}
