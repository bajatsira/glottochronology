package config

import (
	"fmt"
	_ "os"
	_ "strconv"
	"time"

	_ "github.com/joho/godotenv"
	//_ log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type JwtConfig struct {
	Token         string        `mapstructure:"token"`
	ExpiresIn     time.Duration `mapstructure:"expires_in"`
	SigningMethod string        `mapstructure:"signing_method"`
}

type Config struct {
	ServiceHost string    `mapstructure:"ServiceHost"`
	ServicePort int       `mapstructure:"ServicePort"`
	JWT         JwtConfig `mapstructure:"jwt"`
	Redis       RedisConfig
}

type RedisConfig struct {
	Host        string
	Password    string
	Port        int
	User        string
	DialTimeout time.Duration
	ReadTimeout time.Duration
}

/*
func NewConfig() (*Config, error) {
	var err error

	configName := "config"
	_ = godotenv.Load()
	if os.Getenv("CONFIG_NAME") != "" {
		configName = os.Getenv("CONFIG_NAME")
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	viper.WatchConfig()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	fmt.Println("Ищу конфиг в:", viper.ConfigFileUsed())

	cfg := &Config{}           // создаем объект конфига
	err = viper.Unmarshal(cfg) // читаем информацию из файла,
	// конвертируем и затем кладем в нашу переменную cfg
	if err != nil {
		return nil, err
	}

	// парсим JWT настройки
	expiresInStr := viper.GetString("JWT.ExpiresIn")
	expiresIn, err := time.ParseDuration(expiresInStr)
	if err != nil {
		return nil, err
	}
	cfg.JWT.ExpiresIn = expiresIn

	// парсим Redis настройки из .env
	cfg.Redis.Host = os.Getenv("REDIS_HOST")
	cfg.Redis.Port, err = strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		return nil, fmt.Errorf("redis port must be int value: %w", err)
	}
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")
	cfg.Redis.User = os.Getenv("REDIS_USER")

	// парсим Redis таймауты из config.toml
	dialTimeoutStr := viper.GetString("Redis.DialTimeout")
	dialTimeout, err := time.ParseDuration(dialTimeoutStr)
	if err != nil {
		return nil, err
	}
	cfg.Redis.DialTimeout = dialTimeout

	readTimeoutStr := viper.GetString("Redis.ReadTimeout")
	readTimeout, err := time.ParseDuration(readTimeoutStr)
	if err != nil {
		return nil, err
	}
	cfg.Redis.ReadTimeout = readTimeout

	log.Info("config parsed")

	return cfg, nil
}*/

func NewConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	fmt.Println("Ищу конфиг в:", viper.ConfigFileUsed())

	var conf Config
	// Именно этот метод переносит данные из viper в вашу структуру
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, err
	}

	// Если после Unmarshal поле пустое, выдаем ошибку вручную для проверки
	if conf.JWT.ExpiresIn == 0 {
		return nil, fmt.Errorf("expires_in is not set or not parsed correctly")
	}

	return &conf, nil
}
