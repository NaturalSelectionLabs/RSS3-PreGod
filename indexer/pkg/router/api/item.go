package api

import (
	"context"
	"net/http"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/autoupdater"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler_handler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/gin-gonic/gin"
)

type GetItemRequest struct {
	Identity   string               `form:"proof" binding:"required"`
	PlatformID constants.PlatformID `form:"platform_id" binding:"required"`
	NetworkID  constants.NetworkID  `form:"network_id"`
	Limit      int                  `form:"limit"`
	Timestamp  int64                `form:"timestamp"`
}

type itemsResult struct {
	NoteItems  []model.Note  `json:"note"`
	AssetItems []model.Asset `json:"asset"`
}

type GetItemResponse struct {
	util.ErrorBase `json:"error"`
	ItemsResult    itemsResult `json:"data"`
}

func GetItemHandlerFunc(c *gin.Context) {
	request := GetItemRequest{}
	response := GetItemResponse{
		util.ErrorBase{},
		itemsResult{},
	}

	// bind request
	if err := c.ShouldBind(&request); err != nil {
		logger.Errorf("request bind error: %s", err.Error())

		return
	}

	// set default request
	if request.Limit == 0 {
		request.Limit = 100 // TODO: constants.DefaultLimit?
	}

	if request.Timestamp == 0 {
		request.Timestamp = time.Now().Unix()
	}

	// request validation
	var paramErrMsg string

	if request.Identity == "" {
		paramErrMsg = "identity is empty; "
	}

	if !constants.IsValidPlatformSymbol(string(request.PlatformID.Symbol())) {
		paramErrMsg += "platform_id is invalid; "
	}

	if !constants.IsValidNetworkName(string(request.NetworkID.Symbol())) {
		paramErrMsg += "network_id is invalid; "
	}

	if paramErrMsg != "" {
		response.ErrorBase = util.GetErrorBase(util.ErrorCodeParameterError)
		response.ErrorBase.ErrorMsg += ": " + util.ErrorMsg(paramErrMsg)

		c.JSON(http.StatusOK, response)

		return
	}

	// get items from crawler
	result, errorBase := getItemsResult(request)
	response.ErrorBase = errorBase

	if response.ErrorBase.ErrorCode == 0 {
		response.ItemsResult = *result

		c.JSON(http.StatusOK, response)
	}

	if response.ErrorBase.ErrorCode != util.ErrorCodeSuccess {
		logger.Errorf("[%s] get item error", request.Identity)

		response.ErrorBase = util.GetErrorBase(response.ErrorBase.ErrorCode)
		c.JSON(http.StatusOK, response)
	}
}

func addToRecentVisit(ctx context.Context, req *GetItemRequest) error {
	param := &crawler.WorkParam{
		Identity:   req.Identity,
		NetworkID:  req.NetworkID,
		PlatformID: req.PlatformID,
		// NOTE looks like only for misskey
		Limit:     req.Limit,
		Timestamp: time.Unix(req.Timestamp, 0),
	}

	return autoupdater.AddToRecentVisitQueue(ctx, param)
}

func getItemsResultFromOneNetwork(identity string,
	platformID constants.PlatformID,
	networkID constants.NetworkID,
	limit int,
	Timestamp time.Time,
) (*itemsResult, util.ErrorBase) {
	getItemHandler := crawler_handler.NewGetItemsHandler(crawler.WorkParam{
		Identity:   identity,
		PlatformID: platformID,
		NetworkID:  networkID,
		Limit:      limit,
		Timestamp:  Timestamp,
	})

	handlerResult, err := getItemHandler.Excute()
	if err != nil {
		logger.Errorf("get items from crawler error: %s", err.Error())

		return nil, util.GetErrorBase(util.ErrorCodeNotFoundData)
	}

	if handlerResult == nil || handlerResult.Result == nil {
		return nil, util.GetErrorBase(util.ErrorCodeNotFoundData)
	}

	if handlerResult.Error.ErrorCode != util.ErrorCodeSuccess {
		logger.Errorf("[%s] get item error", identity)

		return nil, util.GetErrorBase(handlerResult.Error.ErrorCode)
	}

	return &itemsResult{
		NoteItems:  handlerResult.Result.Notes,
		AssetItems: handlerResult.Result.Assets,
	}, util.GetErrorBase(util.ErrorCodeSuccess)
}

func getItemsResult(request GetItemRequest) (*itemsResult, util.ErrorBase) {
	result := new(itemsResult)
	result.AssetItems = make([]model.Asset, 0)
	result.NoteItems = make([]model.Note, 0)
	errorBase := util.GetErrorBase(util.ErrorCodeSuccess)

	if request.NetworkID == constants.NetworkIDUnknown {
		networkIDs := constants.GetEthereumPlatformNetworks()
		for _, networkID := range networkIDs {
			currResult, currErrorBase := getItemsResultFromOneNetwork(
				request.Identity, request.PlatformID, networkID,
				request.Limit, time.Unix(request.Timestamp, 0),
			)

			if currErrorBase.ErrorCode != util.ErrorCodeSuccess {
				logger.Errorf("[%s] get item error, network[%s],error reason:%s",
					request.Identity, networkID.Symbol(), currErrorBase.ErrorMsg)
			}

			result.AssetItems = append(result.AssetItems, currResult.AssetItems...)
			result.NoteItems = append(result.NoteItems, currResult.NoteItems...)
		}
	} else {
		result, errorBase = getItemsResultFromOneNetwork(
			request.Identity, request.PlatformID, request.NetworkID,
			request.Limit, time.Unix(request.Timestamp, 0),
		)
	}

	return result, errorBase
}
