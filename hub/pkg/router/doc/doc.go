package doc

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/docs"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func Doc(basePath string, version string) gin.HandlerFunc {
	docs.SwaggerInfo_swagger.Title = "RSS3-Hub API"
	docs.SwaggerInfo_swagger.Description = "RSS3-Hub API"
	docs.SwaggerInfo_swagger.Version = version
	docs.SwaggerInfo_swagger.BasePath = basePath
	docs.SwaggerInfo_swagger.Schemes = []string{"http", "https"}

	return ginSwagger.WrapHandler(swaggerFiles.Handler)
}
