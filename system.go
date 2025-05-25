package system

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// DeviceInfo and SystemMetrics structs from api/client.go would also be defined here.
// For simplicity, let's reuse them.

// GetDeviceInfo gathers basic information about the device
func GetDeviceInfo() (*DeviceInfo, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %w", err)
	}
	return &DeviceInfo{
		Hostname:     hostname,
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
	}, nil
}

// GetSystemMetrics gathers CPU and memory usage (simplified example)
func GetSystemMetrics() (*SystemMetrics, error) {
	// This is highly OS-dependent and complex.
	// On Linux, you'd parse /proc/meminfo, /proc/stat.
	// On Windows, use WMI or performance counters.
	// On macOS, use sysctl or equivalent.
	// For a real agent, you'd use external libraries or detailed system calls.
	return &SystemMetrics{
		CPUUsage:    0.5, // Placeholder
		MemoryUsage: 0.7, // Placeholder
	}, nil
}

// ExecuteCommand runs a shell command
func ExecuteCommand(cmd string, args []string) (string, error) {
	command := exec.Command(cmd, args...)
	output, err := command.CombinedOutput()
	return string(output), err
}

// ApplyPolicy (simplified, requires specific implementations per OS)
func ApplyPolicy(p Policy) error {
	switch runtime.GOOS {
	case "windows":
		return applyPolicyWindows(p)
	case "darwin":
		return applyPolicyMacOS(p)
	case "linux":
		return applyPolicyLinux(p)
	default:
		return fmt.Errorf("unsupported OS for policy application: %s", runtime.GOOS)
	}
}
