package config

import "github.com/go-playground/validator/v10"

func validate(cfg *Config) error {
	return validator.New().Struct(cfg)
}
