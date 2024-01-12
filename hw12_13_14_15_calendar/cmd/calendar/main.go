package main

import (
	"context"
	"flag"
	"fmt"
	sqlstorage "github.com/zaytcevcom/otus-go/hw12_13_14_15_calendar/internal/storage/sql"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zaytcevcom/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/zaytcevcom/otus-go/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/zaytcevcom/otus-go/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/zaytcevcom/otus-go/hw12_13_14_15_calendar/internal/storage/memory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.json", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := LoadConfig(configFile)
	if err != nil {
		fmt.Println("Error loading config: ", err)
		return
	}

	fmt.Println(config)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logg := logger.New(config.Logger.Level, nil)

	storage, err := getStorage(ctx, config)
	if err != nil {
		fmt.Println("Failed to get storage: ", err)
		return
	}

	if sqlStorage, ok := storage.(*sqlstorage.Storage); ok {
		defer func(sqlStorage *sqlstorage.Storage, _ context.Context) {
			err := sqlStorage.Close(ctx)
			if err != nil {
				fmt.Println("Cannot close psql connection", err)
			}
		}(sqlStorage, ctx)
	}

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar, config.Server.Host, config.Server.Port)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info(fmt.Sprintf("calendar is running on %s:%d", config.Server.Host, config.Server.Port))

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

func getStorage(ctx context.Context, config Config) (app.Storage, error) {
	if config.IsMemoryStorage {
		storage := memorystorage.New()
		return storage, nil
	}

	storageSql := sqlstorage.New(config.Postgres.Dsn)

	if err := storageSql.Connect(ctx); err != nil {
		return nil, fmt.Errorf("cannot connect to psql: %w", err)
	}

	return storageSql, nil
}
