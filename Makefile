SHELL := /bin/bash

# Binary/output config
BIN_NAME       ?= cloudflare-ddns
OUT_DIR        ?= dist

# Build config
CGO_ENABLED    ?= 0
DEBUG          ?= 0

# Also support a positional "--debug" goal (use: `make build -- --debug`)
ifneq (,$(filter --debug,$(MAKECMDGOALS)))
DEBUG := 1
endif

# Production vs Debug flags
ifeq ($(DEBUG),1)
	# Debug: disable optimizations and inlining for easier debugging
	BUILD_FLAGS := -gcflags "all=-N -l"
else
	# Production: smaller binaries
	BUILD_FLAGS := -trimpath -ldflags "-s -w"
endif

.PHONY: build clean help targets --debug

.DEFAULT_GOAL := build

# Core build target
# Usage:
#   make                        # build for host (production)
#   make build                  # same as above
#   make build TARGET=linux-arm64
#   make build DEBUG=1          # debug build
build:
	@set -euo pipefail; \
	os=""; arch=""; \
	if [ -n "$(TARGET)" ]; then \
		os="$(TARGET)"; arch="$(TARGET)"; \
		os="$${os%%-*}"; arch="$${arch#*-}"; \
	else \
		os="$$(go env GOOS)"; arch="$$(go env GOARCH)"; \
	fi; \
	case "$$arch" in \
		x64|x86_64) arch=amd64 ;; \
		aarch64)   arch=arm64 ;; \
		*)         ;; \
	esac; \
	ext=""; if [ "$$os" = "windows" ]; then ext=".exe"; fi; \
	mkdir -p "$(OUT_DIR)"; \
	printf "Building %s for %s-%s (%s)\n" "$(BIN_NAME)" "$$os" "$$arch" "$$([ "$(DEBUG)" = "1" ] && echo debug || echo release)"; \
	GOOS="$$os" GOARCH="$$arch" CGO_ENABLED=$(CGO_ENABLED) \
	go build $(BUILD_FLAGS) -o "$(OUT_DIR)/$(BIN_NAME)-$$os-$$arch$$ext" .

clean:
	rm -rf "$(OUT_DIR)"

help:
	@echo "Build (default: production)"; \
	echo; \
	echo "Usage:"; \
	echo "  make [build]                     # host OS/arch"; \
	echo "  make build TARGET=linux-arm64     # cross-compile"; \
	echo "  make build DEBUG=1                # debug build"; \
	echo "  make build -- --debug             # debug build (alt syntax)"; \
	echo "  make targets                      # list valid TARGET values on this toolchain"; \
	echo; \
	echo "Notes:"; \
	echo "  TARGET format: <os>-<arch> (e.g., linux-amd64, linux-arm64, windows-x64)"; \
	echo "  Artifacts at $(OUT_DIR)/$(BIN_NAME)-<os>-<arch>[.exe]"

# No-op recipe for the positional flag
--debug:
	@true

# List all available TARGET values from this Go toolchain
targets:
	@echo "Available TARGET values (os-arch):"; \
	go tool dist list | awk -F/ '{print $$1 "-" $$2}' | sort -u

