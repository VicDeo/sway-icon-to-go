// main package for the sway-icon-to-go app

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sway-icon-to-go/internal/cache"
	"sway-icon-to-go/internal/config"
	"sway-icon-to-go/internal/proc"
	"sway-icon-to-go/internal/service"
	"sway-icon-to-go/internal/sway"
	"syscall"

	swayClient "github.com/joshuarubin/go-sway"
)

const (
	fontAwesomeStylesUri = "https://github.com/FortAwesome/Font-Awesome/raw/6.x/css/all.css"
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
	config        *config.Config
	delim         string
	uniq          bool
	length        int
	configPath    string
}

// ConfigIconProvider is a struct that provides the icon for the given pid and node name
type ConfigIconProvider struct{}

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
	name, found := proc.GetProcessName(pid)
	if !found || name == "" {
		name = nodeName
	}

	icon, found := cache.GetIcon(name)
	if found {
		return icon, found
	}

	icon, found = config.GetAppIcon(name)
	if found {
		cache.SetIcon(name, icon)
	}
	return icon, found
}

// reloadConfig reloads the configuration from files
func (h *handler) reloadConfig() error {
	log.Println("Reloading configuration...")

	cache.Clear()

	// Reload configuration
	newConfig, err := config.GetConfig(h.delim, h.uniq, h.length, h.configPath)
	if err != nil {
		return fmt.Errorf("failed to reload config: %w", err)
	}

	h.config = newConfig
	h.nameFormatter = NewConfigNameFormatter(newConfig)

	log.Println("Configuration reloaded successfully")
	return nil
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
	configPath := flag.String("c", "", "path to the app-icons.yaml config file")
	flag.Parse()
	if flag.NArg() == 0 {
	} else if flag.Arg(0) == "awesome" {
		if err := service.FindFonts(); err != nil {
			log.Fatalf("Error while finding fonts: %v", err)
		}
		return
	} else if flag.Arg(0) == "help" {
		help()
		return
	} else if flag.Arg(0) == "parse" {
		if err := service.Dump(fontAwesomeStylesUri); err != nil {
			log.Fatalf("Error while parsing Font Awesome CSS file: %v", err)
		}
		return
	}
	appConfig, configErr := config.GetConfig(*delim, *uniq, *length, *configPath)
	if configErr != nil {
		log.Fatalf("Error while getting config: %v", configErr)
	}
	nameFormatter := NewConfigNameFormatter(appConfig)
	iconProvider := &ConfigIconProvider{}
	h := handler{
		EventHandler:  swayClient.NoOpEventHandler(),
		nameFormatter: nameFormatter,
		iconProvider:  iconProvider,
		config:        appConfig,
		delim:         *delim,
		uniq:          *uniq,
		length:        *length,
		configPath:    *configPath,
	}

	// Set up signal handling for SIGHUP (configuration reload)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP)
	log.Printf("Signal handler set up, PID: %d", os.Getpid())

	// go-sway event loop that listens for window events
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := swayClient.Subscribe(ctx, h, swayClient.EventTypeWindow)
		if err != nil {
			log.Fatalf("failed to connect to sway: %v", err)
		}
	}()

	// Wait for events or signals
	for {
		select {
		case <-ctx.Done():
			return
		case sig := <-sigChan:
			fmt.Printf("Received signal: %v", sig)
			if sig == syscall.SIGHUP {
				fmt.Println("Reloading configuration...")
				if err := h.reloadConfig(); err != nil {
					log.Printf("Failed to reload configuration: %v", err)
				}
			}
		}
	}
}

// Print the help message
func help() {
	fmt.Println(`usage: sway-icon-to-go [-u] [-l LENGTH] [-d DELIMITER] [-c CONFIG_PATH] [help|awesome|parse]
  awesome    check if Font Awesome is available on your system (via fc-list)
  parse      parse Font Awesome CSS file to match icon names with their UTF-8 representation  
  help       print help
  -c         path to the app-icons.yaml config file
  -u         display only unique icons. True by default
  -l         trim app names to this length. 12 by default
  -d         app delimiter. "|" by default

Configuration can be reloaded at runtime by sending SIGHUP signal:
  kill -HUP <pid>
	`)
}
