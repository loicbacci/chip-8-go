package emu

type Config struct {
	// Is the emulated machine COSMAC.
	// Used to interpret some instructions differently.
	cosmac bool
}

func NewConfig(cosmac bool) *Config {
	return &Config{cosmac}
}
