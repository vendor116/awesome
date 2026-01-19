package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/vendor116/awesome/internal"
	"github.com/vendor116/awesome/internal/config"
	"github.com/vendor116/awesome/internal/web"
	"github.com/vendor116/awesome/internal/web/grpc/awesome"
	"github.com/vendor116/awesome/internal/web/rest/router"
	v1 "github.com/vendor116/awesome/internal/web/rest/v1"
	awesomepb "github.com/vendor116/awesome/pkg/protobuf/awesome"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

var version = "dev"

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Default().Error("failed to load config file", "error", err)
		os.Exit(1)
	}

	err = internal.SetupLog(cfg.LogLevel, version)
	if err != nil {
		slog.Default().Error("failed to set up json logger", "error", err)
		os.Exit(1)
	}

	v1RestServer := router.AttachHandlers(v1.NewServer())

	grpcServer := grpc.NewServer()
	awesomepb.RegisterAwesomeServer(grpcServer, awesome.NewServer())

	eg, egCtx := errgroup.WithContext(context.Background())
	ctx, cancel := signal.NotifyContext(egCtx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	web.RunRESTServer(ctx, eg, v1RestServer, cfg.Rest)
	web.RunGrpcServer(ctx, eg, grpcServer, cfg.Grpc)

	if cfg.Pprof.Enabled {
		web.RunPprofServer(ctx, eg, cfg.Pprof)
	}

	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		slog.Default().Error("application completed", "error", err)
	}
}
