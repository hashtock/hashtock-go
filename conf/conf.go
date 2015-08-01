package conf

import (
	"time"
)

type Config struct {
	General GeneralConf

	Jobs struct {
		BankOrders time.Duration
		TagValues  time.Duration
	}
}

type GeneralConf struct {
	AppAddress  string
	AuthAddress string
	ServeAddr   string
	DB          string
	DBName      string
}

var cfg *Config = nil

func GetConfig() *Config {
	if cfg == nil {
		loadConfig()
	}

	return cfg
}
