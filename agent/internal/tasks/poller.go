package tasks

import (
	"context"
	"log"
	"time"

	"github.com/haytamxp/redlab/agent/internal/api"
	"github.com/haytamxp/redlab/agent/internal/config"
	"github.com/haytamxp/redlab/agent/internal/modules"
)

type Poller struct {
	client       *api.Client
	agentID      string
	pollInterval time.Duration
	executor     *Executor
}

func NewPoller(
	client *api.Client,
	agentID string,
	pollInterval time.Duration,
) *Poller {
	cfg, err := config.Load()
	if err != nil {
		log.Printf(
			"[tasks] executor configuration warning: %v",
			err,
		)

		return &Poller{
			client:       client,
			agentID:      agentID,
			pollInterval: pollInterval,
			executor:     NewExecutor(nil),
		}
	}

	registry := modules.NewRegistry(cfg.LDAP)

	return &Poller{
		client:       client,
		agentID:      agentID,
		pollInterval: pollInterval,
		executor:     NewExecutor(registry),
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

	result := p.executor.Execute(ctx, task)

	if err := p.client.SubmitTaskResult(
		ctx,
		p.agentID,
		task.ID,
		result,
	); err != nil {
		log.Printf(
			"[tasks] submit result failed task=%s: %v",
			task.ID,
			err,
		)
		return
	}

	log.Printf(
		"[tasks] task=%s completed with status=%s",
		task.ID,
		result.Status,
	)
}
