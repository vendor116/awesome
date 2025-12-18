package internal

import (
	"log/slog"
	"os"
)

func SetJSONLogger(level string) error {
	var l slog.Level
	if err := l.UnmarshalText([]byte(level)); err != nil {
		return err
	}

	slog.SetDefault(
		slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: l,
			}),
		),
	)
	return nil
}
