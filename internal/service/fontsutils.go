package service

import (
	"log/slog"
	"os/exec"
	"strings"
)

// Check if Font Awesome is available on the system
func FindFonts() (string, error) {
	fonts := make([]string, 0)

	cmd := exec.Command("fc-list")
	cmdOutput, err := cmd.Output()
	if err != nil {
		slog.Error("Error executing fc-list", "error", err)
		return "", err
	}

	for _, font := range strings.Split(strings.TrimSpace(string(cmdOutput)), "\n") {
		if strings.Contains(font, "Awesome") {
			fonts = append(fonts, font)
		}
	}
	return strings.Join(fonts, "\n"), nil
}
