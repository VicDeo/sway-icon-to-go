package service

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// Get the icon names from the Font Awesome CSS file
func Dump(fontAwesomeStylesUri string) error {
	// 138Kb is expected so we do it this way
	resp, err := http.Get(fontAwesomeStylesUri)
	if err != nil {
		log.Printf("HTTP request failed: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v\n", err)
		return err
	}
	re := regexp.MustCompile(`\.fa-([a-zA-Z0-9\-]+)\s*{\s*--fa:\s*"(\\[a-fA-F0-9]+)"`)
	for _, match := range re.FindAllStringSubmatch(string(data), -1) {
		char := strings.Replace(match[2], "\\", "\\u", 1)
		fmt.Printf("%s: %s\n", match[1], char)
	}
	return nil
}
