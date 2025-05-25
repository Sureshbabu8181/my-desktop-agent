//go:build darwin
// +build darwin

package system

import (
	"fmt"
	"os/exec"
)

func createUserMacOS(username, password string) error {
	// Using dscl (Directory Service Command Line)
	cmd := exec.Command("sudo", "dscl", ".", "-create", "/Users/"+username)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create user entry: %w", err)
	}

	cmd = exec.Command("sudo", "dscl", ".", "-create", "/Users/"+username, "UserShell", "/bin/bash")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set user shell: %w", err)
	}

	// Set password
	cmd = exec.Command("sudo", "dscl", ".", "-passwd", "/Users/"+username, password)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set user password: %w", err)
	}
	return nil
}

// ApplyPolicy for macOS
func applyPolicyMacOS(p Policy) error {
	switch p.Type {
	case "filevault_encryption":
		// Requires `fdesetup` and potentially user interaction/passwords.
		// This is a complex area.
		return fmt.Errorf("filevault encryption not implemented for macOS")
	case "screen_lock_timeout":
		// Using `defaults write` to modify system preferences
		cmd := exec.Command("defaults", "write", "com.apple.screensaver", "idleTime", p.Value) // p.Value could be "300" for 5 minutes
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to set screen lock timeout on macOS: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("unknown policy type for macOS: %s", p.Type)
	}
}
