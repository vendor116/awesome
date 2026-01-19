package web

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"strconv"

	"github.com/vendor116/awesome/internal/config"
	"golang.org/x/sync/errgroup"
)

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
