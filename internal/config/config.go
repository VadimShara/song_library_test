package config

import (
	"sync"
	"song-lib/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug bool `env:"LISTEN_IS_DEBUG"`
	Listen  struct {
		Type   string `env:"LISTEN_TYPE"`
		BindIP string `env:"LISTEN_BIND_IP"`
		Port   string `env:"LISTEN_PORT"`
	}
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application configuraion")
		instance = &Config{}
		if err := cleanenv.ReadConfig("C:/Users/vadim/song-lib/.env", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}