package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/autoupdater"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler_handler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/gin-gonic/gin"
)

type GetItemRequest struct {
	Identity   string               `form:"proof" binding:"required"`
	PlatformID constants.PlatformID `form:"platform_id" binding:"required"`
	NetworkID  constants.NetworkID  `form:"network_id"`
	ItemType   constants.ItemType   `form:"item_type"`
	Limit      int                  `form:"limit"`
	Timestamp  int64                `form:"timestamp"`
}

type itemsResult struct {
	NoteItems  *[]model.Item `json:"note"`
	AssetItems *[]model.Item `json:"asset"`
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

	if request.ItemType == "" {
		paramErrMsg += "item_type is empty; "
	}

	if request.ItemType != constants.ItemTypeAsset &&
		request.ItemType != constants.ItemTypeNote {
		paramErrMsg += "itemtype should be one of 'note' or 'asset'; "
	}

	if paramErrMsg != "" {
		response.ErrorBase = util.GetErrorBase(util.ErrorCodeParameterError)
		response.ErrorBase.ErrorMsg += ": " + util.ErrorMsg(paramErrMsg)

		c.JSON(http.StatusOK, response)

		return
	}

	// get items from db
	dbResult, err := getItemsFromDB(c.Request.Context(), request)
	if err != nil {
		logger.Errorf("get items from db error: %s", err.Error())
	}

	if dbResult != nil {
		response.ItemsResult = *dbResult

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

func getItemsFromDB(context context.Context, request GetItemRequest) (*itemsResult, error) {
	ai := rss3uri.NewAccountInstance(request.Identity, request.PlatformID.Symbol())

	var err error

	var result = new(itemsResult)

	isExisted, err := db.Exists(ai)
	if err != nil {
		return nil, fmt.Errorf("find db exists false:%s", err)
	}

	addToRecentVisit(context, &request)

	if !isExisted {
		return nil, nil
	}

	switch request.ItemType {
	case constants.ItemTypeNote:
		result.NoteItems, err = db.GetAccountItems(ai, constants.ItemTypeNote)
		if err != nil {
			return nil, fmt.Errorf("get notes items error:%s", err)
		}
	case constants.ItemTypeAsset:
		result.AssetItems, err = db.GetAccountItems(ai, constants.ItemTypeAsset)
		if err != nil {
			return nil, fmt.Errorf("get asset items error:%s", err)
		}
	default:
		result.NoteItems, err = db.GetAccountItems(ai, constants.ItemTypeNote)
		if err != nil {
			return nil, fmt.Errorf("get notes items error:%s", err)
		}

		result.AssetItems, err = db.GetAccountItems(ai, constants.ItemTypeAsset)
		if err != nil {
			return nil, fmt.Errorf("get asset items error:%s", err)
		}
	}

	return result, err
}

func getItemsFromCrawlerHandler(crawlerResult []*model.Item, itemType constants.ItemType) *itemsResult {
	result := new(itemsResult)

	switch itemType {
	case constants.ItemTypeNote:
		result.NoteItems = crawlerResult2ItemsResult(crawlerResult)
	case constants.ItemTypeAsset:
		result.AssetItems = crawlerResult2ItemsResult(crawlerResult)
	default:
		result.NoteItems = crawlerResult2ItemsResult(crawlerResult)
		result.AssetItems = crawlerResult2ItemsResult(crawlerResult)
	}

	return result
}

func crawlerResult2ItemsResult(itemsPointArr []*model.Item) *[]model.Item {
	itemsArrPoint := new([]model.Item)
	for _, itemsPoint := range itemsPointArr {
		*itemsArrPoint = append(*itemsArrPoint, *itemsPoint)
	}

	return itemsArrPoint
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
	itemType constants.ItemType,
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

	return getItemsFromCrawlerHandler(handlerResult.Result.Items, itemType), util.GetErrorBase(util.ErrorCodeSuccess)
}

func setItemsResult(originalItems *[]model.Item, addedItems *[]model.Item) {
	if originalItems == nil {
		originalItems = new([]model.Item)
	}

	if addedItems != nil {
		curraddedItems := *addedItems
		*originalItems = append(*originalItems, curraddedItems...)
	}
}

func getItemsResult(request GetItemRequest) (*itemsResult, util.ErrorBase) {
	result := new(itemsResult)
	errorBase := util.GetErrorBase(util.ErrorCodeSuccess)

	if request.NetworkID == constants.NetworkIDUnknown {
		networkIDs := constants.GetEthereumPlatformNetworks()
		for _, networkID := range networkIDs {
			currResult, currErrorBase := getItemsResultFromOneNetwork(
				request.Identity, request.PlatformID, networkID, request.ItemType,
				request.Limit, time.Unix(request.Timestamp, 0),
			)

			if currErrorBase.ErrorCode != util.ErrorCodeSuccess {
				logger.Errorf("[%s] get item error, network[%s],error reason:%s",
					request.Identity, networkID.Symbol(), currErrorBase.ErrorMsg)
			}

			setItemsResult(result.AssetItems, currResult.AssetItems)
			setItemsResult(result.NoteItems, currResult.NoteItems)
		}
	} else {
		result, errorBase = getItemsResultFromOneNetwork(
			request.Identity, request.PlatformID, request.NetworkID, request.ItemType,
			request.Limit, time.Unix(request.Timestamp, 0),
		)
	}

	return result, errorBase
}
