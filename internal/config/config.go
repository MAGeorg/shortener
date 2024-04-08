package config

import (
	"flag"
	"os"
)

const (
	serverAdderss = "SERVER_ADDRESS"
	baseAddress   = "BASE_URL"
	pathToFile    = "FILE_STORAGE_PATH"
)

type Config struct {
	Address         string
	BaseAddress     string
	StorageFileName string
}

func NewConfig() *Config {
	return &Config{}
}

func Parse(conf *Config) {
	flag.StringVar(&conf.Address, "a", "localhost:8080", "Address for run server")
	flag.StringVar(&conf.BaseAddress, "b", "http://localhost:8080", "Base URL for shortener address")
	flag.StringVar(&conf.StorageFileName, "f", "/tmp/short-url-db.json", "Full path to the file where the created shortened URLs are stored")
	flag.Parse()

	if a := os.Getenv(serverAdderss); a != "" {
		conf.Address = a
	}
	if b := os.Getenv(baseAddress); b != "" {
		conf.BaseAddress = b
	}

	if f := os.Getenv(pathToFile); f != "" {
		conf.StorageFileName = f
	}
}
