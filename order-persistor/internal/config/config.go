package config

import "time"

type Log struct {
	Level string `yaml:"level"`
}

type Cache struct {
	Size int `yaml:"size"`
}

type KafkaConsumer struct {
	Servers     []string      `yaml:"servers"`
	GroupID     string        `yaml:"group_id"`
	Topic       string        `yaml:"topic"`
	ReadTimeout time.Duration `yaml:"read_timeout"`
}

type Postgres struct {
	ConnString string `yaml:"conn_string"`
}

type Config struct {
	Log           Log           `yaml:"log"`
	Cache         Cache         `yaml:"cache"`
	KafkaConsumer KafkaConsumer `yaml:"kafka_consumer"`
	Postgres      Postgres      `yaml:"postgres"`
}
