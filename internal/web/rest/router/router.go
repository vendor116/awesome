package router

import (
	"github.com/gin-gonic/gin"
	"github.com/vendor116/awesome/internal/web/rest/middleware/logger"
	v1 "github.com/vendor116/awesome/internal/web/rest/v1"
	openapiv1 "github.com/vendor116/awesome/pkg/openapi/v1"
)

func AttachHandlers(s openapiv1.StrictServerInterface) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(
		gin.Recovery(),
		logger.Middleware(),
	).
		GET("/health", v1.HealthHandler)

	apiV1 := router.Group("/api/v1")
	openapiv1.RegisterHandlers(apiV1, openapiv1.NewStrictHandler(s, nil))

	return router
}
