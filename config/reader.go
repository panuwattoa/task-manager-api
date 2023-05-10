package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Read reads configuration data and unmarshall to models.Config
func Read(config interface{}) error {
	currentDir, _ := os.Getwd()
	rViper := viper.New()
	rViper.AddConfigPath(currentDir)             //default older structure
	rViper.AddConfigPath(currentDir + "/../../") //default new structure
	rViper.SetConfigName("config")
	rViper.SetConfigType("yml")
	if err := rViper.ReadInConfig(); err != nil {
		fmt.Printf("cannot read configuration config.yml : %v\n", err)
		return err
	}

	//override
	rViper.SetConfigFile(currentDir + "/config/config.yml")
	err := rViper.MergeInConfig()
	if err != nil {
		rViper.SetConfigFile(currentDir + "/../../config/config.yml")
		rViper.MergeInConfig()
	}

	// prevention hidden issue when unmarshal fail
	if err := rViper.Unmarshal(config); err != nil {
		fmt.Printf("cannot unmarshal configuration, %v\n", err)
		return err
	}

	return nil
}
