package config

import (
	"github.com/spf13/viper"
)

type Elastic struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Config struct {
	Elastic Elastic `yaml:"elastic"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var cfg Config

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
