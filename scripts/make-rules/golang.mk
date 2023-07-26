MODULE_NAME = $(shell awk '$$1=="module"{print $$2}' go.mod)

GO           ?= go
GOFMT        ?= $(GO)fmt
GO_VERSION        ?= $(shell $(GO) version)
GO_VERSION_NUMBER ?= $(word 3, $(GO_VERSION))
CGO_ENABLED    ?= 0

GO_BUILDFLAGS  ?= -gcflags="all=-trimpath=$(PWD) $(shell if [ "$(DEBUG)" = 1 ]; then echo -n -N -l; fi)" -asmflags "all=-trimpath=$(PWD)"
EXTRA_LDFLAGS  ?=
BUILD_LDFLAGS  := -ldflags "-X $(MODULE_NAME)/pkg/version.Version=$(VERSION)        \
	        -X $(MODULE_NAME)/pkg/version.gitCommit=$(COMMIT)                       \
			-X $(MODULE_NAME)/pkg/version.buildDate=$(BUILD_DATE)                   \
	        $(EXTRA_LDFLAGS)"
GO_BUILDFLAGS += $(BUILD_LDFLAGS)

GO_IMG         ?= golang:1.19-alpine

.PHONY: golangci-lint
golangci-lint:
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
