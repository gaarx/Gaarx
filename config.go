package gaarx

import (
	"errors"
	confLib "github.com/micro/go-config"
	"github.com/micro/go-config/encoder"
	"github.com/micro/go-config/encoder/json"
	"github.com/micro/go-config/encoder/toml"
	"github.com/micro/go-config/encoder/xml"
	"github.com/micro/go-config/encoder/yaml"
	"github.com/micro/go-config/source"
	"github.com/micro/go-config/source/file"
)

type (
	config struct {
		c          confLib.Config
		configData interface{}
	}

	ConfigSource string
)

const (
	ConfigSourceToml = "toml"
	ConfigSourceYaml = "yaml"
	ConfigSourceJson = "json"
	ConfigSourceXML  = "xml"
)

func (c *config) initConfig(config interface{}) {
	c.c = confLib.NewConfig()
	c.configData = config
}

func (c *config) loadConfig(way ConfigSource, path string) error {
	var e encoder.Encoder
	switch way {
	case ConfigSourceToml:
		e = toml.NewEncoder()
		break
	case ConfigSourceYaml:
		e = yaml.NewEncoder()
		break
	case ConfigSourceJson:
		e = json.NewEncoder()
		break
	case ConfigSourceXML:
		e = xml.NewEncoder()
		break
	default:
		return errors.New("unsupported config source")
	}
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
