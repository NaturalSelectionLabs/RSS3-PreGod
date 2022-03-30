package api

import (
	"context"
	"net/http"

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
	NetworkID  constants.NetworkID  `form:"network_id" binding:"required"`
	ItemType   constants.ItemType   `form:"item_type" binding:"required"`
}

type itemsResult struct {
	NoteItems  *[]model.Item `json:"note"`
	AssetItems *[]model.Item `json:"asset"`
}

type GetItemResponse struct {
	util.ErrorBase `json:"error"`
	ItemsResult    itemsResult
}

func GetItemHandlerFunc(c *gin.Context) {
	request := GetItemRequest{}
	if err := c.ShouldBind(&request); err != nil {
		logger.Errorf("request bind error: %s", err.Error())

		return
	}

	response := GetItemResponse{
		util.ErrorBase{},
		itemsResult{},
	}

	if len(request.Identity) <= 0 ||
		!constants.IsValidPlatformSymbol(string(request.PlatformID.Symbol())) ||
		!constants.IsValidNetworkName(string(request.NetworkID.Symbol())) ||
		request.ItemType == constants.ItemTypeUnknown {
		logger.Errorf("parameter error")

		response.ErrorBase = util.GetErrorBase(util.ErrorCodeParameterError)
		c.JSON(http.StatusOK, response)
	}

	dbResult, err := getItemsFromDB(c, request)
	if err != nil {
		logger.Errorf("get items from db error: %s", err.Error())
	}

	if dbResult != nil {
		response.ItemsResult = *dbResult

		c.JSON(http.StatusOK, response)
		return
	}

	getItemHandler := crawler_handler.NewGetItemsHandler(crawler.WorkParam{
		Identity:   request.Identity,
		PlatformID: request.PlatformID,
		NetworkID:  platform2Network[request.PlatformID],
	})
	handlerResult := getItemHandler.Excute()

	if handlerResult == nil || handlerResult.Result == nil {
		logger.Errorf("[%s] get item error", request.Identity)

		response.ErrorBase = util.GetErrorBase(util.ErrorCodeNotFoundData)
		c.JSON(http.StatusOK, response)

		return
	}

	if handlerResult.Error.ErrorCode != util.ErrorCodeSuccess {
		logger.Errorf("[%s] get item error", request.Identity)

		response.ErrorBase = util.GetErrorBase(handlerResult.Error.ErrorCode)
		c.JSON(http.StatusOK, response)

		return
	}

	response.ItemsResult = *getItemsFromCrawlerHandler(handlerResult.Result.Items, request.ItemType)

	c.JSON(http.StatusOK, response)
}

func getItemsFromDB(c *gin.Context, request GetItemRequest) (*itemsResult, error) {
	ai := rss3uri.NewAccountInstance(request.Identity, request.PlatformID.Symbol())

	var err error

	var result = new(itemsResult)

	isOld, err := db.Exists(ai)
	if err != nil {
		return nil, err
	}

	addToRecentVisit(c.Request.Context(), &request)

	if !isOld {
		return nil, nil
	}

	switch request.ItemType {
	case constants.ItemTypeNote:
		result.NoteItems, err = db.GetAccountItems(ai, constants.ItemTypeNote)
		if err != nil {
			return nil, err
		}
	case constants.ItemTypeAsset:
		result.AssetItems, err = db.GetAccountItems(ai, constants.ItemTypeAsset)
		if err != nil {
			return nil, err
		}
	default:
		result.NoteItems, err = db.GetAccountItems(ai, constants.ItemTypeNote)
		if err != nil {
			return nil, err
		}

		result.AssetItems, err = db.GetAccountItems(ai, constants.ItemTypeAsset)
		if err != nil {
			return nil, err
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
		// Limit:      ?,
		// TimeStamp:  ?,
	}

	return autoupdater.AddToRecentVisitQueue(ctx, param)
}
