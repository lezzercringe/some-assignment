package main

import (
	"context"
	"flag"
	"log/slog"
	"order-persistor/internal/api"
	"order-persistor/internal/config"
	"order-persistor/internal/inmemory"
	"order-persistor/internal/kafka"
	"order-persistor/internal/log"
	"order-persistor/internal/postgres"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	if err := config.LoadAndValidate(cfgFile, &cfg); err != nil {
		slog.Error("failure loading config", "err", err)
		os.Exit(1)
	}

	logger, err := log.NewLogger(cfg.Log)
	if err != nil {
		slog.Error("could not create logger", "err", err)
		os.Exit(1)
	}

	pool, err := pgxpool.New(context.Background(), cfg.Postgres.ConnString)
	if err != nil {
		logger.Error("creating pg pool", "err", err)
		return
	}
	defer pool.Close()

	paymentsDAO := postgres.PaymentsDAO{
		Pool: pool,
	}

	itemsDAO := postgres.ItemsDAO{
		Pool: pool,
	}

	ordersRepository := postgres.OrdersRepository{
		ItemsDAO:    &itemsDAO,
		PaymentsDAO: &paymentsDAO,
		Pool:        pool,
	}

	cachingOrdersRepository, err := inmemory.NewOrdersCache(cfg.Cache, &ordersRepository, logger)
	if err != nil {
		logger.Error("creating orders cache", "err", err)
		return
	}

	ordersConsumer, err := kafka.NewOrdersConsumer(cfg.KafkaConsumer, cachingOrdersRepository, logger)
	if err != nil {
		logger.Error("failure creating order consumer", "err", err)
		return
	}

	srv := api.NewServer(cfg.API, api.Params{
		Logger:           logger,
		OrdersRepository: cachingOrdersRepository,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	setupSignalHandler(func() {
		logger.Info("received shutdown signal")
		cancel()
	})

	if cfg.Cache.Prefill.Enabled {
		ctx, cancel := context.WithTimeout(ctx, cfg.Cache.Prefill.Timeout)
		defer cancel()

		if err := cachingOrdersRepository.Prefill(ctx); err != nil {
			logger.Error("error pre-filling orders cache", "err", err)
			os.Exit(1)
		}
	}

	go func() {
		err := srv.ListenAndServe()
		logger.Info("api server stopped", "err", err)
		cancel()
	}()

	go func() {
		err := ordersConsumer.Run()
		logger.Info("kafka consumer stopped", "err", err)
		cancel()
	}()

	<-ctx.Done()
	logger.Info("shutting down...")

	ordersConsumer.Stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}

func setupSignalHandler(signalHandler func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-c
		signalHandler()
	}()
}
