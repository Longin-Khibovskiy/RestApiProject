package sl

import (
	"log/slog"

	_ "github.com/Longin-Khibovskiy/RestApiProject.git/internal/lib/logger/handlers/slogdiscard"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
