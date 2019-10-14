package gaarx

import (
	confLib "github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/encoder"
	"github.com/micro/go-micro/config/source"
	"github.com/micro/go-micro/config/source/file"
)

type (
	config struct {
		c          confLib.Config
		configData interface{}
	}

	ConfigSource string
)

func (c *config) initConfig(config interface{}) {
	c.c = confLib.NewConfig()
	c.configData = config
}

func (c *config) loadConfig(path string) error {
	var e encoder.Encoder
	fileSource := file.NewSource(
		file.WithPath(path),
		source.WithEncoder(e),
	)
	err := c.c.Load(fileSource)
	if err != nil {
		return err
	}
	return c.c.Scan(&c.configData)
}

func (c *config) Config() interface{} {
	return c.configData
}
