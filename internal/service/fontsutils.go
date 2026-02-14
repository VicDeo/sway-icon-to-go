package service

import (
	"io"
	"log/slog"
	"os/exec"
	"strings"
)

// Check if Font Awesome is available on the system
func FindFonts() (string, error) {
	cmd1 := exec.Command("fc-list")
	cmd2 := exec.Command("grep", "Awesome")
	cmd3 := exec.Command("sort")

	cmd1Stdout, err := cmd1.StdoutPipe()
	if err != nil {
		slog.Warn("Error creating stdout pipe", "error", err)
		return "", err
	}

	cmd2.Stdin = cmd1Stdout
	cmd2Stdout, err := cmd2.StdoutPipe()
	if err != nil {
		slog.Warn("Error creating stdout pipe", "error", err)
		return "", err
	}

	cmd3.Stdin = cmd2Stdout
	cmd3Output, err := cmd3.StdoutPipe()
	if err != nil {
		slog.Warn("Error creating stdout pipe", "error", err)
		return "", err
	}

	if err := cmd3.Start(); err != nil {
		slog.Warn("Error starting command", "error", err)
		return "", err
	}
	if err := cmd2.Start(); err != nil {
		slog.Warn("Error starting command", "error", err)
		return "", err
	}
	if err := cmd1.Start(); err != nil {
		slog.Warn("Error starting command", "error", err)
		return "", err
	}

	cmd3Result, err := io.ReadAll(cmd3Output)
	if err != nil {
		slog.Warn("Error reading command output", "error", err)
		return "", err
	}

	// Wait for all commands to finish
	err = cmd1.Wait()
	if err != nil {
		slog.Warn("Error waiting for command", "error", err)
		return "", err
	}
	err = cmd2.Wait()
	if err != nil {
		slog.Warn("Error waiting for command", "error", err)
		return "", err
	}
	err = cmd3.Wait()
	if err != nil {
		slog.Warn("Error waiting for command", "error", err)
		return "", err
	}

	return strings.TrimSpace(string(cmd3Result)), nil
}
