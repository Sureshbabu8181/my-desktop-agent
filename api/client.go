package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	serverURL  string
	connectKey string
	httpClient *http.Client
}

type DeviceInfo struct {
	Hostname     string `json:"hostname"`
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
	// Add more device details
}

type SystemMetrics struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	// Add more metrics
}

type Command struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

type Policy struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`  // e.g., "password_policy", "firewall_rule"
	Value string `json:"value"` // Policy configuration
}

func NewClient(serverURL, connectKey string) *Client {
	return &Client{
		serverURL:  serverURL,
		connectKey: connectKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) sendRequest(ctx context.Context, method, path string, payload interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.serverURL, path)

	var reqBody []byte
	if payload != nil {
		var err error
		reqBody, err = json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request payload: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Connect-Key", c.connectKey) // JumpCloud uses this

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned non-OK status: %s", resp.Status)
	}

	body, err := os.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return body, nil
}

func (c *Client) RegisterDevice(ctx context.Context, info *DeviceInfo) error {
	_, err := c.sendRequest(ctx, "POST", "/register", info)
	return err
}

func (c *Client) ReportMetrics(ctx context.Context, metrics *SystemMetrics) error {
	_, err := c.sendRequest(ctx, "POST", "/metrics", metrics)
	return err
}

func (c *Client) FetchCommands(ctx context.Context) ([]Command, error) {
	body, err := c.sendRequest(ctx, "GET", "/commands", nil)
	if err != nil {
		return nil, err
	}
	var commands []Command
	if err := json.Unmarshal(body, &commands); err != nil {
		return nil, fmt.Errorf("failed to unmarshal commands: %w", err)
	}
	return commands, nil
}

func (c *Client) ReportCommandResult(ctx context.Context, commandID, status, output string) error {
	result := map[string]string{
		"id":     commandID,
		"status": status,
		"output": output,
	}
	_, err := c.sendRequest(ctx, "POST", "/command-result", result)
	return err
}

func (c *Client) FetchPolicies(ctx context.Context) ([]Policy, error) {
	body, err := c.sendRequest(ctx, "GET", "/policies", nil)
	if err != nil {
		return nil, err
	}
	var policies []Policy
	if err := json.Unmarshal(body, &policies); err != nil {
		return nil, fmt.Errorf("failed to unmarshal policies: %w", err)
	}
	return policies, nil
}
