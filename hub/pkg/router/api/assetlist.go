package api

import (
	"fmt"
	"net/http"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/protocol/file"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/status"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/web"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/gin-gonic/gin"
)

type GetAssetListRequest struct {
	PageIndex int `uri:"page_index"`
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

	identifier := rss3uri.New(platformInstance).String()

	assetListFile := file.AssetList{
		SignedBase: protocol.SignedBase{
			Base: protocol.Base{
				Version:    protocol.Version,
				Identifier: fmt.Sprintf("%s/list/asset/%d", identifier, request.PageIndex),
			},
		},
	}

	// TODO No test data available

	assetListFile.Total = len(assetListFile.List)

	c.JSON(http.StatusOK, &assetListFile)
}
