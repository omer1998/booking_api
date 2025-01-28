package config

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Env string

type Config struct {
	DbName     string `env:"DB_NAME"`
	DbUser     string `env:"DB_USER"`
	DbPassword string `env:"DB_PASSWORD"`
	DbHost     string `env:"DB_HOST"`
	DbPort     string `env:"DB_PORT"`
	DbPortTest string `env:"DB_PORT_TEST"`
	Env        Env    `env:"ENV" envDefault:"dev"`
}

const (
	DevEnv  Env = "dev"
	TestEnv Env = "test"
)

// postgres://omer:12345678@127.0.0.1:5433/booking?sslmode=disable
func (c *Config) GetConnectionDbUrl() string {
	// localDbUrl:= "postgres://username:password@host:port/dbName"
	if c.Env == "test" {
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.DbUser, c.DbPassword, c.DbHost, c.DbPortTest, c.DbName)
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.DbUser, c.DbPassword, c.DbHost, c.DbPort, c.DbName)
}
func New() (*Config, error) {
	envErr := godotenv.Load()
	if envErr != nil {
		panic("kssssssss")
	}
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
func NewWithEnvPath(path string) (*Config, error) {
	envErr := godotenv.Load(path)
	if envErr != nil {
		panic("kssssssss")
	}
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
func PrintConfig(cfg *Config) {
	fmt.Println("DB_NAME:", cfg.DbName)
	fmt.Println("DB_USER:", cfg.DbUser)
	fmt.Println("DB_PASSWORD:", cfg.DbPassword)
	fmt.Println("DB_HOST:", cfg.DbHost)
	fmt.Println("DB_PORT:", cfg.DbPort)
}
