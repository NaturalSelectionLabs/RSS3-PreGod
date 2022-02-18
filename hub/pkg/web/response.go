package web

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/status"
	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code status.Code `json:"code"`
	Msg  status.Msg  `json:"message"`
	Data interface{} `json:"data"`
}

func (g *Gin) JSONResponse(httpCode int, errCode status.Code, data interface{}) {
	g.C.JSON(httpCode, Response{
		Code: errCode,
		Msg:  status.GetMsg(errCode),
		Data: data,
	})
}
