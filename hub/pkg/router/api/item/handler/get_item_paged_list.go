package item

import (
	"net/http"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/status"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/web"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/gin-gonic/gin"
)

type GetItemPagedListRequestUri struct {
	Authority string             `uri:"authority" binding:"required"`
	ItemType  constants.ItemType `uri:"item_type" binding:"required"`
	PageIndex int                `uri:"page_index" binding:"required"`
}

type GetItemPagedListResponseData struct {
	Authority rss3uri.Instance `json:"authority"`
}

func GetItemPagedList(c *gin.Context) {
	w := web.Gin{C: c}

	w.JSONResponse(http.StatusOK, status.SUCCESS, gin.H{})
}
