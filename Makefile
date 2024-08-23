GO_MK_REF := v2.0.3

# make go.mk a dependency for all targets
.EXTRA_PREREQS = go.mk

ifndef MAKE_RESTARTS
# This section will be processed the first time that make reads this file.

# This causes make to re-read the Makefile and all included
# makefiles after go.mk has been cloned.
Makefile:
	@touch Makefile
endif

.PHONY: go.mk
.ONESHELL:
go.mk:
	@if [ ! -d "go.mk" ]; then
		git clone https://github.com/exoscale/go.mk.git
	fi
	@cd go.mk
	@if ! git show-ref --quiet --verify "refs/heads/${GO_MK_REF}"; then
		git fetch
	fi
	@if ! git show-ref --quiet --verify "refs/tags/${GO_MK_REF}"; then
		git fetch --tags
	fi
	git checkout --quiet ${GO_MK_REF}

## Project

PACKAGE := github.com/exoscale/packer-plugin-exoscale
PROJECT_URL := https://$(PACKAGE)
GO_BIN_OUTPUT_NAME := packer-plugin-exoscale

API_VERSION := $(shell go run . describe | jq -r '.api_version')
EXTRA_ARGS := -parallel 3 -count=1 -failfast

PACKER_PLUGINS_DIR := $(HOME)/.packer.d/plugins

# Dependencies

# Requires: https://github.com/exoscale/go.mk
go.mk/init.mk:
include go.mk/init.mk
go.mk/public.mk:
include go.mk/public.mk

# GoLang

GO_VERSION := $(shell go version | sed -nE 's|^.*\s+go([0-9]+\.[0-9]+)[^0-9].*$$|\1|p')
GO_MOD_VERSION := $(shell sed -nE 's|^go\s+([0-9]+\.[0-9]+)$$|\1|p' go.mod)
ifneq ($(GO_VERSION), $(GO_MOD_VERSION))
$(warning GoLang versions mismatch (Toolchain: $(GO_VERSION); go.mod: $(GO_MOD_VERSION)))
endif

# Packer SDK
# REF: https://github.com/hashicorp/packer-plugin-sdk

PACKER_SDK_VERSION := v0.4.0

PACKER_SDK_MOD_VERSION := $(shell sed -nE 's|^\s*github.com/hashicorp/packer-plugin-sdk\s+(v[.0-9]+)$$|\1|p' go.mod)
ifneq ($(PACKER_SDK_VERSION), $(PACKER_SDK_MOD_VERSION))
$(warning Packer SDK versions mismatch (Makefile: $(PACKER_SDK_VERSION); go.mod: $(PACKER_SDK_MOD_VERSION)))
endif


## Targets

# Dependencies

.PHONY: install-packer-sdc
install-packer-sdc:  ## Packer Software Development Command (SDC)
	'$(GO)' install github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc@$(PACKER_SDK_VERSION)

# Artefacts

.PHONY: generate
generate: install-packer-sdc
	'$(GO)' generate ./...

# Tests

.PHONY: test-acc test-verbose test
test: GO_TEST_EXTRA_ARGS=${EXTRA_ARGS}
test-verbose: GO_TEST_EXTRA_ARGS+=$(EXTRA_ARGS)
test-acc: GO_TEST_EXTRA_ARGS=-v $(EXTRA_ARGS)
test-acc: ## Runs acceptance tests (requires valid Exoscale API credentials)
	PACKER_ACC=1 '$(GO)' test \
	  -race \
	  -timeout 60m \
	  --tags=testacc \
	  $(GO_TEST_EXTRA_ARGS) \
	  $(GO_TEST_PKGS)

# Install (locally)

$(PACKER_PLUGINS_DIR):
	mkdir -p '$(PACKER_PLUGINS_DIR)'

.PHONY: install
install $(PACKER_PLUGINS_DIR)/$(GO_BIN_OUTPUT_NAME): $(GO_BIN_OUTPUT_DIR)/$(GO_BIN_OUTPUT_NAME) $(PACKER_PLUGINS_DIR)
	cp -v '$(GO_BIN_OUTPUT_DIR)/$(GO_BIN_OUTPUT_NAME)' '$(PACKER_PLUGINS_DIR)/$(GO_BIN_OUTPUT_NAME)'

.PHONY: uninstall
uninstall:
	rm -fv "$${HOME}/.packer.d/plugins/$(GO_BIN_OUTPUT_NAME)"

# Release

.PHONY: release
release: release-precheck release-notes
	API_VERSION='$(API_VERSION)' '$(GORELEASER)' release $(GORELEASER_OPTS)

# Clean

clean::
	rm -f '$(GO_BIN_OUTPUT_NAME)'
