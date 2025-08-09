package main

import (
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	"producer"

	"gopkg.in/yaml.v3"
)

var (
	configPath   string
	messageCount int
)

func init() {
	flag.StringVar(&configPath, "config", "config.yaml", "config file path")
	flag.IntVar(&messageCount, "count", 1, "count of messages to send (0, 1, 2...)")
}

func main() {
	flag.Parse()

	cfgFile, err := os.Open(configPath)
	if err != nil {
		slog.Error("could not open config file", "err", err)
		os.Exit(1)
	}

	var cfg producer.Config
	if err := yaml.NewDecoder(cfgFile).Decode(&cfg); err != nil {
		slog.Error("could not decode config", "err", err)
		os.Exit(1)
	}

	p, err := producer.NewProducer(cfg)
	if err != nil {
		slog.Error("could not create producer", "err", err)
		os.Exit(1)
	}

	defer p.Close()
	go p.Start()

	for range messageCount {
		order := producer.GenerateOrder()
		msg, _ := json.Marshal(order)

		if err := p.Produce(msg); err != nil {
			slog.Error("failed producing message", "err", err)
			continue
		}

		slog.Info("produced order", "id", order.ID)
	}
}
