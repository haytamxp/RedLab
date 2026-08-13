package ldap

import (
	"fmt"
	"strings"

	goldap "github.com/go-ldap/ldap/v3"
)

type User struct {
	DN                string
	SAMAccountName    string
	UserPrincipalName string
	Email             string
	FirstName         string
	LastName          string
	DisplayName       string
	Enabled           bool
	Groups            []string
}

func (c *Client) FindUser(username string) (*User, error) {
	if err := c.BindServiceAccount(); err != nil {
		return nil, err
	}

	username = strings.TrimSpace(username)

	if username == "" {
		return nil, fmt.Errorf("username is empty")
	}

	filter := fmt.Sprintf(
		"(&(objectCategory=person)(objectClass=user)(sAMAccountName=%s))",
		goldap.EscapeFilter(username),
	)

	request := goldap.NewSearchRequest(
		c.cfg.BaseDN,
		goldap.ScopeWholeSubtree,
		goldap.NeverDerefAliases,
		1,
		10,
		false,
		filter,
		[]string{
			"distinguishedName",
			"sAMAccountName",
			"userPrincipalName",
			"mail",
			"givenName",
			"sn",
			"displayName",
			"userAccountControl",
			"memberOf",
		},
		nil,
	)

	result, err := c.Search(request)

	if err != nil {
		return nil, fmt.Errorf("LDAP user search failed: %w", err)
	}

	if len(result.Entries) == 0 {
		return nil, fmt.Errorf("LDAP user not found: %s", username)
	}

	entry := result.Entries[0]

	return &User{
		DN:                entry.DN,
		SAMAccountName:    entry.GetAttributeValue("sAMAccountName"),
		UserPrincipalName: entry.GetAttributeValue("userPrincipalName"),
		Email:             entry.GetAttributeValue("mail"),
		FirstName:         entry.GetAttributeValue("givenName"),
		LastName:          entry.GetAttributeValue("sn"),
		DisplayName:       entry.GetAttributeValue("displayName"),
		Enabled:           accountEnabled(
			entry.GetAttributeValue("userAccountControl"),
		),
		Groups: entry.GetAttributeValues("memberOf"),
	}, nil
}

func accountEnabled(value string) bool {
	if value == "" {
		return true
	}

	var control int

	if _, err := fmt.Sscanf(value, "%d", &control); err != nil {
		return true
	}

	const accountDisable = 0x0002

	return control&accountDisable == 0
}