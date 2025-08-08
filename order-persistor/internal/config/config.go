package config

import "time"

type Log struct {
	Level string `yaml:"level" validate:"required"`
}

type Cache struct {
	Size int `yaml:"size" validate:"required,gte=0"`
}

type Prefill struct {
	Enabled bool          `yaml:"enabled"`
	Timeout time.Duration `yaml:"timeout"`
}

type KafkaConsumer struct {
	Servers            string        `yaml:"servers" validate:"required"`
	GroupID            string        `yaml:"group_id"`
	Topic              string        `yaml:"topic" validate:"required"`
	ReadTimeout        time.Duration `yaml:"read_timeout" validate:"required"`
	ProcessTimeout     time.Duration `yaml:"process_timeout" validate:"required"`
	ReadFailureBackoff time.Duration `yaml:"read_failure_backoff" validate:"required"`
}

type API struct {
	Host    string        `yaml:"host" validate:"required"`
	Port    string        `yaml:"port" validate:"required"`
	Timeout time.Duration `yaml:"timeout" validate:"required"`
}

type Postgres struct {
	ConnString string `yaml:"conn_string" validate:"required"`
}

type Config struct {
	Log           Log           `yaml:"log" validate:"required"`
	Cache         Cache         `yaml:"cache" validate:"required"`
	KafkaConsumer KafkaConsumer `yaml:"kafka_consumer" validate:"required"`
	Postgres      Postgres      `yaml:"postgres" validate:"required"`
	API           API           `yaml:"api" validate:"required"`
	Prefill       Prefill       `yaml:"prefill" validate:"required"`
}
