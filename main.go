package main

import (
	"raznar.id/static-serve-metadata/config"
	"raznar.id/static-serve-metadata/flags"
	"raznar.id/static-serve-metadata/logger"
)

func main() {
	flagsData := flags.New()

	flagsData.Parse()
	if flagsData.Debug {
		logger.InitDebugValue()
	}

	appConfig, err := config.Load()
	if err != nil {
		logger.System.LogError(err)
		return
	}

	err = RunWeb(appConfig)
	if err != nil {
		logger.System.LogError(err)
		return
	}

}