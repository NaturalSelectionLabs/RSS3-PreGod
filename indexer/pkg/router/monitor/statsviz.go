package monitor

import (
	"github.com/arl/statsviz"
	"github.com/gin-gonic/gin"
)

func Statsviz(c *gin.Context) {
	if c.Param("filepath") == "/ws" {
		statsviz.Ws(c.Writer, c.Request)

		return
	}

	statsviz.IndexAtRoot("/debug/statsviz").ServeHTTP(c.Writer, c.Request)
}
