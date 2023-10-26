package logger

import (
	"fmt"
	"github.com/rs/zerolog/log"
	//"github.com/charmbracelet/log"
)

func LogStartup(module string) {
	log.Info().Str(module, "launched")
	//log.Info(fmt.Sprintf("[%s] %s", module, "Monitor Started!"))
}

func LogShutDown(module string) {
	log.Warn().Str(module, "stopped")
}

func LogInfo(module, msg string) {
	log.Info().Str(module, msg)
}

func LogError(module string, err error) {
	log.Error().Str(module, fmt.Sprint(err))
}

func LogFatal(module string, msg string) {
	log.Fatal().Str(module, msg)
}
