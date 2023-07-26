
include scripts/make-rules/common.mk # must be the first to include
include scripts/make-rules/golang.mk

help:
	@echo "Usage: make <target>"
	@echo
	@echo " * 'lint' - Run golangci to lint source codes"
	@echo " * 'build' - Build harbor-cleaner with a container"
	@echo " * 'clean' - Clean artifacts"

.PHONY: lint 
lint: golangci-lint
	golangci-lint run -c $(ROOT_DIR)/.golangci.yml

.PHONY: test 
test:
	$(GO) test $(PKGS)

.PHONY: test 
build-local:
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(GO_BUILDFLAGS) -o $(OUTPUT_DIR)/harbor-cleaner ./cmd

.PHONY: build
build:
	$(CONTAINER_RUNTIME) run --rm                                   \
	  -v $(PWD):/opt/$(MODULE_NAME)                                 \
	  -w /opt/$(MODULE_NAME)                                        \
	  -e CGO_ENABLED=$(CGO_ENABLED)                                 \
	  $(GO_IMG)                                                     \
	      $(GO) build $(GO_BUILDFLAGS) 								\
		    -o $(OUTPUT_DIR)/harbor-cleaner $(CMD_DIR)/main.go;

image: build
	$(CONTAINER_RUNTIME) build -t zhangguanzhang/harbor-cleaner:$(VERSION) -f ./build/Dockerfile .

.PHONY: clean
clean:
	-rm -vrf ${OUTPUT_DIR}

