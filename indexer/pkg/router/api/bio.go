package api

import (
	"net/http"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/gin-gonic/gin"
)

type GetBioRequest struct {
	Identity   string               `form:"proof" binding:"required"`
	PlatformId constants.PlatformID `form:"platform_id" binding:"required"`
}

func GetBioHandlerFunc(c *gin.Context) {
	request := GetBioRequest{}
	if err := c.ShouldBind(&request); err != nil {
		logger.Errorf("%s", err.Error())

		return
	}

	// TODO

	c.JSON(http.StatusOK, request)
}
