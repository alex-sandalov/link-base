package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		HTTP     HTTPConfig
		Postgres PostgresConfig
		Redis    RedisConfig
		JWT      JWTConfig
		SMPT     SMPTConfig
	}

	HTTPConfig struct {
		Port           string        `yaml:"port"`
		MaxHeaderBytes int           `yaml:"maxHeaderBytes"`
		ReadTimeout    time.Duration `yaml:"readTimeout"`
		WriteTimeout   time.Duration `yaml:"writeTimeout"`
	}

	PostgresConfig struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `env:"POSTGRES_USER" env-default:"postgres"`
		Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
		Database string `yaml:"database"`
		SSLMode  string `yaml:"sslMode"`
	}

	RedisConfig struct {
		Address        string        `yaml:"addr"`
		DatabaseNumber int           `yaml:"numberDB"`
		DialTimeout    time.Duration `yaml:"dialTimeout"`
		ReadTimeout    time.Duration `yaml:"readTimeout"`
		WriteTimeout   time.Duration `yaml:"writeTimeout"`
		PoolSize       int           `yaml:"poolSize"`
		MinIdleConns   int           `yaml:"minIdleConns"`
	}

	JWTConfig struct {
		AccessTokenTTL  time.Duration `yaml:"accessTokenTTL"`
		RefreshTokenTTL time.Duration `yaml:"refreshTokenTTL"`
		SigningKey      string        `env:"SIGNING_KEY"`
	}

	SMPTConfig struct {
		SMPTHost     string `yaml:"smptHost"`
		SMPTPort     string `yaml:"smptPort"`
		SMPTUser     string `yaml:"smptUser"`
		SMPTPassword string `yaml:"smptPassword"`
	}
)

// MustLoad loads the configuration from the file specified in the CONFIG_PATH environment variable.
// It terminates the program with a fatal log if CONFIG_PATH is not set, the file cannot be found,
// or if there is an error reading the configuration.
func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalf("CONFIG_PATH IS NOT SET")
	}

	envPath := os.Getenv("ENV_PATH")
	if envPath == "" {
		log.Fatalf("ENV_PATH IS NOT SET")
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("Configuration file not found: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(envPath, &cfg); err != nil {
		log.Fatalf("Error reading environment variables: %s", err)
	}

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Error reading configuration: %s", err)
	}

	return &cfg
}
