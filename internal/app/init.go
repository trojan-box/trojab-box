package app

import (
	"fmt"
	"github.com/aresprotocols/trojan-box/internal/pkg/config"
)

func Init() *Config {
	cfg, err := loadConf()
	if err != nil {
		panic(fmt.Sprintf("load app conf err: %v", err))
	}
	Conf = cfg
	return cfg
}

// loadConf load app config
func loadConf() (ret *Config, err error) {
	var cfg Config
	if err := config.Load("app", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
