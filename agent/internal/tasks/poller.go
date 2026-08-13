package tasks

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/haytamxp/redlab/agent/internal/api"
)

type Poller struct {
	client       *api.Client
	agentID      string
	pollInterval time.Duration
}

func NewPoller(
	client *api.Client,
	agentID string,
	pollInterval time.Duration,
) *Poller {
	return &Poller{
		client:       client,
		agentID:      agentID,
		pollInterval: pollInterval,
	}
}

func (p *Poller) Start(ctx context.Context) {
	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()

	p.poll(ctx)

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			p.poll(ctx)
		}
	}
}

func (p *Poller) poll(ctx context.Context) {
	task, err := p.client.NextTask(
		ctx,
		p.agentID,
	)

	if err != nil {
		log.Printf(
			"[tasks] poll failed: %v",
			err,
		)
		return
	}

	if task == nil {
		return
	}

	log.Printf(
		"[tasks] received id=%s type=%s priority=%d",
		task.ID,
		task.Type,
		task.Priority,
	)

	result := json.RawMessage(`{
		"message": "Task received but executor is not implemented yet"
	}`)

	err = p.client.SubmitTaskResult(
		ctx,
		p.agentID,
		task.ID,
		api.TaskResultRequest{
			Status: "FAILED",
			Result: result,
			Error:  "task executor is not implemented in this agent build",
		},
	)

	if err != nil {
		log.Printf(
			"[tasks] submit result failed: %v",
			err,
		)
		return
	}

	log.Printf(
		"[tasks] task=%s marked FAILED",
		task.ID,
	)
}
