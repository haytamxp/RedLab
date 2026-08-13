package modules

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"

	"github.com/haytamxp/redlab/agent/internal/config"
)

type LDAPClient struct {
	config config.LDAPConfig
}

func NewLDAPClient(
	cfg config.LDAPConfig,
) *LDAPClient {
	return &LDAPClient{
		config: cfg,
	}
}

func (c *LDAPClient) connect() (*ldap.Conn, error) {
	if c.config.URL == "" {
		return nil, fmt.Errorf(
			"REDLAB_LDAP_URL is required",
		)
	}

	conn, err := ldap.DialURL(
		c.config.URL,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"LDAP connection failed: %w",
			err,
		)
	}

	if c.config.Username == "" ||
		c.config.Password == "" {
		conn.Close()

		return nil, fmt.Errorf(
			"LDAP service credentials are not configured",
		)
	}

	if err := conn.Bind(
		c.config.Username,
		c.config.Password,
	); err != nil {
		conn.Close()

		return nil, fmt.Errorf(
			"LDAP bind failed: %w",
			err,
		)
	}

	return conn, nil
}

func (c *LDAPClient) baseDN() string {
	return c.config.BaseDN
}
