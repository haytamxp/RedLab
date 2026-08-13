package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

type LDAPConfig struct {
	URL      string
	Username string
	Password string
	BaseDN   string
}

type Config struct {
	ServerURL    string
	AgentID      string
	AgentToken   string
	Heartbeat    time.Duration
	PollInterval time.Duration
	LDAP         LDAPConfig
}

func Load() (Config, error) {
	cfg := Config{
		ServerURL: strings.TrimRight(
			os.Getenv("REDLAB_SERVER_URL"),
			"/",
		),
		AgentID: strings.TrimSpace(
			os.Getenv("REDLAB_AGENT_ID"),
		),
		AgentToken: strings.TrimSpace(
			os.Getenv("REDLAB_AGENT_TOKEN"),
		),
		Heartbeat: loadDuration(
			"REDLAB_HEARTBEAT_SECONDS",
			30,
		),
		PollInterval: loadDuration(
			"REDLAB_POLL_SECONDS",
			10,
		),
		LDAP: LDAPConfig{
			URL: strings.TrimRight(
				strings.TrimSpace(
					os.Getenv("REDLAB_LDAP_URL"),
				),
				"/",
			),
			Username: strings.TrimSpace(
				os.Getenv("REDLAB_LDAP_USERNAME"),
			),
			Password: os.Getenv("REDLAB_LDAP_PASSWORD"),
			BaseDN: strings.TrimSpace(
				os.Getenv("REDLAB_LDAP_BASE_DN"),
			),
		},
	}

	if cfg.ServerURL == "" {
		return Config{}, errors.New(
			"REDLAB_SERVER_URL is required",
		)
	}

	if cfg.AgentID == "" {
		return Config{}, errors.New(
			"REDLAB_AGENT_ID is required",
		)
	}

	if cfg.AgentToken == "" {
		return Config{}, errors.New(
			"REDLAB_AGENT_TOKEN is required",
		)
	}

	return cfg, nil
}

func loadDuration(
	name string,
	defaultSeconds int,
) time.Duration {
	value := strings.TrimSpace(
		os.Getenv(name),
	)

	if value == "" {
		value = strconv.Itoa(defaultSeconds)
	}

	seconds, err := strconv.Atoi(value)
	if err != nil || seconds <= 0 {
		seconds = defaultSeconds
	}

	return time.Duration(seconds) * time.Second
}
