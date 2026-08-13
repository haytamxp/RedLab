package heartbeat

import (
	"context"
	"log"
	"time"

	"github.com/haytamxp/redlab/agent/internal/api"
)

type Service struct {
	client   *api.Client
	agentID  string
	interval time.Duration
}

func New(
	client *api.Client,
	agentID string,
	interval time.Duration,
) *Service {
	return &Service{
		client:   client,
		agentID:  agentID,
		interval: interval,
	}
}

func (s *Service) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.send(ctx)

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			s.send(ctx)
		}
	}
}

func (s *Service) send(ctx context.Context) {
	response, err := s.client.Heartbeat(
		ctx,
		s.agentID,
	)

	if err != nil {
		log.Printf(
			"[heartbeat] failed: %v",
			err,
		)
		return
	}

	log.Printf(
		"[heartbeat] agent=%s status=%s",
		response.Data.ID,
		response.Data.Status,
	)
}