package main

import (
	"log"

	clear "github.com/NaturalSelectionLabs/RSS3-PreGod/clean_gitcoin/internal"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
)

func init() {
	if err := database.Setup(); err != nil {
		log.Fatalf("database.Setup err: %v", err)
	}
}

const GetNotesLimit = 2000
const platformID = constants.PlatformID(1300)
const crawlerID = "gitcoin-recovery-script"

func main() {
	offset, err := util.GetCrawlerMetadata(crawlerID, platformID)
	if err != nil {
		logger.Errorf("get crawler metadata error:%v", err)

		offset = 0
	}

	notes, err := clear.GetDataFromDB(GetNotesLimit, int(offset))
	if err != nil {
		logger.Infof("get data from db err:%v", err)

		return
	}

	logger.Infof("get %d notes", len(notes))

	/*
		for {
			// get data from db
			notes, err := clear.GetDataFromDB(GetNotesLimit, int(offset))
			if err != nil {
				logger.Infof("get data from db err:%v", err)

				return
			}

			if len(notes) == 0 {
				logger.Infof("mission completed")
			}

			// change db
			clear.ClearGitCoinData(notes)

			//save in db
			tx := database.DB.Begin()

			if _, err := database.CreateNotes(tx, notes, true); err != nil {
				continue
			}

			offset += GetNotesLimit
		}
	*/
}
