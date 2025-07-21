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

	appConfig, err := config.New(flagsData.ConfigPath, true)
	if err != nil {
		logger.System.LogError(err)
		return
	}

	logger.System.LogInfo(appConfig.IsExists())
	logger.System.LogInfo(flagsData.Debug)

	err = RunWeb(appConfig)
	if err != nil {
		logger.System.LogError(err)
		return
	}

}