package grpc

import (
	"context"
	"log/slog"
	"net"

	"github.com/vendor116/awesome/internal/config"
	"github.com/vendor116/awesome/pkg/protobuf/awesome"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunServer(
	ctx context.Context,
	eg *errgroup.Group,
	server awesome.AwesomeServer,
	cfg config.GRPCServer,
	loggerFactory func() *slog.Logger,
) {
	s := grpc.NewServer()

	if cfg.ReflectEnabled {
		reflection.Register(s)
	}

	awesome.RegisterAwesomeServer(s, server)

	logger := loggerFactory()

	eg.Go(func() error {
		lc := &net.ListenConfig{}

		listen, err := lc.Listen(ctx, "tcp", cfg.GetAddress())
		if err != nil {
			return err
		}

		logger.Info("starting server", "address", listen.Addr().String())
		err = s.Serve(listen)
		if err != nil {
			return err
		}

		logger.Info("server shutdown gracefully")
		return nil
	})

	eg.Go(func() error {
		<-ctx.Done()
		s.GracefulStop()
		return nil
	})
}
