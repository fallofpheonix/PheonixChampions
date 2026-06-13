# Phoenix Game Layer - Development Orchestration
# Use this Makefile to run the Python core and Godot client.

.PHONY: run-core run-client run-all

# Paths
GODOT_BIN = /Applications/Godot.app/Contents/MacOS/Godot
CORE_SCRIPT = ./core/main.py
CLIENT_DIR = ./client

# Run the Python core directly
run-core:
	python3 $(CORE_SCRIPT)

# Run the Godot client
run-client:
	$(GODOT_BIN) --path $(CLIENT_DIR)

# Run both (Godot will launch the Python core)
run-all:
	$(GODOT_BIN) --path $(CLIENT_DIR)
