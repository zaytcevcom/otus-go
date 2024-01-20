package main

import "github.com/spf13/viper"

type Config struct {
	Logger    LoggerConf
	Postgres  PostgresConf
	Rabbit    RabbitConf
	Scheduler SchedulerConf
}

type LoggerConf struct {
	Level string
}

type PostgresConf struct {
	Dsn string
}

type RabbitConf struct {
	Uri      string
	Exchange string
	Queue    string
}

type SchedulerConf struct {
	Interval int
}

func LoadSchedulerConfig(path string) (Config, error) {
	config := Config{}

	viper.SetConfigFile(path)

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	return config, err
}
