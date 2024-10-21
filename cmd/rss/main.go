package main

import (
	"fmt"
	"log"

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
	data, err := services.FetchNewsItems()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(len(data))
}
