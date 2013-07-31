package broadcastd

type Config struct {
	Port uint
}

func NewConfig() *Config {
	return &Config{}
}
