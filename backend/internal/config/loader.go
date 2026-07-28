package config

import (
	"fmt"
	"os"
	"strconv"
)

// Load loads the application configuration.
func Load() (*Config, error) {

	// Load the .env file
	if err := LoadEnv(); err != nil {
		return nil, err
	}

	cfg := &Config{

		App: AppConfig{
			Name:        getEnv("APP_NAME"),
			Version:     getEnv("APP_VERSION"),
			Environment: getEnv("APP_ENV"),
		},

		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST"),
			Port:         getEnv("SERVER_PORT"),
			ReadTimeout:  getEnvAsInt("SERVER_READ_TIMEOUT"),
			WriteTimeout: getEnvAsInt("SERVER_WRITE_TIMEOUT"),
			IdleTimeout:  getEnvAsInt("SERVER_IDLE_TIMEOUT"),
		},

		Logger: LoggerConfig{
			Level:       getEnv("LOG_LEVEL"),
			Format:      getEnv("LOG_FORMAT"),
			Output:      getEnv("LOG_OUTPUT"),
			Development: getEnvAsBool("LOG_DEVELOPMENT"),
		},

		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST"),
			Port:            getEnv("DB_PORT"),
			User:            getEnv("DB_USER"),
			Password:        getEnv("DB_PASSWORD"),
			Name:            getEnv("DB_NAME"),
			SSLMode:         getEnv("DB_SSLMODE"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS"),
			ConnMaxLifetime: getEnvAsInt("DB_CONN_MAX_LIFETIME"),
		},

		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET"),
			Expiration: getEnvAsInt("JWT_EXPIRATION"),
		},

		LDAP: LDAPConfig{
			Host:     getEnv("LDAP_HOST"),
			Port:     getEnv("LDAP_PORT"),
			Username: getEnv("LDAP_USERNAME"),
			Password: getEnv("LDAP_PASSWORD"),
			BaseDN:   getEnv("LDAP_BASE_DN"),
			UseTLS:   getEnvAsBool("LDAP_USE_TLS"),
		},
	}

	if err := validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func validate(cfg *Config) error {

	if cfg.App.Name == "" {
		return fmt.Errorf("APP_NAME is required")
	}

	if cfg.Server.Port == "" {
		return fmt.Errorf("SERVER_PORT is required")
	}

	if cfg.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	return nil
}

func getEnv(key string) string {
	return os.Getenv(key)
}

func getEnvAsInt(key string) int {

	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return 0
	}

	return value
}

func getEnvAsBool(key string) bool {

	value, err := strconv.ParseBool(os.Getenv(key))
	if err != nil {
		return false
	}

	return value
}