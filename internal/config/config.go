package config

import (
	"flag"
)

type Config struct {
	Address     string
	BaseAddress string
}

var Conf Config

func Parse() {
	flag.StringVar(&Conf.Address, "a", "localhost:8080", "Address for run server")
	flag.StringVar(&Conf.BaseAddress, "b", "http://localhost:8080", "Base URL for shortener address")
	flag.Parse()
}
