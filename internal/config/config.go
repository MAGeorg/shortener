package config

import (
	"flag"
	"os"
)

type Config struct {
	Address     string
	BaseAddress string
}

var Conf Config

func Parse() {
	if os.Getenv("SERVER_ADDRESS") != "" || os.Getenv("BASE_URL") != "" {
		Conf.Address = os.Getenv("SERVER_ADDRESS")
		Conf.BaseAddress = os.Getenv("BASE_URL")
		return
	}

	flag.StringVar(&Conf.Address, "a", "localhost:8080", "Address for run server")
	flag.StringVar(&Conf.BaseAddress, "b", "http://localhost:8080", "Base URL for shortener address")
	flag.Parse()
}
