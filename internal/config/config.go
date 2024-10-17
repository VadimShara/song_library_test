package config

import (
	"sync"
	"song-lib/pkg/logging"
	
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	Listen  struct {
		Type   string `yaml:"type"`
		BindIP string `yaml:"bind_ip"`
		Port   string `yaml:"port"`
	} `yaml:"listen"`
	//Storage StorageConfig `yaml:"storage"`
}

// type StorageConfig struct {
// 	Host 		string	`json:"host"`
// 	Port		string	`json:"port"`
// 	Database	string	`json:"database"`
// 	Username	string	`json:"username"`
// 	Password	string	`json:"password"`
// }

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application configuraion")
		instance = &Config{}
		if err := cleanenv.ReadConfig("C:/Users/vadim/song-lib/config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}