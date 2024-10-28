package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

func CreateFolderIfNotExists(folderPath string) error {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.MkdirAll(folderPath, 0755)
		if err != nil {
			return err
		}

		log.Println("FOLDER CREATED", folderPath)
		return nil
	}

	log.Println("Folder Already exists:", folderPath)
	return nil
}

func CreateSubfolderAndNames(basefolder string) (string, string, string) {
	currentTime := time.Now()
	year, month, day := currentTime.Date()
	fileName := fmt.Sprintf("%d-%d-%d.json", day, month, year)
	folderName := fmt.Sprintf("%s-%d", month, year)
	beforeMonth := fmt.Sprintf("%s-%d", month-1, year)

	return fileName, folderName, beforeMonth
}
