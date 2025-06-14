package logging

import (
	"io"
	"log/slog"
	"net/http"
	"os"
)

const LevelRequest = slog.Level(-2)

func NewRequestLogger(writer io.Writer) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:       LevelRequest,
		ReplaceAttr: replaceAttr,
	})

	return slog.New(handler)
}

func LogRequest(logger *slog.Logger, r *http.Request) {
	logger.Log(
		r.Context(),
		LevelRequest,
		"",
		"method", r.Method,
		"path", r.URL.Path,
	)
}

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	levelNames := map[slog.Leveler]string{
		LevelRequest: "REQUEST",
	}
	if a.Key == slog.LevelKey {
		level := a.Value.Any().(slog.Level)
		if name, exists := levelNames[level]; exists {
			a.Value = slog.StringValue(name)
		}
	}
	return a
}
