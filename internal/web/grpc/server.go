package grpc

import (
	"context"
	"log/slog"
	"net"
	"strconv"

	"github.com/vendor116/awesome/internal/config"
	"github.com/vendor116/awesome/internal/web/grpc/awesome"
	awesomepb "github.com/vendor116/awesome/pkg/protobuf/awesome"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunServer(ctx context.Context, g *errgroup.Group, cfg config.GRPCServer) {
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
