// Main package for the sway-icon-to-go app.

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
	"sway-icon-to-go/internal/display"
	"sway-icon-to-go/internal/proc"
	"sway-icon-to-go/internal/service"
	"sway-icon-to-go/internal/sway"
	"syscall"
)

const (
	fontAwesomeCSSURL = "https://github.com/FortAwesome/Font-Awesome/raw/6.x/css/all.css"
	procPath          = "/proc"
)

func main() {
	var verbose bool

	// Until we have a real log level let our logger be non-verbose
	setupLogger(false)

	// Set up the flags
	format := config.DefaultFormat()
	flag.BoolVar(&format.Uniq, "u", format.Uniq, "display only unique icons (default true)")
	flag.IntVar(&format.Length, "l", format.Length, "trim app names to this length, -1 = no trim (default 12)")
	flag.StringVar(&format.Delimiter, "d", format.Delimiter, "app separator (default \"|\")")
	flag.BoolVar(&verbose, "v", false, "enable verbose/debug logging")

	// Set up the config path
	configPath := flag.String("c", "", "path to app-icons.yaml (auto-detect from ~/.config/sway or ~/.config/i3 if empty)")
	flag.Usage = help
	flag.Parse()

	// Adjust the log level according to the verbose flag
	setupLogger(verbose)

	if format.Length < -1 {
		slog.Error("Length can not be less than -1")
		os.Exit(1)
	}

	// Validate the arguments
	if flag.NArg() > 0 {
		switch flag.Arg(0) {
		case "awesome":
			fonts, err := service.FindFonts()
			if err != nil {
				slog.Error("Error while finding fonts", "error", err)
				os.Exit(1)
			}
			fmt.Println(fonts)
			return
		case "help":
			help()
			return
		case "parse":
			dump, err := service.Dump(fontAwesomeCSSURL)
			if err != nil {
				slog.Error("Error while parsing Font Awesome CSS file", "error", err)
				os.Exit(1)
			}
			fmt.Println(dump)
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
	run(appConfig, configPath)
}

func setupLogger(verbose bool) {
	logLevel := new(slog.LevelVar)
	if verbose {
		logLevel.Set(slog.LevelDebug)
	}
	// Set up the logger
	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel, AddSource: verbose})
	logger := slog.New(h)
	slog.SetDefault(logger)
}

// run runs the application.
func run(appConfig *config.Config, configPath *string) {
	nameFormatter := display.NewNameFormatter(appConfig.Format)

	// Set up the pid to name resolver
	resolver := proc.LinuxResolver{ProcPath: procPath}
	processManager := proc.NewProcessManager(&resolver)

	// Set up the icon provider
	iconCache := cache.NewCache()
	iconProvider := display.NewIconProvider(processManager, display.AppToIconMap(appConfig.AppToIcon), iconCache)

	// Set up signal handling for SIGHUP (configuration reload)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP)
	slog.Info("Signal handler set up", "pid", os.Getpid())

	// go-sway event loop that listens for window events
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	h := sway.NewHandler(nameFormatter, iconProvider, appConfig)

	go func(cancel context.CancelFunc) {
		err := sway.Subscribe(ctx, h)
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
				newConfig, err := config.NewConfig(*configPath, appConfig.Format)
				if err != nil {
					slog.Error("Failed to reload configuration", "error", err)
					continue
				}

				if err := h.ReloadConfig(newConfig); err != nil {
					slog.Error("Failed to reload configuration", "error", err)
				}
			}
		}
	}
}

// help prints the help message.
func help() {
	fmt.Fprintf(os.Stderr, `Renames sway workspaces by window names with Font Awesome icons.

Usage:
  sway-icon-to-go [options] [help|awesome|parse]

With no command, runs the workspace daemon.

Commands:
  awesome    list Font Awesome fonts installed on the system (empty output means not installed)
  parse      dump icon name → UTF-8 mapping (pipe to fa-icons.yaml)
  help       show this help

Flags:
  -c         path to app-icons.yaml (auto-detect from ~/.config/sway or ~/.config/i3 if empty)
  -u         display only unique icons (default true)
  -l         trim app names to this length, -1 = no trim (default 12)
  -d         app separator (default "|")
  -v         enable verbose/debug logging

Configuration can be reloaded at runtime by sending SIGHUP signal:
  pkill -HUP sway-icon-to-go
`)
}
