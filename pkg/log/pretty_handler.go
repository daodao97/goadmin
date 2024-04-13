package log

import (
	"context"
	"io"
	"log"
	"log/slog"
	"time"

	"github.com/fatih/color"
)

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String()

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	_log := []any{
		color.New(color.Faint).Sprint(r.Time.Format(time.DateTime)),
		level,
		r.Message,
	}

	r.Attrs(func(a slog.Attr) bool {
		_log = append(_log, slog.Any(color.New(color.FgCyan).Sprintf(a.Key), a.Value.Any()))

		return true
	})

	h.l.Println(_log...)

	return nil
}

func NewPrettyHandler(
	out io.Writer,
	opts PrettyHandlerOptions,
) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}

	return h
}
