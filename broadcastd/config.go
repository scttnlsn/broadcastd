package broadcastd

type Config struct {
	Auth string
	Port uint
}

func NewConfig() *Config {
	return &Config{}
}
