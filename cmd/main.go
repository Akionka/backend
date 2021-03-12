package main

import (
	"github.com/kate-network/backend/cache"
	"github.com/kate-network/backend/internal"
	"github.com/kate-network/backend/storage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

const envPrefix = "KATE"

type Config struct {
	Redis RedisConfig `mapstructure:"redis"`
	HTTP  HTTPConfig  `mapstructure:"http"`
	DB    string      `mapstructure:"db"`
}

type RedisConfig struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
}

type HTTPConfig struct {
	Address string `mapstructure:"address"`
}

func parseCfg() *Config {
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("redis.address", "127.0.0.1:6379")
	viper.SetDefault("redis.password", "")

	viper.SetDefault("http.address", "127.0.0.1:7454")
	viper.SetDefault("db", "root:@tcp(127.0.0.1:3306)/kate?parseTime=true")

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		logrus.Fatal(err)
	}
	return &c
}

func main() {
	cfg := parseCfg()
	stor, err := storage.Open(cfg.DB)
	if err != nil {
		panic(err)
	}
	ch, err := cache.New(cfg.Redis.Address)
	if err != nil {
		panic(err)
	}

	service := internal.NewServer(stor, ch)
	service.Init()

	if err := service.Listen(cfg.HTTP.Address); err != nil {
		panic(err)
	}
}
