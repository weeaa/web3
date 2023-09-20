package logger

import (
	"fmt"
	"github.com/charmbracelet/log"
)

func LogStartup(module string) {
	log.Info(fmt.Sprintf("[%s] %s", module, "Monitor Started!"))
}

func LogShutDown(module string) {
	log.Warn(fmt.Sprintf("[%s] %s", module, "Monitor Stopped!"))
}

func LogInfo(module, msg string) {
	log.Info(fmt.Sprintf("[%s] %s", module, msg))
}

func LogError(module string, err error) {
	log.Error(fmt.Sprintf("[%s]", module), "error", err)
}

func LogFatal(module string, msg any) {
	log.Fatal(fmt.Sprintf("[%s] %v", module, msg))
}
