package gaarx

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
)

type (
	config struct {
		configData interface{}
	}

	ConfigSource string
)

func (c *config) initConfig(config interface{}) {
	c.configData = config
}

func (c *config) loadConfig(path string) error {
	viper.SetConfigName("config")
	viper.AddConfigPath(path)
	if err := viper.ReadInConfig(); err != nil {
		return errors.New(fmt.Sprintf("error reading config file, %s", err))
	}
	err := viper.Unmarshal(&c.configData)
	if err != nil {
		return errors.New(fmt.Sprintf("unable to decode into struct, %v", err))
	}
	return nil
}

func (c *config) Config() interface{} {
	return c.configData
}
