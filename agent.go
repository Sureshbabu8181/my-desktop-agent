package agent

import (
	"context"
	"time"

	"github.com/yourusername/my-desktop-agent/api"
	"github.com/yourusername/my-desktop-agent/config"
	"github.com/yourusername/my-desktop-agent/logger"
	"github.com/yourusername/my-desktop-agent/system" // Import your OS-specific system package
)

type Agent struct {
	cfg    *config.Config
	client *api.Client // HTTP client for server communication
}

func NewAgent(cfg *config.Config) *Agent {
	return &Agent{
		cfg:    cfg,
		client: api.NewClient(cfg.ServerURL, cfg.ConnectKey),
	}
}

func (a *Agent) Start(ctx context.Context) error {
	logger.Info("Agent starting...")

	// Initial registration/check-in
	err := a.registerDevice(ctx)
	if err != nil {
		logger.Error("Initial device registration failed: %v", err)
		return err
	}
	logger.Info("Device registered successfully.")

	// Main agent loop
	ticker := time.NewTicker(time.Duration(a.cfg.PollInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Agent main loop stopping due to context cancellation.")
			return nil
		case <-ticker.C:
			a.performRoutineTasks(ctx)
		}
	}
}

func (a *Agent) registerDevice(ctx context.Context) error {
	// Get device information
	deviceInfo, err := system.GetDeviceInfo()
	if err != nil {
		return err
	}

	// Send registration request to server
	return a.client.RegisterDevice(ctx, deviceInfo)
}

func (a *Agent) performRoutineTasks(ctx context.Context) {
	logger.Info("Performing routine agent tasks...")

	// 1. Report system metrics
	metrics, err := system.GetSystemMetrics()
	if err != nil {
		logger.Error("Failed to get system metrics: %v", err)
	} else {
		if err := a.client.ReportMetrics(ctx, metrics); err != nil {
			logger.Error("Failed to report metrics: %v", err)
		} else {
			logger.Debug("Metrics reported.")
		}
	}

	// 2. Check for new policies/commands
	commands, err := a.client.FetchCommands(ctx)
	if err != nil {
		logger.Error("Failed to fetch commands: %v", err)
		return
	}

	for _, cmd := range commands {
		logger.Info("Executing command: %s", cmd.Name)
		output, cmdErr := system.ExecuteCommand(cmd.Command, cmd.Args) // A simplified example
		if cmdErr != nil {
			logger.Error("Command '%s' failed: %v, Output: %s", cmd.Name, cmdErr, output)
			a.client.ReportCommandResult(ctx, cmd.ID, "failed", output)
		} else {
			logger.Info("Command '%s' executed successfully. Output: %s", cmd.Name, output)
			a.client.ReportCommandResult(ctx, cmd.ID, "success", output)
		}
	}

	// 3. Apply policies (more complex, involves OS-specific actions)
	policies, err := a.client.FetchPolicies(ctx)
	if err != nil {
		logger.Error("Failed to fetch policies: %v", err)
	} else {
		for _, p := range policies {
			logger.Info("Applying policy: %s", p.Name)
			if err := system.ApplyPolicy(p); err != nil {
				logger.Error("Failed to apply policy %s: %v", p.Name, err)
				// Report policy application failure
			} else {
				logger.Debug("Policy %s applied.", p.Name)
				// Report policy application success
			}
		}
	}
}
