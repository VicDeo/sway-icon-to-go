// main package for the sway-icon-to-go app

package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
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
	procPath             = "/proc"
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
	logger        *slog.Logger
	nameFormatter *NameFormatter
	iconProvider  *IconProvider
	config        *config.Config
	delim         string
	uniq          bool
	length        int
	configPath    string
}

// reloadConfig reloads the configuration from files
func (h *handler) reloadConfig() error {
	h.logger.Info("Reloading configuration...")

	cache.Clear()

	// Reload configuration
	newConfig, err := config.NewConfig(h.delim, h.uniq, h.length, h.configPath)
	if err != nil {
		return fmt.Errorf("failed to reload config: %w", err)
	}

	h.config = newConfig
	h.nameFormatter = NewNameFormatter(newConfig)

	h.logger.Info("Configuration reloaded successfully")
	return nil
}

// Window event handler
func (h handler) Window(ctx context.Context, event swayClient.WindowEvent) {
	for _, b := range WindowChangeTypes {
		if b == event.Change {
			if err := sway.ProcessWorkspaces(ctx, h.iconProvider, h.nameFormatter); err != nil {
				h.logger.Error("Error while processing the event", "error", err)
			}
		}
	}
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	uniq := flag.Bool("u", config.DefaultUniq, "display only unique icons. True by default")
	length := flag.Int("l", config.DefaultLength, "trim app names to this length. 12 by default")
	delim := flag.String("d", config.DefaultDelimiter, "app separator. \"|\" by default")
	configPath := flag.String("c", "", "path to the app-icons.yaml config file")
	flag.Parse()
	if flag.NArg() > 0 {
		if flag.Arg(0) == "awesome" {
			if err := service.FindFonts(); err != nil {
				logger.Error("Error while finding fonts", "error", err)
				os.Exit(1)
			}
			return
		} else if flag.Arg(0) == "help" {
			help()
			return
		} else if flag.Arg(0) == "parse" {
			if err := service.Dump(fontAwesomeStylesUri); err != nil {
				logger.Error("Error while parsing Font Awesome CSS file", "error", err)
				os.Exit(1)
			}
			return
		}
	}
	appConfig, configErr := config.NewConfig(*delim, *uniq, *length, *configPath)
	if configErr != nil {
		logger.Error("Error while getting config", "error", configErr)
		os.Exit(1)
	}
	nameFormatter := NewNameFormatter(appConfig)
	resolver := proc.LinuxResolver{ProcPath: procPath}
	processManager := proc.NewProcessManager(&resolver)
	iconProvider := NewIconProvider(processManager, appConfig)

	h := handler{
		EventHandler:  swayClient.NoOpEventHandler(),
		logger:        logger,
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
	logger.Info("Signal handler set up", "pid", os.Getpid())

	// go-sway event loop that listens for window events
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := swayClient.Subscribe(ctx, h, swayClient.EventTypeWindow)
		if err != nil {
			logger.Error("failed to connect to sway", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for events or signals
	for {
		select {
		case <-ctx.Done():
			return
		case sig := <-sigChan:
			logger.Info("Received signal", "signal", sig)
			if sig == syscall.SIGHUP {
				logger.Info("Reloading configuration...")
				if err := h.reloadConfig(); err != nil {
					logger.Warn("Failed to reload configuration", "error", err)
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
