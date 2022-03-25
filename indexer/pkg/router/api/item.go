package api

import (
	"fmt"
	"net/http"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/gin-gonic/gin"
)

type GetItemRequest struct {
	Identity   string               `form:"proof" binding:"required"`
	PlatformID constants.PlatformID `form:"platform_id" binding:"required"`
	NetworkID  constants.NetworkID  `form:"network_id"`
	ItemType   constants.ItemType   `form:"item_type"`
}

func makeGinResponse(ai rss3uri.Instance, ItemType constants.ItemType) interface{} {
	var res *[]model.Item

	var err error

	switch ItemType {
	case constants.ItemTypeNote:
		res, err = db.GetNotes(ai)
	case constants.ItemTypeAsset:
		res, err = db.GetAssets(ai)
	default:
		err = fmt.Errorf("unsupported item type")
	}

	if err != nil {
		return gin.H{
			"error":   0, // TODO: consistent with bio api
			"message": err.Error(),
		}
	}

	return gin.H{
		"result": res,
	}
}

func GetItemHandlerFunc(c *gin.Context) {
	request := GetItemRequest{}
	if err := c.ShouldBind(&request); err != nil {
		return
	}

	// TODO Query data
	// 旧用户查询数据库，新用户拉取
	ai := rss3uri.NewAccountInstance(request.Identity, request.PlatformID.Symbol())

	isOld, err := db.Exists(ai)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error":   0, // TODO: consistent with bio api
			"message": err.Error(),
		})

		return
	}

	if isOld {
		switch request.Identity {
		case string(constants.ItemTypeNote):
			obj := makeGinResponse(ai, constants.ItemTypeNote)

			c.JSON(http.StatusOK, gin.H{
				string(constants.ItemTypeNote): obj,
			})
		case string(constants.ItemTypeAsset):
			obj := makeGinResponse(ai, constants.ItemTypeAsset)

			c.JSON(http.StatusOK, gin.H{
				string(constants.ItemTypeAsset): obj,
			})
		default:
			obj1 := makeGinResponse(ai, constants.ItemTypeNote)
			obj2 := makeGinResponse(ai, constants.ItemTypeAsset)

			c.JSON(http.StatusOK, gin.H{
				string(constants.ItemTypeNote):  obj1,
				string(constants.ItemTypeAsset): obj2,
			})
		}
	} else {
		//TODO: Work() on this task
		c.JSON(http.StatusOK, request)
	}
}
