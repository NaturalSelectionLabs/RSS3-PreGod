package ping

import (
	"net/http"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/web"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/status"
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	w := web.Gin{C: c}

	w.JSONResponse(http.StatusOK, status.CodeOK, "pong")
}
