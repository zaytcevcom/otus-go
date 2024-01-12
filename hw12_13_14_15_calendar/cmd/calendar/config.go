package main

import "github.com/spf13/viper"

type Config struct {
	Logger          LoggerConf
	ServerHTTP      ServerHTTPConf
	ServerGRPC      ServerGRPCConf
	IsMemoryStorage bool
	Postgres        PostgresConf
}

type LoggerConf struct {
	Level string
}

type ServerHTTPConf struct {
	Host string
	Port int
}

type ServerGRPCConf struct {
	Host string
	Port int
}

type PostgresConf struct {
	Dsn string
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
