package logger

import (
	"fmt"
	"github.com/charmbracelet/log"
)

func LogStartup(module string) {

	log.Info(module + " Monitor Started")
}

func LogInfo(module, msg string) {
	log.Info(fmt.Sprintf("%s %s", module, msg))
}

func LogError(module string, err error) {
	log.Error(module, "error", err)
}
