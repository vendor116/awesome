package internal_test

import (
	"log/slog"
	"testing"

	"github.com/vendor116/awesome/internal"
)

func TestSetJSONLogger(t *testing.T) {
	if err := internal.SetJSONLogger("debug", "dev"); err != nil {
		t.Fatal(err)
	}

	h := slog.Default().Handler()
	if _, ok := h.(*slog.JSONHandler); !ok {
		t.Error("expected slog.JSONHandler")
	}
}
