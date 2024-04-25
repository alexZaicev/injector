package logging

import "log/slog"

const (
	Error = "error"
	Queue = "queue"
)

func WithError(err error) slog.Attr {
	return slog.String(Error, err.Error())
}
