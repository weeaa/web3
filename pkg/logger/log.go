package logger

import (
	"github.com/charmbracelet/log"
)

func LogStartup(module string) {
	log.Info(module + " Monitor Started")
}

func LogError(module string, err error) {
	log.Error(module, "error", err)
}
