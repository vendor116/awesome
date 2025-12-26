package logger_test

import (
	"log/slog"
	"testing"

	"github.com/vendor116/awesome/pkg/logger"
)

func TestSetJSONLogger(t *testing.T) {
	if err := logger.SetupJSONLogger("debug", "dev"); err != nil {
		t.Fatal(err)
	}

	h := slog.Default().Handler()
	if _, ok := h.(*slog.JSONHandler); !ok {
		t.Error("expected slog.JSONHandler")
	}
}
