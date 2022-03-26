package api

import (
	"net/http"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler_handler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/gin-gonic/gin"
)

type GetBioRequest struct {
	Identity   string               `form:"proof" binding:"required"`
	PlatformID constants.PlatformID `form:"platform_id" binding:"required"`
}

type GetBioResponse struct {
	util.ErrorBase `json:"error"`
	UserBio        string `json:"user_bio"`
}

var (
	// Since the transmitted parameter is only PlatformID
	// Currently, the platform and network for pulling bio are the same
	// , so there is a need for a place to transfer to each other.
	platform2Network = map[constants.PlatformID]constants.NetworkID{
		constants.PlatformIDTwitter: constants.NetworkIDTwitter,
		constants.PlatformIDJike:    constants.NetworkIDJike,
		constants.PlatformIDMisskey: constants.NetworkIDMisskey,
	}
)

func GetBioHandlerFunc(c *gin.Context) {
	request := GetBioRequest{}
	if err := c.ShouldBind(&request); err != nil {
		logger.Errorf("request bind error: %s", err.Error())

		return
	}

	response := GetBioResponse{
		util.ErrorBase{},
		"",
	}

	if len(request.Identity) > 0 || !constants.IsValidPlatformSymbol(string(request.PlatformID.Symbol())) {
		logger.Errorf("parameter error")

		response.ErrorBase = util.GetErrorBase(util.ErrorCodeParameterError)
		c.JSON(http.StatusOK, response)

		return
	}

	getuserBioHandler := crawler_handler.NewGetUserBioHandler(
		crawler.WorkParam{
			Identity:   request.Identity,
			PlatformID: request.PlatformID,
			NetworkID:  platform2Network[request.PlatformID],
		})
	handlerResult := getuserBioHandler.Excute()

	if handlerResult == nil {
		logger.Errorf("[%s] get user bio result error", request.Identity)

		response.ErrorBase = util.GetErrorBase(util.ErrorCodeNotFoundData)
		c.JSON(http.StatusOK, response)

		return
	}

	if handlerResult.Error.ErrorCode != util.ErrorCodeSuccess {
		logger.Errorf("[%s] get user bio result error", request.Identity)

		response.ErrorBase = util.GetErrorBase(handlerResult.Error.ErrorCode)
		c.JSON(http.StatusOK, response)

		return
	}

	response.UserBio = handlerResult.UserBio

	c.JSON(http.StatusOK, response)
}
