package crawler_handler

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util"
)

type GetUserBioHandler struct {
	CrawlerHandlerBase
}

type GetUserBioResult struct {
	CrawlerHandlerResultBase

	UserBio string
}

func NewGetUserBioHandler(workParam crawler.WorkParam) *GetUserBioHandler {
	return &GetUserBioHandler{
		CrawlerHandlerBase{
			WorkParam: workParam,
		},
	}
}

func NewGetUserBioResult() *GetUserBioResult {
	return &GetUserBioResult{
		CrawlerHandlerResultBase{
			Error: util.GetErrorBase(util.ErrorCodeSuccess),
		},
		"",
	}
}

func (pt *GetUserBioHandler) Excute() (*GetUserBioResult, error) {
	var err error

	var c crawler.Crawler

	var userBio string

	result := NewGetUserBioResult()

	c = MakeCrawlers(pt.WorkParam.PlatformID)
	if c == nil {
		result.Error = util.GetErrorBase(util.ErrorCodeNotSupportedNetwork)

		return result, fmt.Errorf("unsupported platform id[%d]", pt.WorkParam.PlatformID)
	}

	userBio, err = c.GetUserBio(pt.WorkParam.Identity)

	if err != nil {
		result.Error = util.GetErrorBase(util.ErrorCodeNotFoundData)

		return result, fmt.Errorf("[%s] can't find", pt.WorkParam.Identity)
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

	return result, nil
}
