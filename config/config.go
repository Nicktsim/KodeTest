package config

import (
    "log"
    "os"
    "time"

    "github.com/ilyakaznacheev/cleanenv"
    "github.com/joho/godotenv"
)

type Config struct {
    Env         string `yaml:"env" env-default:"local"`
    HTTPServer  `yaml:"http_server"`
    Database    `yaml:"database"`
    PgAdmin     `yaml:"pgadmin"`
}

type HTTPServer struct {
    Address     string        `yaml:"address" env-default:"localhost:8080"`
    Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
    IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
    // User        string        `yaml:"user" env-required:"true"`
    // Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

type Database struct {
    Host     string `env:"POSTGRES_HOST" env-default:"127.0.0.1"`
    Port     string `env:"POSTGRES_PORT" env-default:"6500"`
    User     string `env:"POSTGRES_USER" env-default:"admin"`
    Password string `env:"POSTGRES_PASSWORD" env-default:"password123"`
    DBName   string `env:"POSTGRES_DB" env-default:"kode-test-project"`
}

type PgAdmin struct {
    DefaultEmail    string `env:"PGADMIN_DEFAULT_EMAIL" env-default:"admin@admin.com"`
    DefaultPassword string `env:"PGADMIN_DEFAULT_PASSWORD" env-default:"password123"`
}

func MustLoad() *Config {
    if err := godotenv.Load(".env"); err != nil {
        log.Fatal(err)
    }

    configPath := os.Getenv("CONFIG_PATH")
    if configPath == "" {
        log.Fatal("CONFIG_PATH is not set")
    }

    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        log.Fatalf("config file does not exist: %s", configPath)
    }

    var cfg Config

    if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
        log.Fatalf("cannot read config: %s", err)
    }

    return &cfg
}