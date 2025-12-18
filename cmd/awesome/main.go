package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/vendor116/awesome/internal"
	"github.com/vendor116/awesome/internal/config"
	httpserver "github.com/vendor116/awesome/internal/http_server"
	"github.com/vendor116/awesome/internal/http_server/handlers"
	"github.com/vendor116/awesome/pkg/version"
	"golang.org/x/sync/errgroup"
)

func main() {
	cfg, err := config.LoadAndValidate()
	if err != nil {
		log.Fatalf("failed to load config:%v", err)
	}

	err = internal.SetJSONLogger(cfg.LogLevel)
	if err != nil {
		log.Fatalf("failed to set json logger:%v", err)
	}

	slog.Default().Info("starting awesome", slog.String("version", version.GetVersion()))

	vh := handlers.NewBaseHandler()

	eg, egCtx := errgroup.WithContext(context.Background())
	ctx, cancel := signal.NotifyContext(egCtx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	httpserver.Run(
		ctx,
		eg,
		httpserver.RegisterHandlers(vh),
		cfg.HTTPServer,
	)

	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		slog.Default().Error("application completed", "error", err)
	}
}
