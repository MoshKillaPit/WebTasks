package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DB struct {
		Host       string `yaml:"host" validate:"required"`        // Хост базы данных
		Port       int    `yaml:"port" validate:"required"`        // Порт базы данных
		User       string `yaml:"user" validate:"required"`        // Имя пользователя
		Password   string `yaml:"password" validate:"required"`    // Пароль
		DBName     string `yaml:"dbname" validate:"required"`      // Имя базы данных
		SSLMode    string `yaml:"sslmode" validate:"required"`     // Режим SSL
		SearchPath string `yaml:"search_path" validate:"required"` // Схема поиска
	} `yaml:"db"`

	Server struct {
		IP   string `yaml:"ip" validate:"required"`   // IP-адрес сервера
		Port int    `yaml:"port" validate:"required"` // Порт сервера
	} `yaml:"server"`
}

func ViperConfig() (*Config, error) {
	viper.SetConfigName("db")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	var dbConfig Config

	err = viper.Unmarshal(&dbConfig)
	if err != nil {
		log.Fatalf("unable to decode config, %v", err)
	}

	return &dbConfig, nil
}
