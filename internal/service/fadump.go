package service

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
)

// Get the icon names from the Font Awesome CSS file
func Dump(fontAwesomeStylesUri string) error {
	// 138Kb is expected so we do it this way
	resp, err := http.Get(fontAwesomeStylesUri)
	if err != nil {
		slog.Warn("HTTP request failed", "error", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		slog.Warn("HTTP request failed", "status", resp.StatusCode)
		return fmt.Errorf("HTTP request failed with status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Warn("Failed to read response body", "error", err)
		return err
	}
	re := regexp.MustCompile(`\.fa-([a-zA-Z0-9\-]+)\s*{\s*--fa:\s*"(\\[a-fA-F0-9]+)"`)
	for _, match := range re.FindAllStringSubmatch(string(data), -1) {
		hexValue := strings.TrimPrefix(match[2], "\\")
		//Format as an 4-digit Go Unicode literal (\u0000)
		escaped := fmt.Sprintf("\\u%04s", hexValue)
		fmt.Printf("%s: %s\n", match[1], escaped)
	}
	return nil
}
