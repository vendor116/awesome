package rlog

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Middleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		logAttrs := []any{
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", status,
			"latency_ms", latency.Milliseconds(),
			"latency", latency.String(),
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		}

		if size := c.Writer.Size(); size > 0 {
			logAttrs = append(logAttrs, "response_size", size)
		}

		if errors := c.Errors.ByType(gin.ErrorTypePrivate).String(); errors != "" {
			logAttrs = append(logAttrs, "errors", errors)
		}

		switch {
		case status >= http.StatusInternalServerError:
			logger.Error("HTTP request error", logAttrs...)
		case status >= http.StatusBadRequest:
			logger.Warn("HTTP client error", logAttrs...)
		default:
			logger.Debug("HTTP request", logAttrs...)
		}
	}
}
