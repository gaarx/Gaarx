package gaarx

import (
	configlib "github.com/micro/go-config"
	"github.com/micro/go-config/encoder/toml"
	"github.com/micro/go-config/source"
	"github.com/micro/go-config/source/file"
)

type (
	config struct {
		c          configlib.Config
		configData ConfigAble
	}

	ConfigSource string
)

const (
	ConfigSourceToml = "toml"
)

func (c *config) initConfig(config ConfigAble) {
	c.c = configlib.NewConfig()
	c.configData = config
}

func (c *config) loadConfig(way ConfigSource, path string) error {
	switch way {
	case ConfigSourceToml:
		e := toml.NewEncoder()
		fileSource := file.NewSource(
			file.WithPath(path),
			source.WithEncoder(e),
		)
		err := c.c.Load(fileSource)
		_ = c.c.Scan(&c.configData)
		return err
	}
	return nil
}

func (c *config) Config() ConfigAble {
	return c.configData
}
