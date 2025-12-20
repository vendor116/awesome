package httpserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func Run(
	ctx context.Context,
	g *errgroup.Group,
	address string,
	h http.Handler,
	logger *slog.Logger,
	opts ...OptionFunc,
) {
	o := getDefaultOptions()

	for _, opt := range opts {
		if err := opt(&o); err != nil {
			logger.ErrorContext(ctx, "failed to configure option", "error", err)
		}
	}

	server := &http.Server{
		Addr:              address,
		Handler:           h,
		ReadHeaderTimeout: o.readHeaderTimeout,
		WriteTimeout:      o.writeTimeout,
		IdleTimeout:       o.idleTimeout,
		ReadTimeout:       o.readTimeout,
		MaxHeaderBytes:    o.maxHeaderBytes,
	}

	g.Go(func() error {
		logger.InfoContext(ctx, "starting server", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		logger.InfoContext(ctx, "server shutdown gracefully")
		return nil
	})

	g.Go(func() error {
		<-ctx.Done()
		if err := context.Cause(ctx); err != nil && !errors.Is(err, context.Canceled) {
			return err
		}

		shutdownCtx, cancel := context.WithTimeout(context.Background(), o.shutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.InfoContext(ctx, "failed to shutdown server", "error", err)
		}
		return nil
	})
}
