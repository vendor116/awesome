package internal

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"strconv"

	"github.com/vendor116/awesome/internal/config"
	"github.com/vendor116/awesome/internal/web/grpc/awesome"
	"github.com/vendor116/awesome/internal/web/rest/router"
	v1 "github.com/vendor116/awesome/internal/web/rest/v1"
	awesomepb "github.com/vendor116/awesome/pkg/protobuf/awesome"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunGrpcServer(ctx context.Context, g *errgroup.Group, cfg config.GrpcConfig) {
	s := grpc.NewServer()

	awesomepb.RegisterAwesomeServer(s, awesome.NewServer())

	if cfg.Reflect {
		reflection.Register(s)
	}

	g.Go(func() error {
		lc := &net.ListenConfig{}

		listen, err := lc.Listen(ctx, "tcp", net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)))
		if err != nil {
			return err
		}

		slog.Default().Info("starting grpc server", "address", listen.Addr().String())
		err = s.Serve(listen)
		if err != nil {
			return err
		}

		slog.Default().Info("grpc server stopped")
		return nil
	})
	g.Go(func() error {
		<-ctx.Done()
		s.GracefulStop()
		return nil
	})
}

func RunRESTServer(ctx context.Context, g *errgroup.Group, cfg config.RESTConfig) {
	v1Server := v1.NewServer()

	s := http.Server{
		Addr:              net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		Handler:           router.AttachHandlers(v1Server),
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	g.Go(func() error {
		slog.Default().InfoContext(ctx, "starting rest server", "address", s.Addr)
		err := s.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		slog.Default().InfoContext(ctx, "rest server stopped")
		return nil
	})
	g.Go(func() error {
		<-ctx.Done()

		sdCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		err := s.Shutdown(ctx)
		if err != nil {
			slog.Default().ErrorContext(sdCtx, "failed to shutdown rest server", "error", err)
		}
		return nil
	})
}

func RunPprofServer(ctx context.Context, g *errgroup.Group, cfg config.PprofConfig) {
	s := http.Server{
		Addr:              net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		Handler:           http.DefaultServeMux,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	g.Go(func() error {
		slog.Default().InfoContext(ctx, "starting pprof server", "address", s.Addr)
		err := s.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		slog.Default().InfoContext(ctx, "pprof server stopped")
		return nil
	})
	g.Go(func() error {
		<-ctx.Done()

		sdCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		err := s.Shutdown(ctx)
		if err != nil {
			slog.Default().ErrorContext(sdCtx, "failed to shutdown pprof server", "error", err)
		}
		return nil
	})
}
