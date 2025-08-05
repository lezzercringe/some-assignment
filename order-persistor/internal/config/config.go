package config

type Cache struct {
	Size int `yaml:"size"`
}

type Config struct {
	Cache Cache `yaml:"cache"`
}
