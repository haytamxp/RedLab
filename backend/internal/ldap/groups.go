package ldap

import (
	"fmt"
	"strings"

	goldap "github.com/go-ldap/ldap/v3"
)

type Group struct {
	DN          string
	Name        string
	Description string
}

func (c *Client) FindUserGroups(userDN string) ([]Group, error) {
	if err := c.BindServiceAccount(); err != nil {
		return nil, err
	}

	userDN = strings.TrimSpace(userDN)

	if userDN == "" {
		return nil, fmt.Errorf("user DN is empty")
	}

	filter := fmt.Sprintf(
		"(&(objectClass=group)(member=%s))",
		goldap.EscapeFilter(userDN),
	)

	request := goldap.NewSearchRequest(
		c.cfg.BaseDN,
		goldap.ScopeWholeSubtree,
		goldap.NeverDerefAliases,
		0,
		10,
		false,
		filter,
		[]string{
			"distinguishedName",
			"sAMAccountName",
			"name",
			"description",
		},
		nil,
	)

	result, err := c.Search(request)

	if err != nil {
		return nil, fmt.Errorf("LDAP group search failed: %w", err)
	}

	groups := make([]Group, 0, len(result.Entries))

	for _, entry := range result.Entries {
		groups = append(groups, Group{
			DN:          entry.DN,
			Name:        entry.GetAttributeValue("name"),
			Description: entry.GetAttributeValue("description"),
		})
	}

	return groups, nil
}