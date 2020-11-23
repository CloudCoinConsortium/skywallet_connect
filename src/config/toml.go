package config

import (
	"github.com/BurntSushi/toml"
	"error"
)

type MainSection struct {
	Timeout int `toml:"timeout"`
	Domain string `toml:"domain"`
	MaxFixTransferNotes int `toml:"max_fixtransfer_notes"`
}

type RConfig struct {
	Title string `toml:"title"`
	Main MainSection `toml:"main"`
	Help string `toml:"help"`
}

func Apply(data string) *error.Error {
	var conf RConfig
	if _, err := toml.Decode(data, &conf); err != nil {
		return &error.Error{ERROR_CONFIG_PARSE, err.Error()}
	}

	if conf.Main.Timeout != 0 {
		DEFAULT_TIMEOUT = conf.Main.Timeout
	}

	if conf.Main.Domain != "" {
		DEFAULT_DOMAIN = conf.Main.Domain
	}

	if conf.Main.MaxFixTransferNotes != 0 {
		MAX_FIXTRANSFER_NOTES = conf.Main.MaxFixTransferNotes
	}

	return nil
}
