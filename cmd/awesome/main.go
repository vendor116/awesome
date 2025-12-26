package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/vendor116/awesome/internal/config"
	"github.com/vendor116/awesome/internal/web/grpc"
	"github.com/vendor116/awesome/internal/web/rest"
	"golang.org/x/sync/errgroup"
)

var (
	version      = "dev"
	path, prefix string
)

func main() {
	cfg, err := config.Load[config.App](path, prefix)
	if err != nil {
		slog.Default().Error("failed to load config file", "error", err)
		os.Exit(1)
	}

	err = setDefaultLogger(cfg.LogLevel, version)
	if err != nil {
		slog.Default().Error("failed to set up json logger", "error", err)
		os.Exit(1)
	}

	eg, egCtx := errgroup.WithContext(context.Background())
	ctx, cancel := signal.NotifyContext(egCtx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	rest.RunServer(ctx, eg, cfg.RESTServer)
	grpc.RunServer(ctx, eg, cfg.GRPCServer)

	if cfg.PprofEnable {
		rest.RunPprofServer(ctx, eg, cfg.PprofServer)
	}

	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		slog.Default().Error("application completed", "error", err)
	}
}

func init() {
	flag.StringVar(&path, "config", "", "path to config file")
	flag.StringVar(&prefix, "prefix", "", "environment variable prefix")
	flag.Parse()
}

func setDefaultLogger(level, version string) error {
	var l slog.Level
	if err := l.UnmarshalText([]byte(level)); err != nil {
		return err
	}

	slog.SetDefault(
		slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: l,
			}),
		).With(
			slog.String("version", version),
		),
	)
	return nil
}
