package main

import (
	"github.com/heyrovsky/rsscurator/common/services"
	"github.com/heyrovsky/rsscurator/config"
	"go.uber.org/zap"
)

func init() {
	config.IntitilizeConfigs()
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	services.InitServices(logger)
}
