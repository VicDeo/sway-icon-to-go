// main package for the sway-icon-to-go app

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sway-icon-to-go/internal/config"
	"sway-icon-to-go/internal/sway"

	swayClient "github.com/joshuarubin/go-sway"
)

const (
	fontAwesomeStylesUri = "https://github.com/FortAwesome/Font-Awesome/raw/6.x/css/all.css"
	procPath             = "/proc"
	DefaultCacheSize     = 30
)

var (
	// WindowChangeTypes is a list of window event changes that we are interested in
	WindowChangeTypes = [...]swayClient.WindowEventChange{
		swayClient.WindowMove,
		swayClient.WindowNew,
		swayClient.WindowTitle,
		swayClient.WindowClose,
	}
)

// handler is a struct that handles the sway events
type handler struct {
	swayClient.EventHandler
	nameFormatter *ConfigNameFormatter
	iconProvider  *ConfigIconProvider
}

// ConfigIconProvider is a struct that provides the icon for the given pid and node name
type ConfigIconProvider struct {
	// Cache the executable name by pid
	pidCache map[uint32]string
	// Cache the icon by executable name or window name
	nameCache map[string]string
}

// ConfigNameFormatter is a struct that formats the workspace name according to the config
type ConfigNameFormatter struct{}

// NewConfigNameFormatter creates a new ConfigNameFormatter with the given config
func NewConfigNameFormatter(config *config.Config) *ConfigNameFormatter {
	return &ConfigNameFormatter{}
}

// Format the workspace name according to the config
func (c ConfigNameFormatter) Format(workspaceNumber int64, appIcons []string) string {
	return config.BuildName(workspaceNumber, appIcons)
}

// Get the icon for the given pid and node name
func (c ConfigIconProvider) GetIcon(pid *uint32, nodeName string) (string, bool) {
	var name string
	if pid != nil {
		if cachedName, ok := c.pidCache[*pid]; ok {
			name = cachedName
		} else {
			filename, err := getExecutableName(pid)
			if err != nil || filename == "" {
				log.Printf("Error while getting executable name: %v\n", err)
				name = nodeName
			} else {
				c.pidCache[*pid] = filename
				name = filename
			}
		}
	}
	if icon, ok := c.nameCache[name]; ok {
		return icon, true
	}
	icon, found := config.GetAppIcon(name)
	if found {
		c.nameCache[name] = icon
	}
	return icon, found
}

// Window event handler
func (h handler) Window(ctx context.Context, event swayClient.WindowEvent) {
	for _, b := range WindowChangeTypes {
		if b == event.Change {
			if err := sway.ProcessWorkspaces(ctx, h.iconProvider, h.nameFormatter); err != nil {
				log.Printf("Error while processing the event : %s\n", err)
			}
		}
	}
}

func main() {
	uniq := flag.Bool("u", config.DefaultUniq, "display only unique icons. True by default")
	length := flag.Int("l", config.DefaultLength, "trim app names to this length. 12 by default")
	delim := flag.String("d", config.DefaultDelimiter, "app separator. \"|\" by default")
	flag.Parse()
	if flag.NArg() == 0 {
	} else if flag.Arg(0) == "awesome" {
		if err := findFonts(); err != nil {
			log.Fatalf("Error while finding fonts: %v", err)
		}
		return
	} else if flag.Arg(0) == "help" {
		help()
		return
	} else if flag.Arg(0) == "parse" {
		if err := dump(); err != nil {
			log.Fatalf("Error while parsing Font Awesome CSS file: %v", err)
		}
		return
	}
	appConfig, configErr := config.GetConfig(*delim, *uniq, *length, "")
	if configErr != nil {
		log.Fatalf("Error while getting config: %v", configErr)
	}
	nameFormatter := NewConfigNameFormatter(appConfig)
	iconProvider := &ConfigIconProvider{
		pidCache:  make(map[uint32]string, DefaultCacheSize),
		nameCache: make(map[string]string, DefaultCacheSize),
	}
	h := handler{
		EventHandler:  swayClient.NoOpEventHandler(),
		nameFormatter: nameFormatter,
		iconProvider:  iconProvider,
	}

	// go-sway event loop that listens for window events
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := swayClient.Subscribe(ctx, h, swayClient.EventTypeWindow)
	if err != nil {
		log.Fatalf("failed to connect to sway: %v", err)
	}

	// Wait indefinitely
	select {}
}

// Print the help message
func help() {
	fmt.Println(`usage: sway-icon-to-go [-uc] [-l LENGTH] [-d DELIMITER] [help|awesome|parse]
  awesome    check if Font Awesome is available on your system (via fc-list)
  parse      parse Font Awesome CSS file to match icon names with their UTF-8 representation  
  help       print help
  -c         path to the app-icons.yaml config file
  -u         display only unique icons. True by default
  -l         trim app names to this length. 12 by default
  -d         app delimiter. "|" by default
	`)
}

// Check if Font Awesome is available on the system
func findFonts() error {
	cmd1 := exec.Command("fc-list")
	cmd2 := exec.Command("grep", "Awesome")
	cmd3 := exec.Command("sort")
	cmd2.Stdin, _ = cmd1.StdoutPipe()
	cmd3.Stdin, _ = cmd2.StdoutPipe()
	cmd3Output, _ := cmd3.StdoutPipe()

	if err := cmd3.Start(); err != nil {
		log.Printf("Error starting command: %v\n", err)
		return err
	}
	if err := cmd2.Start(); err != nil {
		log.Printf("Error starting command: %v\n", err)
		return err
	}
	if err := cmd1.Start(); err != nil {
		log.Printf("Error starting command: %v\n", err)
		return err
	}

	cmd3Result, err := io.ReadAll(cmd3Output)
	if err != nil {
		log.Printf("Error reading command output: %v\n", err)
		return err
	}

	// Wait for all commands to finish
	_ = cmd1.Wait()
	_ = cmd2.Wait()
	_ = cmd3.Wait()

	// Print the final result
	fmt.Printf("Result:\n%s\n", cmd3Result)
	return nil
}

// Get the icon names from the Font Awesome CSS file
func dump() error {
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
	fmt.Printf("Result:\n%s\n", string(data))
	// TODO: Fix this regex
	re := regexp.MustCompile(`\.fa-([^:]+):?:before[^"]+"(.*)"`)
	for _, match := range re.FindAllStringSubmatch(string(data), -1) {
		char := strings.Replace(match[2], "\\", "\\u", 1)
		fmt.Printf("%s: %s\n", match[1], char)
	}
	return nil
}

// Resolve the executable name for the given pid
func getExecutableName(pid *uint32) (string, error) {
	if pid == nil {
		return "", fmt.Errorf("pid is nil")
	}
	pidStr := strconv.FormatUint(uint64(*pid), 10)
	exePath := filepath.Join(procPath, pidStr, "exe")
	realPath, err := filepath.EvalSymlinks(exePath)
	if err != nil {
		return "", err
	}
	return filepath.Base(realPath), nil
}
