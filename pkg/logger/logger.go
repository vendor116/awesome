package logger

import (
	"log/slog"
	"os"
)

func SetupJSONLogger(level, version string) error {
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

func newServerLogger(name string) *slog.Logger {
	return slog.Default().With(slog.String("server", name))
}

type Factory func() *slog.Logger

func GrpcLoggerFactory() Factory {
	return func() *slog.Logger {
		return newServerLogger("grpc")
	}
}

func PprofLoggerFactory() Factory {
	return func() *slog.Logger {
		return newServerLogger("pprof")
	}
}

func RestLoggerFactory() Factory {
	return func() *slog.Logger {
		return newServerLogger("rest")
	}
}
