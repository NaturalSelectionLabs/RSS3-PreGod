package main

import (
	"log"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/clean_erc20/internal"
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
const crawlerID = "erc20-recovery-script"

func main() {
	offset, err := util.GetCrawlerMetadata(crawlerID, platformID)
	if err != nil {
		logger.Errorf("get crawler metadata error:%v", err)

		offset = 0
	}

	notes, err := internal.GetDataFromDB(1, int(offset))
	if err != nil {
		logger.Infof("get data from db err:%v", err)

		return
	}

	// logger.Infof("get %d notes", len(notes))
	// logger.Debugf("notes:%v", notes)

	// for {
	// get data from db
	// notes, err := clear.GetDataFromDB(GetNotesLimit, int(offset))
	// if err != nil {
	// 	logger.Infof("get data from db err:%v", err)

	// 	return
	// }

	if len(notes) == 0 {
		logger.Infof("mission completed")
	}

	// change db
	internal.ClearGitCoinData(notes)

	//save in db
	// tx := database.DB.Begin()

	// logger.Infof("notes[0].tags:%v", notes[0].Tags)

	if _, err := database.CreateNotes(database.DB, notes, true); err != nil {
		// continue
	}

	logger.Debugf("note[0].RelatedURLs:%v", notes[0].RelatedURLs)

	// set the current block height as the from height
	if err := util.SetCrawlerMetadata(crawlerID, offset, platformID); err != nil {
		logger.Errorf("create crawler metadata error: %v", err)
	}

	offset += GetNotesLimit

	logger.Infof("offset:%d", offset)
	// }

}
