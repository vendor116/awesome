package rest

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/vendor116/awesome/internal/web/rest/middleware/logger"
	restv1 "github.com/vendor116/awesome/internal/web/rest/v1"
	v1 "github.com/vendor116/awesome/pkg/rest/v1"
)

func registerHandlers(server v1.StrictServerInterface) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(
		gin.Recovery(),
		logger.Middleware(slog.Default().With(slog.String("server", "rest"))),
	).
		GET("/health", restv1.HealthHandler)

	apiV1 := router.Group("/api/v1")
	v1.RegisterHandlers(apiV1, v1.NewStrictHandler(server, nil))

	return router
}
