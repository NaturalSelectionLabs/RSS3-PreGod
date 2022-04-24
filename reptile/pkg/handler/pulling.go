package handler

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/gitcoin"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
)

func addressToLower(project *gitcoin.ProjectInfo) {
	if project == nil {
		return
	}

	project.AdminAddress = strings.ToLower(project.AdminAddress)
	project.ContractAddress = strings.ToLower(project.ContractAddress)
	project.TokenAddress = strings.ToLower(project.TokenAddress)
}

func setProjectInfoInDB(currentPos int64) (bool, error) {
	getProject := false

	projectInfo, err := GetResult(currentPos)
	if err != nil {
		return getProject, err
	}

	addressToLower(projectInfo)

	if projectInfo != nil && err == nil {
		err = SetResultInDB(projectInfo)
		if err != nil {
			logger.Fatal(err)
		}

		getProject = true
	}

	return getProject, nil
}

func printProjectInfoStage(currentPos int64, getProject bool) {
	getProjectStr := "true"

	if !getProject {
		getProjectStr = "false"
	}

	log.Printf("get [%d] project info stage: %s", currentPos, getProjectStr)
	log.Printf("current position: %d", currentPos)
	log.Printf("------------------------------------------\n")
}

// -- pull all results up to a certain value

func PullInformation(endpos int64) {
	lastPos := GetLastPostion(
		GitcoinProjectIdentity,
		GitcoinProjectPlatformID)
	currentPos := lastPos + 1

	for {
		netTag := true
		getProject, err := setProjectInfoInDB(currentPos)

		if err != nil {
			if err.Error() == "get result false:StatusCode [403]" {
				netTag = false
			}
		}

		err = SetLastPostion(
			currentPos,
			GitcoinProjectIdentity,
			GitcoinProjectPlatformID)
		if err != nil {
			logger.Errorf("set last position error:%s", err)
		}

		printProjectInfoStage(currentPos, getProject)

		if lastPos == endpos { // normal is 6000
			break
		}

		if netTag == true {
			lastPos = currentPos
			currentPos = lastPos + 1
		}

		time.Sleep(50 * time.Millisecond)
	}
}

// -- set db data dress to lower needed
// If the addresses that are pulled down at the beginning are all uppercase addresses,
// convert them to lowercase through this

func changeAddressFromDB(lastpos int64, endPos int64) error {
	projects, err := GetResultFromDB(lastpos, endPos)
	if err != nil {
		return fmt.Errorf("get result from db error:%s", err)
	}

	processingProject := []gitcoin.ProjectInfo{}

	for _, project := range projects {
		addressToLower(&project)
	}

	if len(processingProject) > 0 {
		err = SetResultsInDB(processingProject)
		if err != nil {
			logger.Fatal(err)
		}
	}

	return nil
}

func SetDBDataDressToLower() {
	var rangeMax int64 = 100

	resultcount := GetResultMax()

	for {
		lastpos := GetLastPostion(
			GitcoinProjectIdentity2,
			GitcoinProjectPlatformID2)

		if lastpos >= resultcount {
			logger.Infof("lastPos[%d] is the latest pos", lastpos)

			break
		}

		nextPos := lastpos + rangeMax
		if nextPos >= resultcount {
			nextPos = resultcount
		}

		if err := changeAddressFromDB(lastpos, nextPos); err != nil {
			logger.Fatal(err)
		}

		if err := SetLastPostion(
			nextPos,
			GitcoinProjectIdentity2,
			GitcoinProjectPlatformID2); err != nil {
			logger.Fatal(err)
		}
	}
}

// -- timing start to store data

func GetResultByStage() {
	rangeMax := 50

	lastPos := GetResultMax()
	currentPos := lastPos + 1
	countPos := 0

	projectCountArr := []gitcoin.ProjectInfo{}

	for {
		if countPos >= rangeMax {
			break
		}

		getProject := false

		projectInfo, err := GetResult(currentPos)
		if projectInfo != nil && err == nil {
			getProject = true
		} else if projectInfo == nil {
			goto END
		}

		if err != nil {
			logger.Infof("get result error:%s", err)

			if err.Error() == "get result false:StatusCode [403]" {
				break
			} else {
				goto END
			}
		}

		addressToLower(projectInfo)

		projectCountArr = append(projectCountArr, *projectInfo)

	END:
		printProjectInfoStage(currentPos, getProject)

		lastPos = currentPos

		currentPos = lastPos + 1
		countPos++
	}

	if err := SetResultsInDB(projectCountArr); err != nil {
		logger.Fatal(err)
	}
}
