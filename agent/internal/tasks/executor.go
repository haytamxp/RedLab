package tasks

import (
	"context"
	"encoding/json"

	"github.com/haytamxp/redlab/agent/internal/api"
	"github.com/haytamxp/redlab/agent/internal/modules"
)

type Executor struct {
	registry *modules.Registry
}

func NewExecutor(
	registry *modules.Registry,
) *Executor {
	return &Executor{
		registry: registry,
	}
}

func (e *Executor) Execute(
	ctx context.Context,
	task *api.Task,
) api.TaskResultRequest {
	if e.registry == nil {
		return api.TaskResultRequest{
			Status: "FAILED",
			Result: json.RawMessage(`{}`),
			Error:  "executor is not configured",
		}
	}

	result, err := e.registry.Execute(
		ctx,
		task.Type,
		task.Payload,
	)

	if err != nil {
		errorResult := map[string]any{
			"task_type": task.Type,
			"error":     err.Error(),
		}

		payload, marshalErr := json.Marshal(
			errorResult,
		)

		if marshalErr != nil {
			payload = json.RawMessage(`{}`)
		}

		return api.TaskResultRequest{
			Status: "FAILED",
			Result: payload,
			Error:  err.Error(),
		}
	}

	payload, err := json.Marshal(result)
	if err != nil {
		return api.TaskResultRequest{
			Status: "FAILED",
			Result: json.RawMessage(`{}`),
			Error:  "failed to serialize module result",
		}
	}

	return api.TaskResultRequest{
		Status: "COMPLETED",
		Result: payload,
	}
}
