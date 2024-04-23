package logging

import "log/slog"

const (
	Error = "error"

	Path        = "path"
	Method      = "method"
	Protocol    = "protocol"
	Source      = "source"
	Destination = "destination"
)

func WithError(err error) slog.Attr {
	return slog.String(Error, err.Error())
}
