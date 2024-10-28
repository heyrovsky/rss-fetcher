package main

import (
	"log"
	"path/filepath"

	"github.com/heyrovsky/rsscurator/common/services"
	"github.com/heyrovsky/rsscurator/common/utils"
	"github.com/heyrovsky/rsscurator/config"
	"github.com/heyrovsky/rsscurator/pkg/writer"
	"go.uber.org/zap"
)

var (
	folderName          string
	lastMonthfolderName string
	fileName            string
)

func init() {
	config.IntitilizeConfigs()
	fileName, folderName, lastMonthfolderName = utils.CreateSubfolderAndNames("feeds")

	folderName = filepath.Join("feeds", folderName)
	lastMonthfolderName = filepath.Join("feeds", lastMonthfolderName)
	fileName = filepath.Join(folderName, fileName)

	if err := utils.CreateFolderIfNotExists(folderName); err != nil {
		log.Println(err)
	}

}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	services.InitServices(logger)
	data, err := services.FetchNewsItems()
	if err != nil {
		log.Fatalln(err)
	}
	if err := writer.JsonWriter(data, fileName, folderName, lastMonthfolderName); err != nil {
		log.Println(err)
	}

}
