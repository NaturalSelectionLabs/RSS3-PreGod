package web

import (
	"net/http"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/status"
	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code    status.Code    `json:"code"`
	Message status.Message `json:"message"`
	Data    interface{}    `json:"data,omitempty"`
}

func (g *Gin) JSONResponse(httpCode int, errCode status.Code, data interface{}) {
	if httpCode == http.StatusOK {
		g.C.JSON(httpCode, data)

		return
	}

	g.C.JSON(httpCode, Response{
		Code:    errCode,
		Message: errCode.Message(),
		Data:    data,
	})
}
