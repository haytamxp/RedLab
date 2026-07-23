package config

type Config struct {
	App	AppConfig
	Server ServerConfig
	Logger LoggerConfig
	JWT JWTConfig
	LDAP LDAPConfig
	Database DatabaseConfig
}

type AppConfig struct {
	Name string
	Version string
	Environment string
}

type ServerConfig struct {
    Host         string
    Port         string
    ReadTimeout  int
    WriteTimeout int
    IdleTimeout  int
}

type LoggerConfig struct {
    Level       string
    Format      string
    Output      string
    Development bool
}

type DatabaseConfig struct {
    Host            string
    Port            string
    User            string
    Password        string
    Name            string
    SSLMode         string
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime int
}

type LDAPConfig struct {
    Host     string
    Port     string
    Username string
    Password string
    BaseDN   string
    UseTLS   bool
}

type JWTConfig struct {
    Secret     string
    Expiration int
}