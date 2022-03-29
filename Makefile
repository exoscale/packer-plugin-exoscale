## Project

PACKAGE := github.com/exoscale/packer-plugin-exoscale
PROJECT_URL := https://$(PACKAGE)
GO_BIN_OUTPUT_NAME := packer-plugin-exoscale

API_VERSION := $(shell go run . describe | jq -r '.api_version')
EXTRA_ARGS := -parallel 3 -count=1 -failfast

# Dependencies

# Requires: https://github.com/exoscale/go.mk
# - install: git submodule update --init --recursive go.mk
# - update:  git submodule update --remote
include go.mk/init.mk
include go.mk/public.mk

# Packer SDK
# REF: https://github.com/hashicorp/packer-plugin-sdk

PACKER_SDK_VERSION := v0.2.11

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

# Release

#.PHONY: release
#release:
#	see release-default in go.mk/release.mk

# Clean

clean::
	rm -f '$(GO_BIN_OUTPUT_NAME)'
