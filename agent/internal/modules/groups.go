package modules

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-ldap/ldap/v3"

	"github.com/haytamxp/redlab/agent/internal/mitre"
)

type GroupsModule struct {
	ldap *LDAPClient
}

func NewGroupsModule(
	client *LDAPClient,
) *GroupsModule {
	return &GroupsModule{
		ldap: client,
	}
}

func (m *GroupsModule) Name() string {
	return "ad_group_enumeration"
}

func (m *GroupsModule) Technique() mitre.Technique {
	return mitre.Technique{
		ID:   "T1069.002",
		Name: "Permission Groups Discovery: Domain Groups",
	}
}

func (m *GroupsModule) Execute(
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

	filter := "(&(objectClass=group)(groupType=*))"

	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter,
		[]string{
			"cn",
			"sAMAccountName",
			"description",
			"groupType",
		},
		nil,
	)

	result, err := conn.SearchWithPaging(
		searchRequest,
		1000,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"LDAP group enumeration failed: %w",
			err,
		)
	}

	groups := make(
		[]map[string]string,
		0,
		len(result.Entries),
	)

	for _, entry := range result.Entries {
		groups = append(
			groups,
			map[string]string{
				"cn": entry.GetAttributeValue(
					"cn",
				),
				"sam_account_name": entry.GetAttributeValue(
					"sAMAccountName",
				),
				"description": entry.GetAttributeValue(
					"description",
				),
				"group_type": entry.GetAttributeValue(
					"groupType",
				),
			},
		)
	}

	return map[string]any{
		"count":  len(groups),
		"groups": groups,
	}, nil
}
