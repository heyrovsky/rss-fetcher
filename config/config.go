package config

import (
	"github.com/heyrovsky/rsscurator/common/utils"
	"github.com/spf13/viper"
)

var (
	FEEDS []string
)

func IntitilizeConfigs() {
	utils.ImportConfig()

	FEEDS = viper.GetStringSlice("feeds")
}
