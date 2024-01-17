package main

import (
	"context"
	"flag"
	"fmt"
	internalgrpc "github.com/zaytcevcom/otus-go/hw12_13_14_15_calendar/internal/server/grpc"
	sqlstorage "github.com/zaytcevcom/otus-go/hw12_13_14_15_calendar/internal/storage/sql"
	"os"
	"os/signal"
	"sync"
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

	defer func(storage app.Storage, ctx context.Context) {
		err := storage.Close(ctx)
		if err != nil {
			fmt.Println("Cannot close storage connection", err)
		}
	}(storage, ctx)

	calendar := app.New(logg, storage)

	serverHTTP := internalhttp.NewServer(logg, calendar, config.ServerHTTP.Host, config.ServerHTTP.Port)
	serverGRPC := internalgrpc.NewServer(logg, calendar, config.ServerGRPC.Host, config.ServerGRPC.Port)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := serverHTTP.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}

		serverGRPC.Stop()
	}()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		logg.Info(fmt.Sprintf("calendar HTTP is running on %s:%d", config.ServerHTTP.Host, config.ServerHTTP.Port))

		if err := serverHTTP.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
			os.Exit(1) //nolint:gocritic
		}
	}()

	go func() {
		defer wg.Done()
		logg.Info(fmt.Sprintf("calendar GRPC is running on %s:%d", config.ServerGRPC.Host, config.ServerGRPC.Port))

		if err := serverGRPC.Start(ctx); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			cancel()
			os.Exit(1) //nolint:gocritic
		}
	}()

	wg.Wait()
}

func getStorage(ctx context.Context, config Config) (app.Storage, error) {
	var storage app.Storage

	if config.IsMemoryStorage {
		storage = memorystorage.New()
	} else {
		storage = sqlstorage.New(config.Postgres.Dsn)
	}

	if err := storage.Connect(ctx); err != nil {
		return nil, fmt.Errorf("cannot connect to storage: %w", err)
	}

	return storage, nil
}
