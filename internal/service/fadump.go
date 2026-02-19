// Package service provides a way to dump the Font Awesome icons
// and check if Font Awesome is available on the system.

package service

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	timeout = 10 * time.Second
)

// Get the icon names from the Font Awesome CSS file.
func Dump(fontAwesomeCSSURL string) (string, error) {
	data, err := downloadFontAwesomeStyles(fontAwesomeCSSURL)
	if err != nil {
		slog.Error("Failed to download Font Awesome styles", "error", err)
		return "", err
	}
	return parseFontAwesomeStyles(data)
}

func parseFontAwesomeStyles(data []byte) (string, error) {
	parsed := strings.Builder{}
	re := regexp.MustCompile(`\.fa-([a-zA-Z0-9\-]+)\s*{\s*--fa:\s*"\\([a-fA-F0-9]+)"`)
	for _, match := range re.FindAllStringSubmatch(string(data), -1) {
		escaped := formatUnicodeLiteral(match[2])
		parsed.WriteString(fmt.Sprintf("%s: %s\n", match[1], escaped))
	}
	return parsed.String(), nil
}

// formatUnicodeLiteral formats the given unicode string as a 4-digit Go Unicode literal (\u0000)
func formatUnicodeLiteral(unicode string) string {
	return fmt.Sprintf("\\u%04s", unicode)
}

func downloadFontAwesomeStyles(fontAwesomeCSSURL string) ([]byte, error) {
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(fontAwesomeCSSURL)
	if err != nil {
		slog.Error("HTTP request failed", "error", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		slog.Error("HTTP request failed", "status", resp.StatusCode)
		return nil, fmt.Errorf("HTTP request failed with status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read response body", "error", err)
		return nil, err
	}
	return data, nil
}
