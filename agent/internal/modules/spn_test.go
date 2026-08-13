package modules

import (
	"testing"

	"github.com/haytamxp/redlab/agent/internal/config"
)

func TestSPNModuleMetadata(t *testing.T) {
	module := NewSPNModule(
		NewLDAPClient(
			config.LDAPConfig{},
		),
	)

	if module.Name() != "spn_enumeration" {
		t.Fatalf(
			"unexpected module name: %s",
			module.Name(),
		)
	}

	technique := module.Technique()

	if technique.ID != "T1558.003" {
		t.Fatalf(
			"unexpected technique ID: %s",
			technique.ID,
		)
	}
}
