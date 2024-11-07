package pkg

import (
	"io"
	"log/slog"
)

func CreateNewLogger(w io.Writer) (*slog.Logger, error) {
	return slog.New(slog.NewJSONHandler(w, nil)), nil
}
