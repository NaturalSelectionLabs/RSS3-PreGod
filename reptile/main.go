package main

import (
	"log"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/reptile/pkg/handler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
)

func init() {
	if err := database.Setup(); err != nil {
		log.Fatalf("database.Setup err: %v", err)
	}

	if err := handler.Setup(); err != nil {
		log.Fatalf("handler.Setup err: %v", err)
	}
}

var firstTime = true

func main() {
	lastPos := 0
	currentPos := 0

	for {
		if firstTime == true {
			lastPos = handler.GetLastPostion()
			firstTime = false
		}

		currentPos = lastPos + 1
		getProject := false

		projectInfo, err := handler.GetResult(currentPos)
		if err != nil {
			if err == "get result false:StatusCode [403]" {
				continue
			}
			logger.Fatal(err)
		}

		if projectInfo != nil && err == nil {
			err = handler.SetResultInDB(projectInfo)
			if err != nil {
				logger.Fatal(err)
			}

			getProject = true
		}

		err = handler.SetLastPostion(currentPos)
		if err != nil {
			logger.Fatal(err)
		}

		getProjectStr := "false"
		if getProject == true {
			getProjectStr = "true"
		}

		log.Printf("get [%d] project info stage: %s", currentPos, getProjectStr)
		log.Printf("current position: %d", currentPos)
		log.Printf("------------------------------------------\n")

		if lastPos == 6000 {
			break
		}

		lastPos = currentPos

		time.Sleep(50 * time.Millisecond)
	}
}
