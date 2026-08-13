package modules

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-ldap/ldap/v3"

	"github.com/haytamxp/redlab/agent/internal/mitre"
)

type ComputersModule struct {
	ldap *LDAPClient
}

func NewComputersModule(
	client *LDAPClient,
) *ComputersModule {
	return &ComputersModule{
		ldap: client,
	}
}

func (m *ComputersModule) Name() string {
	return "ad_computer_enumeration"
}

func (m *ComputersModule) Technique() mitre.Technique {
	return mitre.Technique{
		ID:   "T1018",
		Name: "Remote System Discovery",
	}
}

func (m *ComputersModule) Execute(
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

	filter := "(&(objectCategory=computer)(objectClass=computer))"

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
			"dNSHostName",
			"operatingSystem",
			"operatingSystemVersion",
		},
		nil,
	)

	result, err := conn.SearchWithPaging(
		searchRequest,
		1000,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"LDAP computer enumeration failed: %w",
			err,
		)
	}

	computers := make(
		[]map[string]string,
		0,
		len(result.Entries),
	)

	for _, entry := range result.Entries {
		computers = append(
			computers,
			map[string]string{
				"sam_account_name": entry.GetAttributeValue(
					"sAMAccountName",
				),
				"dns_hostname": entry.GetAttributeValue(
					"dNSHostName",
				),
				"operating_system": entry.GetAttributeValue(
					"operatingSystem",
				),
				"operating_system_version": entry.GetAttributeValue(
					"operatingSystemVersion",
				),
			},
		)
	}

	return map[string]any{
		"count":     len(computers),
		"computers": computers,
	}, nil
}
