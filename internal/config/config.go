package config

import "os"

type Config struct {
	Address string
}

func Load() Config {
	address := os.Getenv("APP_ADDRESS")
	if address == "" {
		address = ":8080"
	}

	return Config{
		Address: address,
	}
}
