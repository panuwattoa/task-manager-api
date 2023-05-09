package config

import (
	"fmt"
	"io/ioutil"

	"github.com/caarlos0/env/v6"
	"gopkg.in/yaml.v3"
)

const (
	configFileName = "config.yaml"
)

var (
	Env = environmentConfig{}
)

// Env data structure
type environmentConfig struct {
	AppEnv string `env:"APP_ENV,required"`
}

// LoadEnv gets configuration from environment then parse to Env variable
func LoadEnv() error {
	function := "LoadEnv"

	err := env.Parse(&Env)
	if err != nil {
		return fmt.Errorf("%s: unable to read environment configuration, %s", function, err.Error())
	}

	return nil
}

// LoadFile gets configuration from file then parse to App variable
func LoadFile() error {
	function := "LoadFile"

	// read configuration file
	data, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return fmt.Errorf("%s: unable to read file configuration, %s", function, err.Error())
	}

	// unmarshal configuration data to map[string]interface{}
	var dataMap map[string]interface{}
	err = yaml.Unmarshal(data, &dataMap)
	if err != nil {
		return fmt.Errorf("%s: unable to unmarshal configuration data to map, %s", function, err.Error())
	}

	// get configuration by environment(local, dev, or etc.) and marshal to binary
	envData, err := yaml.Marshal(dataMap[Env.AppEnv])
	if err != nil {
		return fmt.Errorf("%s: unable to marshal configuration data, %s", function, err.Error())
	}

	// unmarshal configuration data to App variable
	err = yaml.Unmarshal(envData, &Conf)
	if err != nil {
		return fmt.Errorf("%s: unable to unmarshal configuration data to App, %s", function, err.Error())
	}

	return nil
}
