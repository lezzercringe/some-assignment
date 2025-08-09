package producer

type Config struct {
	Servers string `yaml:"servers"`
	Topic   string `yaml:"topic"`
}
