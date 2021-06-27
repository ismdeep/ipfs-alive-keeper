package config

import "github.com/BurntSushi/toml"

type Config struct {
	Links    []string
	Routines uint
}

var DefaultConf *Config

// Load 加载配置
func Load(configPath string) error {
	DefaultConf = &Config{}

	if _, err := toml.DecodeFile(configPath, DefaultConf); err != nil {
		return err
	}

	return nil
}
