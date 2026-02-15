package service

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
)

// Get the icon names from the Font Awesome CSS file
func Dump(fontAwesomeStylesUri string) error {
	// 138Kb is expected so we do it this way
	resp, err := http.Get(fontAwesomeStylesUri)
	if err != nil {
		slog.Error("HTTP request failed", "error", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		slog.Error("HTTP request failed", "status", resp.StatusCode)
		return fmt.Errorf("HTTP request failed with status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read response body", "error", err)
		return err
	}
	re := regexp.MustCompile(`\.fa-([a-zA-Z0-9\-]+)\s*{\s*--fa:\s*"\\([a-fA-F0-9]+)"`)
	for _, match := range re.FindAllStringSubmatch(string(data), -1) {
		//Format as an 4-digit Go Unicode literal (\u0000)
		escaped := fmt.Sprintf("\\u%04s", match[2])
		fmt.Printf("%s: %s\n", match[1], escaped)
	}
	return nil
}
