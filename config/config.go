package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config represents the top-level structure of the configuration file.
type Config struct {
	Database *DatabaseConfig `yaml:"database"`
	Api      *ApiConfig      `yaml:"api"`
}

// DatabaseConfig represents the database configuration section of the configuration file.
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"ssl_mode"`
}

// ConnString generates a string to pass to sql.Open to connect to the PostgreSQL
// database.
func (dc *DatabaseConfig) ConnString() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s database=%s sslmode=%s",
		dc.Host, dc.User, dc.Password, dc.Database, dc.SSLMode,
	)
}

// DatabaseConfig represents the API configuration section of the configuration file.
type ApiConfig struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

// NewConfig reads and parses the content at the provided file path into an instance of
// the Config struct.
func NewConfig(filePath string) (*Config, error) {
	raw, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	cfg := new(Config)
	if err = yaml.Unmarshal(raw, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
