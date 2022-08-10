package console

import (
	"os"

	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/luckyAkbar/himatro-telegram-bot/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use: "go run main.go",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	setupLogger()
}

func setupLogger() {
	formatter := runtime.Formatter{
		ChildFormatter: &log.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		},
		Line: true,
		File: true,
	}

	log.SetFormatter(&formatter)
	log.SetOutput(os.Stdout)

	logLevel, err := log.ParseLevel(config.LogLevel())

	if err != nil {
		logLevel = log.DebugLevel
	}

	log.SetLevel(logLevel)
}
