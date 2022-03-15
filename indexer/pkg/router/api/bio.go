package api

import (
	"net/http"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/gin-gonic/gin"
)

type GetBioRequest struct {
	Identity   string               `form:"proof" binding:"required"`
	PlatformId constants.PlatformID `form:"platform_id" binding:"required"`
}

func GetBioHandlerFunc(c *gin.Context) {
	request := GetItemRequest{}
	if err := c.ShouldBind(&request); err != nil {
		return
	}

	// TODO Query data

	c.JSON(http.StatusOK, request)
}
