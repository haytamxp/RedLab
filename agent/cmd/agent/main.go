package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/haytamxp/redlab/agent/internal/api"
	"github.com/haytamxp/redlab/agent/internal/auth"
	"github.com/haytamxp/redlab/agent/internal/config"
	"github.com/haytamxp/redlab/agent/internal/heartbeat"
	"github.com/haytamxp/redlab/agent/internal/tasks"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf(
			"[agent] configuration error: %v",
			err,
		)
	}

	provider := auth.NewTokenProvider(
		cfg.AgentToken,
	)

	client := api.NewClient(
		cfg.ServerURL,
		provider,
	)

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	heartbeatService := heartbeat.New(
		client,
		cfg.AgentID,
		cfg.Heartbeat,
	)

	poller := tasks.NewPoller(
		client,
		cfg.AgentID,
		cfg.PollInterval,
	)

	log.Printf(
		"[agent] starting agent=%s",
		cfg.AgentID,
	)

	go heartbeatService.Start(ctx)
	go poller.Start(ctx)

	<-ctx.Done()

	log.Println("[agent] shutting down")
}