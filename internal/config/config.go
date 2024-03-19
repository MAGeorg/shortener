package config

import (
	"flag"
	"os"
)

type Config struct {
	Address     string
	BaseAddress string
}

func (c *Config) GetBaseAddress() string {
	return c.BaseAddress
}

func NewConfig() *Config {
	return &Config{}
}

func Parse(conf *Config) {
	if os.Getenv("SERVER_ADDRESS") != "" || os.Getenv("BASE_URL") != "" {
		conf.Address = os.Getenv("SERVER_ADDRESS")
		conf.BaseAddress = os.Getenv("BASE_URL")
		return
	}

	flag.StringVar(&conf.Address, "a", "localhost:8080", "Address for run server")
	flag.StringVar(&conf.BaseAddress, "b", "http://localhost:8080", "Base URL for shortener address")
	flag.Parse()
}
