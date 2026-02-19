package sway

import (
	"context"
	"fmt"
	"log/slog"
	"sway-icon-to-go/internal/config"
	"sway-icon-to-go/internal/display"

	sc "github.com/joshuarubin/go-sway"
)

var (
	// windowChangeTypes is a map of window event changes that we are interested in.
	windowChangeTypes = map[sc.WindowEventChange]bool{
		sc.WindowMove:  true,
		sc.WindowNew:   true,
		sc.WindowTitle: true,
		sc.WindowClose: true,
	}
)

// handler is a struct that handles the sway events
type handler struct {
	sc.EventHandler
	nameFormatter NameFormatter
	iconProvider  *display.IconProvider
	config        *config.Config
	format        *config.Format
	configPath    string
}

// NewHandler creates a new handler instance.
func NewHandler(nameFormatter NameFormatter, iconProvider *display.IconProvider, config *config.Config, format *config.Format, configPath string) *handler {
	h := &handler{
		EventHandler:  sc.NoOpEventHandler(),
		nameFormatter: nameFormatter,
		iconProvider:  iconProvider,
		config:        config,
		format:        format,
		configPath:    configPath,
	}
	return h
}

// ReloadConfig reloads the configuration from files
func (h *handler) ReloadConfig() error {
	slog.Info("Reloading configuration...")

	// Reload configuration
	newConfig, err := config.NewConfig(h.configPath, h.format)
	if err != nil {
		return fmt.Errorf("failed to reload config: %w", err)
	}

	h.nameFormatter = display.NewNameFormatter(h.format.Delimiter, h.format.Length, h.format.Uniq)
	h.config = newConfig
	h.iconProvider.ClearCache()
	h.iconProvider.SetIconMap(display.AppToIconMap(newConfig.AppToIcon))
	slog.Info("Configuration reloaded successfully")
	return nil
}

// Window event handler
func (h handler) Window(ctx context.Context, event sc.WindowEvent) {
	if _, ok := windowChangeTypes[event.Change]; !ok {
		return
	}
	if err := h.processWorkspaces(ctx); err != nil {
		slog.Error("Error while processing the event", "error", err)
	}
}

// processWorkspaces processes the workspaces and renames them according to the name formatter and icon provider
// basing on the apps running on the workspaces.
func (h *handler) processWorkspaces(ctx context.Context) error {
	sway, err := NewSwayClient(ctx)
	if err != nil {
		return err
	}

	// Then traverse the tree and populate the workspaces map.
	workspaces, err := sway.CollectWorkspaces()
	if err != nil {
		return err
	}

	// Add icons to the all windows of all workspaces.
	if err := h.iconProvider.AddIcons(workspaces); err != nil {
		return err
	}

	// Send all commands at once as there could be a mess otherwise.
	if err := sway.RenameWorkspaces(workspaces, h.nameFormatter); err != nil {
		return err
	}
	return nil
}
