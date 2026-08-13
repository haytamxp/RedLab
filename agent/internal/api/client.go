package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/haytamxp/redlab/agent/internal/auth"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	auth       *auth.TokenProvider
}

type AgentResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    Agent  `json:"data"`
}

type Agent struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	Hostname        string     `json:"hostname"`
	IPAddress       string     `json:"ip_address"`
	OperatingSystem string     `json:"operating_system"`
	Version         string     `json:"version"`
	Status          string     `json:"status"`
	LastSeen        *time.Time `json:"last_seen,omitempty"`
}

type TaskResponse struct {
	Success bool  `json:"success"`
	Data    *Task `json:"data"`
}

type Task struct {
	ID         string          `json:"id"`
	AgentID    string          `json:"agent_id"`
	Type       string          `json:"type"`
	Payload    json.RawMessage `json:"payload"`
	Status     string          `json:"status"`
	Priority   int             `json:"priority"`
}

type TaskResultRequest struct {
	Status string          `json:"status"`
	Result json.RawMessage `json:"result"`
	Error  string          `json:"error"`
}

func NewClient(
	baseURL string,
	provider *auth.TokenProvider,
) *Client {
	return &Client{
		baseURL: strings.TrimRight(
			baseURL,
			"/",
		),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		auth: provider,
	}
}

func (c *Client) Heartbeat(
	ctx context.Context,
	agentID string,
) (*AgentResponse, error) {
	var response AgentResponse

	err := c.doJSON(
		ctx,
		http.MethodPost,
		fmt.Sprintf(
			"/api/v1/agents/%s/heartbeat",
			agentID,
		),
		nil,
		&response,
	)

	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) NextTask(
	ctx context.Context,
	agentID string,
) (*Task, error) {
	var response TaskResponse

	err := c.doJSON(
		ctx,
		http.MethodPost,
		fmt.Sprintf(
			"/api/v1/agents/%s/tasks/next",
			agentID,
		),
		nil,
		&response,
	)

	if err != nil {
		if strings.Contains(
			err.Error(),
			"HTTP 204",
		) {
			return nil, nil
		}

		return nil, err
	}

	return response.Data, nil
}

func (c *Client) SubmitTaskResult(
	ctx context.Context,
	agentID string,
	taskID string,
	result TaskResultRequest,
) error {
	return c.doJSON(
		ctx,
		http.MethodPost,
		fmt.Sprintf(
			"/api/v1/agents/%s/tasks/%s/result",
			agentID,
			taskID,
		),
		result,
		nil,
	)
}

func (c *Client) doJSON(
	ctx context.Context,
	method string,
	path string,
	body any,
	responseTarget any,
) error {
	var requestBody *bytes.Reader

	if body == nil {
		requestBody = bytes.NewReader(nil)
	} else {
		payload, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf(
				"marshal request: %w",
				err,
			)
		}

		requestBody = bytes.NewReader(payload)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		method,
		c.baseURL+path,
		requestBody,
	)
	if err != nil {
		return fmt.Errorf(
			"create HTTP request: %w",
			err,
		)
	}

	request.Header.Set(
		"Authorization",
		c.auth.AuthorizationHeader(),
	)
	request.Header.Set(
		"Content-Type",
		"application/json",
	)
	request.Header.Set(
		"Accept",
		"application/json",
	)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf(
			"HTTP request failed: %w",
			err,
		)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNoContent {
		return fmt.Errorf("HTTP 204 no content")
	}

	if response.StatusCode < 200 ||
		response.StatusCode >= 300 {
		return fmt.Errorf(
			"HTTP %d from %s",
			response.StatusCode,
			path,
		)
	}

	if responseTarget == nil {
		return nil
	}

	if err := json.NewDecoder(
		response.Body,
	).Decode(responseTarget); err != nil {
		return fmt.Errorf(
			"decode response: %w",
			err,
		)
	}

	return nil
}