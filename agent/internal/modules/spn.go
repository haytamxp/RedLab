package modules

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-ldap/ldap/v3"

	"github.com/haytamxp/redlab/agent/internal/mitre"
)

type SPNModule struct {
	ldap *LDAPClient
}

func NewSPNModule(
	client *LDAPClient,
) *SPNModule {
	return &SPNModule{
		ldap: client,
	}
}

func (m *SPNModule) Name() string {
	return "spn_enumeration"
}

func (m *SPNModule) Technique() mitre.Technique {
	return mitre.Technique{
		ID:   "T1558.003",
		Name: "Steal or Forge Kerberos Tickets: Kerberoasting",
	}
}

func (m *SPNModule) Execute(
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

	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		"(&(objectClass=user)(servicePrincipalName=*))",
		[]string{
			"sAMAccountName",
			"userPrincipalName",
			"servicePrincipalName",
			"description",
		},
		nil,
	)

	result, err := conn.SearchWithPaging(
		searchRequest,
		1000,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"SPN enumeration failed: %w",
			err,
		)
	}

	accounts := make(
		[]map[string]any,
		0,
		len(result.Entries),
	)

	totalSPNs := 0

	for _, entry := range result.Entries {
		spns := entry.GetAttributeValues(
			"servicePrincipalName",
		)

		totalSPNs += len(spns)

		accounts = append(
			accounts,
			map[string]any{
				"sam_account_name": entry.GetAttributeValue(
					"sAMAccountName",
				),
				"user_principal_name": entry.GetAttributeValue(
					"userPrincipalName",
				),
				"description": entry.GetAttributeValue(
					"description",
				),
				"service_principal_names": spns,
			},
		)
	}

	return map[string]any{
		"account_count": len(accounts),
		"spn_count":     totalSPNs,
		"accounts":      accounts,
	}, nil
}
