package config

import (
	"crypto-trading-bot/internal/utils"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

// Config содержит все настройки приложения
type Config struct {
	Postgres struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
	} `mapstructure:"postgres"`

	Web struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"web"`

	Binance struct {
		APIKey    string `mapstructure:"apiKey"`
		APISecret string `mapstructure:"apiSecret"`
	} `mapstructure:"binance"`

	Huobi struct {
		APIKey    string `mapstructure:"apiKey"`
		APISecret string `mapstructure:"apiSecret"`
	} `mapstructure:"huobi"`

	Logging struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"logging"`

	Security struct {
		JWTSecret string `mapstructure:"jwtSecret"`
	} `mapstructure:"security"`

	Backtesting struct {
		DefaultCapital    float64 `mapstructure:"defaultCapital"`
		DefaultCommission float64 `mapstructure:"defaultCommission"`
		DefaultSpread     float64 `mapstructure:"defaultSpread"`
	} `mapstructure:"backtesting"`
}

// LoadConfig загружает конфигурацию из файла и переменных окружения
func LoadConfig() *Config {
	// Получение текущего рабочего каталога
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// обрезаем путь до корня приложения
	currentDir, err = utils.TrimPathToDir(currentDir, "crypto-trading-bot")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Текущий каталог: %v\n", currentDir)

	viper.AddConfigPath(currentDir)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Чтение файла конфигурации
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
		os.Exit(1)
	}

	// Чтение переменных окружения
	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
		os.Exit(1)
	}

	return &config
}
