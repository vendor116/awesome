package httpserver

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/vendor116/awesome/internal/http_server/handlers"
	"github.com/vendor116/awesome/internal/http_server/middleware/rlog"
	"github.com/vendor116/awesome/pkg/openapi"
)

func RegisterHandlers(ssi openapi.StrictServerInterface) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(
		gin.Recovery(),
		rlog.Middleware(slog.Default()),
	)

	r.GET("/health", handlers.HealthHandler)

	return registerAPIV1(r, ssi)
}

func registerAPIV1(r *gin.Engine, ssi openapi.StrictServerInterface) *gin.Engine {
	v1 := r.Group("/api/v1")

	openapi.RegisterHandlers(v1, openapi.NewStrictHandler(ssi, nil))

	return r
}
