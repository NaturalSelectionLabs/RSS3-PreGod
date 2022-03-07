package item

import (
	"net/http"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/rss3uri"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/gin-gonic/gin"
)

type GetItemRequestUri struct {
	Authority string             `uri:"authority" binding:"required"`
	ItemType  constants.ItemType `uri:"item_type" binding:"required"`
	ItemUUID  string             `uri:"item_uuid" binding:"required"`
}

type GetItemResponseData struct {
	Authority rss3uri.Instance `json:"authority"`
}

func GetItem(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
