package config

import (
	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port string
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	URL string
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret string
	Issuer string
}

// LoadConfig loads the application configuration from environment variables or config file
func LoadConfig() (*Config, error) {
	// Set defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("jwt.issuer", "ride-sharing-app")

	// Look for config files
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Read environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")

	// Map environment variables
	viper.BindEnv("server.port", "APP_SERVER_PORT")
	viper.BindEnv("database.url", "APP_DB_URL")
	viper.BindEnv("jwt.secret", "APP_JWT_SECRET")
	viper.BindEnv("jwt.issuer", "APP_JWT_ISSUER")

	// Read config file (if exists)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
		// Config file not found, will use defaults and env vars
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
