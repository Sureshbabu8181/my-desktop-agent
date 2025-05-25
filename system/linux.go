//go:build linux
// +build linux

package system

import (
	"fmt"
	"os/exec"
)

func createUserLinux(username, password string) error {
	cmd := exec.Command("useradd", "-m", username) // -m creates home directory
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create user %s: %w", username, err)
	}

	// Set password for the new user
	echoCmd := exec.Command("echo", fmt.Sprintf("%s:%s", username, password))
	passwdCmd := exec.Command("chpasswd")     // Reads from stdin
	passwdCmd.Stdin, _ = echoCmd.StdoutPipe() // Pipe echo output to chpasswd stdin

	if err := passwdCmd.Start(); err != nil {
		return fmt.Errorf("failed to start chpasswd: %w", err)
	}
	if err := echoCmd.Run(); err != nil {
		return fmt.Errorf("failed to run echo for password: %w", err)
	}
	if err := passwdCmd.Wait(); err != nil {
		return fmt.Errorf("failed to set password for user %s: %w", username, err)
	}
	return nil
}

// ApplyPolicy for Linux (highly dependent on policy type)
func applyPolicyLinux(p Policy) error {
	switch p.Type {
	case "password_complexity":
		// Modify /etc/pam.d/system-auth or /etc/login.defs
		return fmt.Errorf("password complexity policy on Linux not implemented")
	case "firewall_rule":
		// Use 'iptables' or 'firewall-cmd' (for firewalld)
		cmd := exec.Command("iptables", "-A", "INPUT", "-p", "tcp", "--dport", p.Value, "-j", "DROP")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to add iptables rule: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("unknown policy type for Linux: %s", p.Type)
	}
}
