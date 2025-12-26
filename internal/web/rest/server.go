package rest

import (
	"context"
	"log/slog"

	"github.com/vendor116/awesome/internal/config"
	httpserver "github.com/vendor116/awesome/internal/web/http_server"
	v1 "github.com/vendor116/awesome/pkg/rest/v1"
	"golang.org/x/sync/errgroup"
)

func RunServer(
	ctx context.Context,
	g *errgroup.Group,
	h v1.StrictServerInterface,
	cfg config.HTTPServer,
	loggerFactory func() *slog.Logger,
) {
	httpserver.Run(
		ctx,
		g,
		cfg.GetAddress(),
		registerHandlers(h),
		loggerFactory(),
	)
}
