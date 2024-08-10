package config

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"os"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
)

var config *Config

type Config struct {
	service.ServiceConf
	ListenOn string
	Cache    cache.CacheConf
	Mongo    struct {
		URL string
		DB  string
	}
}

func Init() {
	config = new(Config)
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "etc/config.yaml"
	}
	err := conf.Load(path, config)
	if err != nil {
		panic(err)
	}
	err = config.SetUp()
	if err != nil {
		panic(err)
	}
}

func Get() *Config {
	if config == nil {
		fmt.Println("config is nil")
	}
	return config
}
