package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost        string
	DBPort        uint16
	DBUser        string
	DBPassword    string
	DBName        string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()
	cfg := &Config{
		DBHost:        viper.GetString("DB_HOST"),
		DBPort:        viper.GetUint16("DB_PORT"),
		DBUser:        viper.GetString("DB_USER"),
		DBPassword:    viper.GetString("DB_PASSWORD"),
		DBName:        viper.GetString("DB_NAME"),
		RedisAddr:     viper.GetString("REDIS_ADDR"),
		RedisPassword: viper.GetString("REDIS_PASSWORD"),
		RedisDB:       viper.GetInt("REDIS_DB"),
	}
	fmt.Print(cfg)
	return cfg, nil
}
