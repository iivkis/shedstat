package config

import (
	"sync"

	"github.com/spf13/viper"
)

type config struct {
	DB struct {
		Name     string
		User     string
		Password string
		Host     string
		Port     string
	}
	ClickHouse struct {
		Name     string
		User     string
		Password string
		Host     string
		Port     string
	}
}

var cfg = &config{}
var once = &sync.Once{}

func Get() *config {
	once.Do(func() {
		viper.BindEnv("DB.Name", "DB_NAME")
		viper.BindEnv("DB.User", "DB_USER")
		viper.BindEnv("DB.Password", "DB_PASSWORD")
		viper.BindEnv("DB.Host", "DB_HOST")
		viper.BindEnv("DB.Port", "DB_PORT")

		viper.BindEnv("ClickHouse.Name", "CLICKHOUSE_NAME")
		viper.BindEnv("ClickHouse.User", "CLICKHOUSE_USER")
		viper.BindEnv("ClickHouse.Password", "CLICKHOUSE_PASSWORD")
		viper.BindEnv("ClickHouse.Host", "CLICKHOUSE_HOST")
		viper.BindEnv("ClickHouse.Port", "CLICKHOUSE_PORT")

		if err := viper.Unmarshal(&cfg); err != nil {
			panic(err)
		}
	})
	return cfg
}
