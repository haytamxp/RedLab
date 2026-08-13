package modules

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-ldap/ldap/v3"

	"github.com/haytamxp/redlab/agent/internal/mitre"
)

type UsersModule struct {
	ldap *LDAPClient
}

func NewUsersModule(
	client *LDAPClient,
) *UsersModule {
	return &UsersModule{
		ldap: client,
	}
}

func (m *UsersModule) Name() string {
	return "ad_user_enumeration"
}

func (m *UsersModule) Technique() mitre.Technique {
	return mitre.Technique{
		ID:   "T1087.002",
		Name: "Domain Account Discovery",
	}
}

func (m *UsersModule) Execute(
	ctx context.Context,
	payload json.RawMessage,
) (any, error) {
	_ = ctx

	var options struct {
		BaseDN string `json:"base_dn"`
	}

	if len(payload) > 0 {
		if err := json.Unmarshal(
			payload,
			&options,
		); err != nil {
			return nil, fmt.Errorf(
				"invalid task payload: %w",
				err,
			)
		}
	}

	baseDN := options.BaseDN

	if baseDN == "" {
		baseDN = m.ldap.baseDN()
	}

	if baseDN == "" {
		return nil, fmt.Errorf(
			"base DN is not configured",
		)
	}

	conn, err := m.ldap.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	filter := "(&(objectCategory=person)(objectClass=user))"

	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter,
		[]string{
			"sAMAccountName",
			"userPrincipalName",
			"displayName",
			"userAccountControl",
		},
		nil,
	)

	result, err := conn.SearchWithPaging(
		searchRequest,
		1000,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"LDAP user enumeration failed: %w",
			err,
		)
	}

	users := make([]map[string]string, 0, len(result.Entries))

	for _, entry := range result.Entries {
		users = append(
			users,
			map[string]string{
				"sam_account_name": entry.GetAttributeValue(
					"sAMAccountName",
				),
				"user_principal_name": entry.GetAttributeValue(
					"userPrincipalName",
				),
				"display_name": entry.GetAttributeValue(
					"displayName",
				),
				"user_account_control": entry.GetAttributeValue(
					"userAccountControl",
				),
			},
		)
	}

	return map[string]any{
		"count": len(users),
		"users": users,
	}, nil
}
