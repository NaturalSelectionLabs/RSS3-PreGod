package instance

import (
	"net/http"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/status"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/web"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/gin-gonic/gin"
)

type GetInstanceRequestUri struct {
	Authority string `uri:"authority" binding:"required"`
}

type GetInstanceResponseData struct {
	Authority rss3uri.Instance `json:"authority"`
}

// GetInstance returns the instance information for the given authority.
//
// @Summary      Get instance information
// @Description  get instance information by authority
// @Tags         authority
// @Accept       json
// @Produce      json
// @Param        authority  path      string  true  "Authority"
// @Success      200        {object}  web.Response{data=GetInstanceResponseData}
// @Router       /{authority} [get]
func GetInstance(c *gin.Context) {
	w := web.Gin{C: c}

	// validate uri
	var uri GetInstanceRequestUri
	if err := c.ShouldBindUri(&uri); err != nil {
		w.JSONResponse(http.StatusBadRequest, status.INVALID_PARAMS, "invalid uri")

		return
	}

	// parse uri
	authority, err := rss3uri.ParseInstance(uri.Authority)
	if err != nil {
		w.JSONResponse(http.StatusBadRequest, status.INVALID_PARAMS, "invalid uri: "+err.Error())

		return
	}

	// TODO: get instance from db

	w.JSONResponse(http.StatusOK, status.SUCCESS, gin.H{
		"authority": authority,
	})
}
