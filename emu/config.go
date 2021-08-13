package emu

type Config struct {
	cosmac bool
}

func NewConfig(cosmac bool) *Config {
	return &Config{cosmac}
}
