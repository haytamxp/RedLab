package ldap

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"

	goldap "github.com/go-ldap/ldap/v3"

	"github.com/haytamxp/redlab/backend/internal/config"
)

type Client struct {
	conn *goldap.Conn
	cfg  config.LDAPConfig
}

func NewClient(cfg config.LDAPConfig) *Client {
	return &Client{
		cfg: cfg,
	}
}

func (c *Client) Connect() error {
	if c.conn != nil {
		return nil
	}

	host := strings.TrimSpace(c.cfg.Host)
	port := strings.TrimSpace(c.cfg.Port)

	if host == "" {
		return fmt.Errorf("LDAP host is empty")
	}

	if port == "" {
		return fmt.Errorf("LDAP port is empty")
	}

	address := net.JoinHostPort(host, port)

	var (
		conn *goldap.Conn
		err  error
	)

	if c.cfg.UseTLS {
		conn, err = goldap.DialURL(
			"ldaps://"+address,
			goldap.DialWithTLSConfig(&tls.Config{
				MinVersion: tls.VersionTLS12,
				ServerName: host,
			}),
		)
	} else {
		conn, err = goldap.DialURL(
			"ldap://"+address,
			goldap.DialWithDialer(
				&net.Dialer{
					Timeout: 5 * time.Second,
				},
			),
		)
	}

	if err != nil {
		return fmt.Errorf("LDAP connection failed: %w", err)
	}

	c.conn = conn

	return nil
}

func (c *Client) BindServiceAccount() error {
	if err := c.Connect(); err != nil {
		return err
	}

	username := strings.TrimSpace(c.cfg.Username)

	if username == "" {
		return fmt.Errorf("LDAP_USERNAME is empty")
	}

	if c.cfg.Password == "" {
		return fmt.Errorf("LDAP_PASSWORD is empty")
	}

	if err := c.conn.Bind(username, c.cfg.Password); err != nil {
		return fmt.Errorf("LDAP service account bind failed: %w", err)
	}

	return nil
}

func (c *Client) BindUser(bindDN, password string) error {
	if err := c.Connect(); err != nil {
		return err
	}

	bindDN = strings.TrimSpace(bindDN)

	if bindDN == "" {
		return fmt.Errorf("LDAP user DN is empty")
	}

	if password == "" {
		return fmt.Errorf("LDAP user password is empty")
	}

	if err := c.conn.Bind(bindDN, password); err != nil {
		return fmt.Errorf("LDAP user bind failed: %w", err)
	}

	return nil
}

func (c *Client) Search(request *goldap.SearchRequest) (*goldap.SearchResult, error) {
	if c.conn == nil {
		return nil, fmt.Errorf("LDAP connection is not established")
	}

	return c.conn.Search(request)
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}