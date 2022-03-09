package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetAssetListRequest struct {
	Instance  string `uri:"instance" binding:"required"`
	PageIndex string `uri:"page_index" binding:"required"`
}

func GetAssetListRequestHandlerFunc(c *gin.Context) {
	request := GetAssetListRequest{}
	if err := c.ShouldBindUri(&request); err != nil {
		return
	}

	// TODO Query data from database
	// TODO No test data available
	c.JSON(http.StatusOK, request)
}
