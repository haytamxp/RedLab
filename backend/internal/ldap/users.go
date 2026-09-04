package ldap

import (
	"fmt"
	"strings"

	goldap "github.com/go-ldap/ldap/v3"
)

func (c *Client) ListUsers(search string) ([]User, error) {
	if err := c.BindServiceAccount(); err != nil {
		return nil, err
	}

	filter := "(&(objectCategory=person)(objectClass=user))"
	search = strings.TrimSpace(search)
	if search != "" {
		value := goldap.EscapeFilter(search)
		filter = fmt.Sprintf(
			"(&(objectCategory=person)(objectClass=user)(|(sAMAccountName=%[1]s)(displayName=*%[1]s*)(mail=*%[1]s*)))",
			value,
		)
	}

	request := goldap.NewSearchRequest(
		c.cfg.BaseDN,
		goldap.ScopeWholeSubtree,
		goldap.NeverDerefAliases,
		100,
		10,
		false,
		filter,
		[]string{
			"distinguishedName", "sAMAccountName",
			"userPrincipalName", "mail", "givenName",
			"sn", "displayName", "userAccountControl",
			"memberOf",
		},
		nil,
	)

	result, err := c.Search(request)
	if err != nil {
		return nil, fmt.Errorf("LDAP user search failed: %w", err)
	}

	users := make([]User, 0, len(result.Entries))
	for _, entry := range result.Entries {
		users = append(users, User{
			DN:                entry.DN,
			SAMAccountName:    entry.GetAttributeValue("sAMAccountName"),
			UserPrincipalName: entry.GetAttributeValue("userPrincipalName"),
			Email:             entry.GetAttributeValue("mail"),
			FirstName:         entry.GetAttributeValue("givenName"),
			LastName:          entry.GetAttributeValue("sn"),
			DisplayName:       entry.GetAttributeValue("displayName"),
			Enabled:           accountEnabled(entry.GetAttributeValue("userAccountControl")),
			Groups:            entry.GetAttributeValues("memberOf"),
		})
	}

	return users, nil
}
