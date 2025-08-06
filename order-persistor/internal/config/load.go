package config

import (
	"io"

	"gopkg.in/yaml.v3"
)

func Load(r io.Reader, out *Config) error {
	if err := yaml.NewDecoder(r).Decode(out); err != nil {
		return err
	}

	return validate(out)
}
