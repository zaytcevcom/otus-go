package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/zaytcevcom/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/zaytcevcom/otus-go/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/zaytcevcom/otus-go/hw12_13_14_15_calendar/internal/sender"
	"os/signal"
	"syscall"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/calendar_sender/config.toml", "Path to configuration file")
}

//nolint:unused
func main() {
	flag.Parse()

	config, err := LoadSenderConfig(configFile)
	if err != nil {
		fmt.Println("Error loading config: ", err)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logg := logger.New(config.Logger.Level, nil)

	broker, err := rabbitmq.NewRabbitMQ(logg, config.Rabbit.Uri, config.Rabbit.Exchange, config.Rabbit.Queue)
	if err != nil {
		fmt.Println("cannot connect to rabbit", err)
		return
	}

	s := sender.New(logg, broker)

	go func() {
		<-ctx.Done()

		s.Stop()
	}()

	if err := s.Start(); err != nil {
		logg.Error("failed to start grpc server: " + err.Error())
	}
}
