package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"log/slog"
	_ "net/http/pprof"
	"os/signal"
	"syscall"

	"github.com/vendor116/awesome/internal"
	"github.com/vendor116/awesome/internal/config"
	"github.com/vendor116/awesome/internal/web/grpc"
	"github.com/vendor116/awesome/internal/web/grpc/awesome"
	"github.com/vendor116/awesome/internal/web/pprof"
	"github.com/vendor116/awesome/internal/web/rest"
	v1 "github.com/vendor116/awesome/internal/web/rest/v1"
	"golang.org/x/sync/errgroup"
)

var version = "dev"

func main() {
	var path, prefix string
	flag.StringVar(&path, "config", "", "path to config file")
	flag.StringVar(&prefix, "prefix", "", "environment variable prefix")
	flag.Parse()

	cfg, err := config.Load[config.App](path, prefix)
	if err != nil {
		log.Fatalf("failed to load config:%v", err)
	}

	err = cfg.Validate()
	if err != nil {
		log.Fatalf("invalid config:%v", err)
	}

	err = internal.SetJSONLogger(cfg.LogLevel, version)
	if err != nil {
		log.Fatalf("failed to set json logger:%v", err)
	}

	restV1Server := v1.NewServer()
	grpcAwesomeServer := awesome.NewServer()

	eg, egCtx := errgroup.WithContext(context.Background())
	ctx, cancel := signal.NotifyContext(egCtx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	rest.RunServer(
		ctx,
		eg,
		restV1Server,
		cfg.RESTServer,
	)

	grpc.RunServer(
		ctx,
		eg,
		grpcAwesomeServer,
		cfg.GRPCServer,
	)

	if cfg.PprofEnabled {
		pprof.RunServer(
			ctx,
			eg,
			cfg.PprofServer,
		)
	}

	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		slog.Default().Error("application completed", "error", err)
	}
}
