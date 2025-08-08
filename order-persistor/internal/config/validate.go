package config

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

func validate(cfg *Config) error {
	if err := validator.New().Struct(cfg); err != nil {
		return err
	}

	if err := validatePrefill(&cfg.Prefill); err != nil {
		return err
	}

	return nil
}

func validatePrefill(p *Prefill) error {
	if p.Enabled && p.Timeout <= 0 {
		return errors.New("prefill timeout should be >= 0 if prefill is enabled")
	}

	return nil
}
