BINARY_NAME=sway-icon-to-go
INSTALL_PATH=/usr/local/bin/$(BINARY_NAME)

.PHONY: build install reload clean

build:
	go build -o $(BINARY_NAME) ./cmd/main.go

install: build
	@echo "Installing to $(INSTALL_PATH)..."
	sudo cp $(BINARY_NAME) $(INSTALL_PATH)
	sudo chmod +x $(INSTALL_PATH)

# Reload the daemon without restarting all of Sway
reload: install
	@echo "Restarting $(BINARY_NAME)..."
	pkill $(BINARY_NAME) || true
	swaymsg exec $(INSTALL_PATH)

clean:
	rm -f $(BINARY_NAME)