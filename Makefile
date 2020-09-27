PROJECTNAME := $(shell basename "$(PWD)")

# Go related variables.
PROJECTROOT := $(shell pwd)
GOBIN := $(PROJECTROOT)/bin

# Shell script related variables.
UTILDIR := $(PROJECTROOT)/make-scripts/utils
SPINNER := $(UTILDIR)/spinner.sh
BUILDIR := $(PROJECTROOT)/make-scripts/build

UNITTESTCOMMAND := $(shell cd cli/cmd/; go test -v)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: default
default: help

## Build the cli binary
build-cli:
	@printf "🔨 Building binary $(GOBIN)/$(PROJECTNAME)\n" 
	@./make-scripts/build/build-cli.sh
	@cp ./cli/keptn $(GOBIN)/
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
	@$(UNITTESTCOMMAND)
	@printf "👍 Done\n"

## Prepare code for PR
prepare-for-pr: fmt lint test-unit-cli
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
	@printf "Keptn"
	@printf "Help is coming soon..."
	@printf ""
