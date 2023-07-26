# Set shell to bash
SHELL := /bin/bash
# Current version of the project.x`
TAG_VERSION ?= $(shell git describe --exact-match --tags --abbrev=0  2> /dev/null || echo untagged)
VERSION := $(TAG_VERSION)$(shell if ! git diff-index --quiet HEAD; then echo -n "-dirty"; fi)
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

# include the common makefile
COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
ifeq ($(origin ROOT_DIR),undefined)
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/../.. && pwd -P))
endif

CONTAINER_RUNTIME ?= docker

# A list of all packages.
PKGS := $(shell go list ./... | grep -v /vendor | grep -v /test)


# Project main package location (can be multiple ones).
CMD_DIR := cmd

# Project output directory.
OUTPUT_DIR := ./bin

# Golang standard bin directory.
BIN_DIR := $(GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter

UNAME := $(shell uname)

