package ldap

import (
	"fmt"
	"strings"
)

func (c *Client) Authenticate(username, password string) (*User, error) {

	username = strings.TrimSpace(username)

	if username == "" {
		return nil, fmt.Errorf("username is empty")
	}

	if password == "" {
		return nil, fmt.Errorf("password is empty")
	}

	// The service account searches AD for the user.
	user, err := c.FindUser(username)
	if err != nil {
		return nil, err
	}

	if !user.Enabled {
		return nil, fmt.Errorf("LDAP account is disabled")
	}

	// Use a separate LDAP connection for the user's authentication bind.
	userClient := NewClient(c.cfg)
	defer userClient.Close()

	if err := userClient.BindUser(user.DN, password); err != nil {
		return nil, fmt.Errorf("LDAP authentication failed: %w", err)
	}

	return user, nil
}
