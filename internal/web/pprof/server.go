package pprof

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/vendor116/awesome/internal/config"
	httpserver "github.com/vendor116/awesome/internal/web/http_server"
	"golang.org/x/sync/errgroup"
)

func RunServer(
	ctx context.Context,
	g *errgroup.Group,
	cfg config.HTTPServer,
	loggerFactory func() *slog.Logger,
) {
	httpserver.Run(
		ctx,
		g,
		cfg.GetAddress(),
		http.DefaultServeMux,
		loggerFactory(),
	)
}
