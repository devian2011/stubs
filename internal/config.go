package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	stubHttp "stubs/internal/stub/http"
)

type LogProvider string

type Configuration struct {
	Http struct {
		Addr string `json:"addr" yaml:"addr"`
	} `json:"http" yaml:"http"`
	Stubs struct {
		Http []stubHttp.Config `json:"http" yaml:"http"`
	} `json:"stubs" yaml:"stubs"`
}

func NewConfiguration(configFile string) (*Configuration, error) {
	if fStat, fErr := os.Stat(configFile); fErr != nil || fStat.IsDir() {
		return nil, errors.New(fmt.Sprintf("file does not exists: %s", configFile))
	}
	cfg := &Configuration{}
	data, readErr := os.ReadFile(configFile)
	if readErr != nil {
		return nil, readErr
	}
	ext := filepath.Ext(configFile)
	switch ext {
	case ".json":
		unmarshalErr := json.Unmarshal(data, cfg)
		if unmarshalErr != nil {
			return nil, unmarshalErr
		}
		break
	case ".yaml":
		unmarshalErr := yaml.Unmarshal(data, cfg)
		if unmarshalErr != nil {
			return nil, unmarshalErr
		}
		break
	default:
		return nil, errors.New(fmt.Sprintf("unsupported config file extension: %s. filePath: %s", ext, configFile))
	}

	return cfg, nil
}
