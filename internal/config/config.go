package config

import (
	"flag"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	address       = "localhost"
	timeout       = 4
	idelTimeout   = 30
	port          = 8080
	loggerLevel   = "debug"
	translateMode = false
	envFile       = ".env"
)

type Config struct {
	Server struct {
		Address     string        `yaml:"address"`
		Timeout     time.Duration `yaml:"timeout"`
		IdleTimeout time.Duration `yaml:"timeout"`
	} `yaml:"server"`
	Database struct {
		User    string `yaml:"user"`
		DbName  string `yaml:"dbname"`
		Host    string `yaml:"host"`
		Port    uint64 `yaml:"port"`
		SslMode string `yaml:"sslmode"`
	} `yaml:"database"`
	EnvFile        string `yaml:"env_file"`
	LoggerLvl      string `yaml:"logger_level"`
	CookieSettings struct {
		Secure     bool `yaml:"secure"`
		HttpOnly   bool `yaml:"http_only"`
		ExpireDate struct {
			Years  uint64 `yaml:"years"`
			Months uint64 `yaml:"months"`
			Days   uint64 `yaml:"days"`
		} `yaml:"expire_date"`
	}
}

func NewConfig() *Config {
	return &Config{
		Server: struct {
			Address     string        `yaml:"address"`
			Timeout     time.Duration `yaml:"timeout"`
			IdleTimeout time.Duration `yaml:"timeout"`
		}(struct {
			Address     string
			Timeout     time.Duration
			IdleTimeout time.Duration
		}{
			Address:     address,
			Timeout:     time.Duration(timeout),
			IdleTimeout: time.Duration(idelTimeout),
		}),
		Database: struct {
			User    string `yaml:"user"`
			DbName  string `yaml:"dbname"`
			Host    string `yaml:"host"`
			Port    uint64 `yaml:"port"`
			SslMode string `yaml:"sslmode"`
		}(struct {
			User    string
			DbName  string
			Host    string
			Port    uint64
			SslMode string
		}{
			User:    "postgres",
			DbName:  "cyber_garden",
			Host:    "localhost",
			Port:    5432,
			SslMode: "disable",
		}),
		EnvFile:   envFile,
		LoggerLvl: loggerLevel,
		CookieSettings: struct {
			Secure     bool `yaml:"secure"`
			HttpOnly   bool `yaml:"http_only"`
			ExpireDate struct {
				Years  uint64 `yaml:"years"`
				Months uint64 `yaml:"months"`
				Days   uint64 `yaml:"days"`
			} `yaml:"expire_date"`
		}(struct {
			Secure     bool
			HttpOnly   bool
			ExpireDate struct {
				Years  uint64
				Months uint64
				Days   uint64
			}
		}{
			Secure:   true,
			HttpOnly: true,
			ExpireDate: struct {
				Years  uint64
				Months uint64
				Days   uint64
			}{
				Years:  0,
				Months: 0,
				Days:   7,
			},
		}),
	}
}

func (c *Config) Open(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = yaml.NewDecoder(file).Decode(c); err != nil {
		return err
	}

	return nil
}

func ParseFlag(path *string) {
	flag.StringVar(path, "ConfigPath", "configs/app/local.yaml", "Path to Config")
}
