package main

import "github.com/spf13/viper"

type Config struct {
	Logger          LoggerConf
	Server          ServerConf
	IsMemoryStorage bool
}

type LoggerConf struct {
	Level string
}

type ServerConf struct {
	Host string
	Port int
}

func LoadConfig(path string) (Config, error) {
	config := Config{}

	viper.SetConfigFile(path)

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	return config, err
}
