package main

import (
	"log"
	"path/filepath"

	"github.com/heyrovsky/rsscurator/common/constants"
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

	folderName = filepath.Join(constants.FEED_FOLDER_NAME, folderName)
	lastMonthfolderName = filepath.Join(constants.FEED_FOLDER_NAME, lastMonthfolderName)
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
		logger.Error(err.Error())
	}
	if err := writer.JsonWriter(data, fileName, folderName, lastMonthfolderName, *logger); err != nil {
		logger.Error(err.Error())
	}

}
