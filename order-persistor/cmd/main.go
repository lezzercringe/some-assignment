package main

import (
	"context"
	"flag"
	"log/slog"
	"order-persistor/internal/config"
	"order-persistor/internal/inmemory"
	"order-persistor/internal/kafka"
	"order-persistor/internal/log"
	"order-persistor/internal/postgres"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "config.yaml", "path to config file (e.g. config.yaml)")
}

func main() {
	flag.Parse()

	cfgFile, err := os.Open(configPath)
	if err != nil {
		slog.Error("could not open config file", "err", err)
		os.Exit(1)
	}
	defer cfgFile.Close()

	var cfg config.Config
	if err := config.Load(cfgFile, &cfg); err != nil {
		slog.Error("fail loading config", "err", err)
		os.Exit(1)
	}

	logger, err := log.NewLogger(cfg.Log)
	if err != nil {
		slog.Error("could not create logger", "err", err)
		os.Exit(1)
	}

	pool, _ := pgxpool.New(context.Background(), cfg.Postgres.ConnString)

	paymentsRepository := postgres.PaymentsRepository{
		Pool: pool,
	}

	itemDAO := postgres.ItemDAO{
		Pool: pool,
	}

	ordersRepository := postgres.OrdersRepository{
		ItemDAO:            &itemDAO,
		PaymentsRepository: &paymentsRepository,
		Pool:               pool,
	}

	cachingOrdersRepository, _ := inmemory.NewOrdersCache(
		cfg.Cache, &ordersRepository, logger,
	)

	ordersConsumer, err := kafka.NewOrdersConsumer(cfg.KafkaConsumer, cachingOrdersRepository, logger)
	if err != nil {
		logger.Error("fail createing order consumer", "err", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	setupSignalHandler(func() {
		cancel()
	})
	ordersConsumer.Run(ctx)
}

func setupSignalHandler(signalHandler func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-c
		signalHandler()
	}()
}
