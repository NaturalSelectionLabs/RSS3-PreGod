package api

import (
	"net/http"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/status"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/web"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/gin-gonic/gin"
)

type GetAssetListRequest struct {
	PageIndex string `uri:"page_index" binding:"required"`
}

func GetAssetListRequestHandlerFunc(c *gin.Context) {
	request := GetAssetListRequest{}
	if err := c.ShouldBindUri(&request); err != nil {
		return
	}

	value, exists := c.Get(middleware.KeyInstance)
	if !exists {
		w := web.Gin{C: c}
		w.JSONResponse(http.StatusBadRequest, status.CodeInvalidParams, nil)

		return
	}

	platformInstance, ok := value.(*rss3uri.PlatformInstance)
	if !ok {
		w := web.Gin{C: c}
		w.JSONResponse(http.StatusBadRequest, status.CodeInvalidParams, nil)

		return
	}

	if platformInstance.Prefix != constants.PrefixNameAccount || platformInstance.Platform != constants.PlatformSymbolEthereum {
		w := web.Gin{C: c}
		w.JSONResponse(http.StatusBadRequest, status.CodeInvalidParams, nil)

		return
	}

	// TODO Query data from database
	// TODO No test data available
	c.JSON(http.StatusOK, request)
}
