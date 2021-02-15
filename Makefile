PROJECTNAME := $(shell basename "$(PWD)")

# Go related variables.
PROJECTROOT := $(shell pwd)
GOBIN := $(PROJECTROOT)/bin

# Shell script related variables.
UTILDIR := $(PROJECTROOT)/make-scripts/utils
SPINNER := $(UTILDIR)/spinner.sh
BUILDIR := $(PROJECTROOT)/make-scripts/build

CREATEBIN := $(shell [ ! -d ./bin ] && mkdir bin)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

# output filename for cli binary
OUTPUT_EXECUTEABLE_NAME := $(PROJECTNAME)

.PHONY: default
default: help

## Build the cli binary
build-cli:
	@printf "🔨 Building binary '$(OUTPUT_EXECUTEABLE_NAME)'\n"
	@./make-scripts/build/build-cli.sh
	@cp cli/$(OUTPUT_EXECUTEABLE_NAME) $(GOBIN)/
	@printf "👍 Done\n"

## Build all docker images
build-docker:
	@printf "🔨 Building docker images\n"
	@./make-scripts/build/build-docker.sh
	@printf "👍 Done\n"

## Start the bridge
start-bridge:
	@printf "🚀 Starting Bridge\n" 
	@./make-scripts/start_bridge.sh

## Install helm
install-helm:
	@printf "🔨 Installing Helm\n" 
	@./make-scripts/install_helm.sh
	@printf "👍 Done\n"

## Lint the code
install-golint:
	@printf "🔨 Installing golint\n" 
	@./make-scripts/install_golint.sh
	@printf "👍 Done\n"

## Format the code
fmt:
	@printf "🔨 Formatting\n" 
	@gofmt -l -s .
	@printf "👍 Done\n"

## Check codebase for style mistakes
lint: install-golint
	@printf "🔨 Linting\n"
	@golint ./...
	@printf "👍 Done\n"

## Clean build files
clean:
	@printf "🔨 Cleaning build cache\n" 
	@go clean .
	@printf "👍 Done\n"
	@-rm $(GOBIN)/* 2>/dev/null

## Run unit tests on the CLI
test-unit-cli:
	@printf "⚙️ Running unit tests on the CLI\n" 
	@./make-scripts/run_cli_test.sh
	@printf "👍 Done\n"

## Prepare code for PR
prepare-for-pr: fmt lint
	@printf "❗️ Remember to run the tests"
	@git diff-index --quiet HEAD -- ||\
	(echo "-----------------" &&\
	echo "NOTICE: There are some files that have not been committed." &&\
	echo "-----------------\n" &&\
	git status &&\
	echo "\n-----------------" &&\
	echo "NOTICE: There are some files that have not been committed." &&\
	echo "-----------------\n"  &&\
	exit 0)

help:
	@echo "KEPTN"
	@echo ""
	@echo "* build-cli: Build the keptn cli and save it in bin/"
	@echo "* build-docker: Build the keptn docker images"
	@echo "* start-bridge: Start the bridge server"
	@echo "* install-helm: Install the helm binary in your local"
	@echo "* install-golint: Install golint for linting the code"
	@echo "* fmt: Formats the codebase"
	@echo "* lint: Lints the codebase"
	@echo "* clean: Cleans the build cache"
	@echo "* test-unit-cli: Run unit tests on the Keptn CLI"
	@echo "* prepare-for-pr: Makes the code ready for PR by formatting, linting and checking for uncommitted files"
	@echo ""
	@echo "Please visit https://keptn.sh for more information."
	@echo "Get in touch with us via Slack: https://slack.keptn.sh"
