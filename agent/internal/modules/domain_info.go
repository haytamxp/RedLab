package modules

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-ldap/ldap/v3"

	"github.com/haytamxp/redlab/agent/internal/mitre"
)

type DomainInfoModule struct {
	ldap *LDAPClient
}

func NewDomainInfoModule(
	client *LDAPClient,
) *DomainInfoModule {
	return &DomainInfoModule{
		ldap: client,
	}
}

func (m *DomainInfoModule) Name() string {
	return "domain_info"
}

func (m *DomainInfoModule) Technique() mitre.Technique {
	return mitre.Technique{
		ID:   "T1482",
		Name: "Domain Trust Discovery",
	}
}

func (m *DomainInfoModule) Execute(
	ctx context.Context,
	payload json.RawMessage,
) (any, error) {
	_ = ctx
	_ = payload

	conn, err := m.ldap.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	searchRequest := ldap.NewSearchRequest(
		"",
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		"(objectClass=*)",
		[]string{
			"defaultNamingContext",
			"dnsHostName",
			"rootDomainNamingContext",
			"forestFunctionality",
			"domainFunctionality",
		},
		nil,
	)

	result, err := conn.Search(
		searchRequest,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"RootDSE query failed: %w",
			err,
		)
	}

	if len(result.Entries) == 0 {
		return nil, fmt.Errorf(
			"RootDSE returned no entries",
		)
	}

	entry := result.Entries[0]

	return map[string]any{
		"default_naming_context": entry.GetAttributeValue(
			"defaultNamingContext",
		),
		"dns_hostname": entry.GetAttributeValue(
			"dnsHostName",
		),
		"root_domain_naming_context": entry.GetAttributeValue(
			"rootDomainNamingContext",
		),
		"forest_functionality": entry.GetAttributeValue(
			"forestFunctionality",
		),
		"domain_functionality": entry.GetAttributeValue(
			"domainFunctionality",
		),
	}, nil
}
