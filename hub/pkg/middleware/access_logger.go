package middleware

import (
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var desugarredLogger = logger.DesugarredLogger

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		c.Next()

		latency := time.Since(t)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		path := c.Request.URL.Path

		switch {
		case statusCode >= 400 && statusCode <= 499:
			{
				desugarredLogger.Warn("[ACCESS]",
					zap.Int("statusCode", statusCode),
					zap.String("latency", latency.String()),
					zap.String("clientIP", clientIP),
					zap.String("method", method),
					zap.String("path", path),
					zap.String("error", c.Errors.String()),
				)
			}
		case statusCode >= 500:
			{
				desugarredLogger.Error("[ACCESS]",
					zap.Int("statusCode", statusCode),
					zap.String("latency", latency.String()),
					zap.String("clientIP", clientIP),
					zap.String("method", method),
					zap.String("path", path),
					zap.String("error", c.Errors.String()),
				)
			}
		default:
			desugarredLogger.Info("[ACCESS]",
				zap.Int("statusCode", statusCode),
				zap.String("latency", latency.String()),
				zap.String("clientIP", clientIP),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("error", c.Errors.String()),
			)
		}
	}
}
