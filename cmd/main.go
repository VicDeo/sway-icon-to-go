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
	// windowChangeTypes is a map of window event changes that we are interested in.
	windowChangeTypes = map[swayClient.WindowEventChange]bool{
		swayClient.WindowMove:  true,
		swayClient.WindowNew:   true,
		swayClient.WindowTitle: true,
		swayClient.WindowClose: true,
	}
)

// handler is a struct that handles the sway events
type handler struct {
	swayClient.EventHandler
	nameFormatter *NameFormatter
	iconProvider  *IconProvider
	config        *config.Config
	format        *config.Format
	configPath    string
}

// reloadConfig reloads the configuration from files
func (h *handler) reloadConfig() error {
	slog.Info("Reloading configuration...")

	// Reload configuration
	newConfig, err := config.NewConfig(h.configPath, h.format)
	if err != nil {
		return fmt.Errorf("failed to reload config: %w", err)
	}

	h.config = newConfig
	h.nameFormatter = NewNameFormatter(h.format)
	h.iconProvider.ClearCache()
	slog.Info("Configuration reloaded successfully")
	return nil
}

// Window event handler
func (h handler) Window(ctx context.Context, event swayClient.WindowEvent) {
	if _, ok := windowChangeTypes[event.Change]; !ok {
		return
	}
	if err := sway.ProcessWorkspaces(ctx, h.iconProvider, h.nameFormatter); err != nil {
		slog.Error("Error while processing the event", "error", err)
	}
}

func main() {
	// Set up the logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	slog.SetDefault(logger)

	// Set up the flags
	format := config.DefaultFormat()
	flag.BoolVar(&format.Uniq, "u", format.Uniq, "display only unique icons. True by default")
	flag.IntVar(&format.Length, "l", format.Length, "trim app names to this length. 12 by default")
	flag.StringVar(&format.Delimiter, "d", format.Delimiter, "app separator. \"|\" by default")

	configPath := flag.String("c", "", "path to the app-icons.yaml config file")
	flag.Parse()

	// Validate the arguments
	if flag.NArg() > 0 {
		if flag.Arg(0) == "awesome" {
			fonts, err := service.FindFonts()
			if err != nil {
				slog.Error("Error while finding fonts", "error", err)
				os.Exit(1)
			}
			fmt.Println(fonts)
			return
		} else if flag.Arg(0) == "help" {
			help()
			return
		} else if flag.Arg(0) == "parse" {
			if err := service.Dump(fontAwesomeStylesUri); err != nil {
				slog.Error("Error while parsing Font Awesome CSS file", "error", err)
				os.Exit(1)
			}
			return
		}
	}
	// Get the configuration
	appConfig, configErr := config.NewConfig(*configPath, format)
	if configErr != nil {
		slog.Error("Error while getting config", "error", configErr)
		os.Exit(1)
	}
	// Run the application
	run(appConfig, format, configPath)
}

// run runs the application.
func run(appConfig *config.Config, format *config.Format, configPath *string) {
	nameFormatter := NewNameFormatter(format)

	// Set up the pid to name resolver
	resolver := proc.LinuxResolver{ProcPath: procPath}
	processManager := proc.NewProcessManager(&resolver)

	// Set up the icon provider
	iconCache := cache.NewCache()
	iconProvider := NewIconProvider(processManager, appConfig, iconCache)

	h := handler{
		EventHandler:  swayClient.NoOpEventHandler(),
		nameFormatter: nameFormatter,
		iconProvider:  iconProvider,
		config:        appConfig,
		format:        format,
		configPath:    *configPath,
	}

	// Set up signal handling for SIGHUP (configuration reload)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP)
	slog.Info("Signal handler set up", "pid", os.Getpid())

	// go-sway event loop that listens for window events
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(cancel context.CancelFunc) {
		err := swayClient.Subscribe(ctx, h, swayClient.EventTypeWindow)
		if err != nil {
			slog.Error("failed to connect to sway", "error", err)
			cancel()
		}
	}(cancel)

	// Wait for events or signals
	for {
		select {
		case <-ctx.Done():
			return
		case sig := <-sigChan:
			slog.Info("Received signal", "signal", sig)
			if sig == syscall.SIGHUP {
				if err := h.reloadConfig(); err != nil {
					slog.Warn("Failed to reload configuration", "error", err)
				}
			}
		}
	}
}

// help prints the help message.
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
