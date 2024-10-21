package config

import (
	"github.com/heyrovsky/rsscurator/common/utils"
	"github.com/spf13/viper"
)

var (
	CYBERSEC   []string
	TECHNOLOGY []string
	BLOCKCHAIN []string
)

func IntitilizeConfigs() {
	utils.ImportConfig()

	CYBERSEC = viper.GetStringSlice("cybersec")
	TECHNOLOGY = viper.GetStringSlice("technology")
	BLOCKCHAIN = viper.GetStringSlice("blockchain")
}
