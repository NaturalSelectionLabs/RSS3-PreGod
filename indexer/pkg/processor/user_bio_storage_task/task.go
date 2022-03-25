package user_bio_storage_task

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/processor"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
)

type UserBioStorageTask struct {
	processor.ProcessTaskParam
	ResultQ chan *UserBioStorageResult
}

type UserBioStorageResult struct {
	processor.ProcessTaskResult

	UserBio string
}

func NewEmptyUserBioStorageResult() *UserBioStorageResult {
	return &UserBioStorageResult{
		processor.ProcessTaskResult{
			TaskType:   processor.ProcessTaskTypeUserBioStorage,
			TaskResult: processor.ProcessTaskErrorCodeSuccess,
		},

		"",
	}
}

func (pt *UserBioStorageTask) Fun() error {
	var err error

	var c crawler.Crawler

	var userBio string

	result := NewEmptyUserBioStorageResult()

	c = processor.MakeCrawler(pt.WorkParam.NetworkID)
	if c == nil {
		result.TaskResult = processor.ProcessTaskErrorCodeNotSupportedNetwork

		logger.Errorf("unsupported network id: %d", pt.WorkParam.NetworkID)

		goto RETURN
	}

	logger.Infof("c:%v", &c)

	userBio, err = c.GetUserBio(pt.WorkParam.Identity)

	if err != nil {
		result.TaskResult = processor.ProcessTaskErrorCodeNotFoundData

		goto RETURN
	}

	if len(userBio) > 0 {
		// TODO：add userbio into redis
		// redis.SetUserBio(userBio)
		// ctx := context.Background()
		// key := fmt.Sprintf("%s_%s_%s", pt.WorkParam.Identity,
		// 	pt.WorkParam.PlatformID.Symbol(),
		// )
		// cache.Set(ctx, key, userBio, 2)
		result.UserBio = userBio
	}

RETURN:
	pt.ResultQ <- result

	if err != nil {
		logger.Error(err)

		return err
	} else {
		return nil
	}
}
