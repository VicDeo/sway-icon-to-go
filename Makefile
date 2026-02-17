BINARY_NAME=sway-icon-to-go
INSTALL_PATH=/usr/local/bin/$(BINARY_NAME)
USER_BIN_DIR=$(HOME)/.local/bin
USER_SERVICE_DIR=$(HOME)/.config/systemd/user
SERVICE_FILE=configs/sway-icon-to-go.service

.PHONY: build install install-service uninstall-service reload-service reload clean

build:
	go build -o $(BINARY_NAME) ./cmd

install: build
	@echo "Installing to $(INSTALL_PATH)..."
	sudo cp $(BINARY_NAME) $(INSTALL_PATH)
	sudo chmod +x $(INSTALL_PATH)

# Reload the daemon without restarting all of Sway (for non-service install)
reload: install
	@echo "Restarting $(BINARY_NAME)..."
	pkill $(BINARY_NAME) || true
	swaymsg exec $(INSTALL_PATH)

# Install as user-level systemd service
install-service: build
	@echo "Installing $(BINARY_NAME) as user service..."
	@mkdir -p $(USER_BIN_DIR) $(USER_SERVICE_DIR)
	cp $(BINARY_NAME) $(USER_BIN_DIR)/
	cp $(SERVICE_FILE) $(USER_SERVICE_DIR)/
	systemctl --user daemon-reload
	systemctl --user enable sway-icon-to-go.service
	@echo "Starting service..."
	systemctl --user start sway-icon-to-go.service
	@echo "Done. Service installed and running."

# Uninstall user-level systemd service
uninstall-service:
	@echo "Uninstalling $(BINARY_NAME) service..."
	systemctl --user stop sway-icon-to-go.service || true
	systemctl --user disable sway-icon-to-go.service || true
	rm -f $(USER_BIN_DIR)/$(BINARY_NAME)
	rm -f $(USER_SERVICE_DIR)/sway-icon-to-go.service
	systemctl --user daemon-reload
	@echo "Done. Service uninstalled."

# Rebuild, reinstall binary, and restart the user service
reload-service: build
	@echo "Reloading $(BINARY_NAME) service..."
	cp $(BINARY_NAME) $(USER_BIN_DIR)/
	systemctl --user restart sway-icon-to-go.service
	@echo "Done. Service restarted with new binary."


clean:
	rm -f $(BINARY_NAME)