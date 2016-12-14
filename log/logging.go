package log

import (
	"github.com/op/go-logging"
	"os"
)

var Logger *logging.Logger

func init() {
	Logger = logging.MustGetLogger("mango")
	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	formattedBackend := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(formattedBackend)
	Logger.Debug("Logging started.")
}

func Critical(args ...interface{}) {
	Logger.Critical(args)
}

func Debug(args ...interface{}) {
	Logger.Debug(args)
}

func Error(args ...interface{}) {
	Logger.Error(args)
}

func Notice(args ...interface{}) {
	Logger.Notice(args)
}

func Warning(args ...interface{}) {
	Logger.Warning(args)
}

func Info(args ...interface{}) {
	Logger.Info(args)
}
