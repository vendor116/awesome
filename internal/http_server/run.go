package httpserver

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/vendor116/awesome/internal/config"
	"golang.org/x/sync/errgroup"
)

const shutdownTimeout = 3 * time.Second

func Run(ctx context.Context, g *errgroup.Group, handler http.Handler, cfg config.HTTPServer) {
	server := &http.Server{
		Addr:        net.JoinHostPort(cfg.Host, cfg.Port),
		Handler:     handler,
		ReadTimeout: cfg.ReadHeaderTimeout,
	}

	logger := slog.Default().With("address", server.Addr)

	g.Go(func() error {
		logger.Info("starting http server")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		logger.Info("http server shutdown gracefully")
		return nil
	})

	g.Go(func() error {
		<-ctx.Done()
		if err := context.Cause(ctx); err != nil && !errors.Is(err, context.Canceled) {
			return err
		}

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		logger.Info("shutting down http server")
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.ErrorContext(shutdownCtx, "failed to shutdown http server", "error", err)
		}
		return nil
	})
}
