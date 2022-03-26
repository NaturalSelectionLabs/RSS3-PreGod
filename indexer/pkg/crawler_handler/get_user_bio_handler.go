package crawler_handler

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
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

func (pt *GetUserBioHandler) Excute() *GetUserBioResult {
	var err error

	var c crawler.Crawler

	var userBio string

	result := NewGetUserBioResult()

	c = MakeCrawlers(pt.WorkParam.NetworkID)
	if c == nil {
		result.Error = util.GetErrorBase(util.ErrorCodeNotSupportedNetwork)

		logger.Errorf("unsupported network id[%d]", pt.WorkParam.NetworkID)

		return result
	}

	userBio, err = c.GetUserBio(pt.WorkParam.Identity)

	if err != nil {
		result.Error = util.GetErrorBase(util.ErrorCodeNotFoundData)

		logger.Errorf("[%d] can't find", pt.WorkParam.Identity)

		return result
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

	return result
}
