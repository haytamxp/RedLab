package modules

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/haytamxp/redlab/agent/internal/config"
	"github.com/haytamxp/redlab/agent/internal/mitre"
)

var ErrUnsupportedTask = errors.New(
	"unsupported assessment task",
)

type Module interface {
	Name() string
	Technique() mitre.Technique
	Execute(
		ctx context.Context,
		payload json.RawMessage,
	) (any, error)
}

type Registry struct {
	modules map[string]Module
}

func NewRegistry(
	ldapConfig config.LDAPConfig,
) *Registry {
	registry := &Registry{
		modules: make(map[string]Module),
	}

	ldapClient := NewLDAPClient(
		ldapConfig,
	)

	registry.Register(
		"DOMAIN_INFO",
		NewDomainInfoModule(ldapClient),
	)

	registry.Register(
		"AD_USER_ENUMERATION",
		NewUsersModule(ldapClient),
	)

	registry.Register(
		"AD_GROUP_ENUMERATION",
		NewGroupsModule(ldapClient),
	)

	return registry
}

func (r *Registry) Register(
	taskType string,
	module Module,
) {
	taskType = strings.ToUpper(
		strings.TrimSpace(taskType),
	)

	r.modules[taskType] = module
}

func (r *Registry) Execute(
	ctx context.Context,
	taskType string,
	payload json.RawMessage,
) (map[string]any, error) {
	taskType = strings.ToUpper(
		strings.TrimSpace(taskType),
	)

	module, ok := r.modules[taskType]
	if !ok {
		return nil, fmt.Errorf(
			"%w: %s",
			ErrUnsupportedTask,
			taskType,
		)
	}

	data, err := module.Execute(
		ctx,
		payload,
	)
	if err != nil {
		return nil, err
	}

	technique := module.Technique()

	return map[string]any{
		"module":    module.Name(),
		"task_type": taskType,
		"mitre_technique": map[string]string{
			"id":   technique.ID,
			"name": technique.Name,
		},
		"data": data,
	}, nil
}
