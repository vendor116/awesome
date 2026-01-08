package main

import (
	"context"
	"errors"
	"log/slog"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/vendor116/awesome/internal"
	"github.com/vendor116/awesome/internal/config"
	"golang.org/x/sync/errgroup"
)

var version = "dev"

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Default().Error("failed to load config file", "error", err)
		os.Exit(1)
	}

	err = internal.SetLogger(cfg.LogLevel, version)
	if err != nil {
		slog.Default().Error("failed to set up json logger", "error", err)
		os.Exit(1)
	}

	eg, egCtx := errgroup.WithContext(context.Background())
	ctx, cancel := signal.NotifyContext(egCtx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	internal.RunRESTServer(ctx, eg, cfg.Rest)
	internal.RunGrpcServer(ctx, eg, cfg.Grpc)

	if cfg.Pprof.Enabled {
		internal.RunPprofServer(ctx, eg, cfg.Pprof)
	}

	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		slog.Default().Error("application completed", "error", err)
	}
}
