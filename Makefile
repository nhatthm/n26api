JSON_CLI_VERSION = 1.8.3
SWAC_VERSION = 0.1.19
GOLANGCI_LINT_VERSION ?= v1.48.0

BIN_DIR = bin
VENDOR_DIR = vendor
OPENAPI = openapi.yaml

GO ?= go
GOLANGCI_LINT ?= $(BIN_DIR)/golangci-lint-$(GOLANGCI_LINT_VERSION)
JSON_CLI ?= ${BIN_DIR}/json-cli
SWAC ?= ${BIN_DIR}/swac

.PHONY: $(VENDOR_DIR)
$(VENDOR_DIR):
	@mkdir -p $VENDOR_DIR
	@$(GO) mod vendor
	@$(GO) mod tidy

.PHONY: deps
deps: $(VENDOR_DIR)

.PHONY: clean
clean:
	@rm -rf "${BIN_DIR}" "${VENDOR_DIR}"

.PHONY: lint
lint: $(GOLANGCI_LINT) $(VENDOR_DIR)
	@$(GOLANGCI_LINT) run -c .golangci.yaml

.PHONY: test
test: test-unit

## Run unit tests
.PHONY: test-unit
test-unit:
	@echo ">> unit test"
	@$(GO) test -gcflags=-l -coverprofile=unit.coverprofile -covermode=atomic -race ./...

.PHONY: generate-transaction
generate-transaction: $(JSON_CLI)
	@$(JSON_CLI) gen-go $(OPENAPI) \
		--patches patch-entities.json \
		--ptr-in-schema \
			'#/components/schemas/Transaction' \
		--def-ptr '#/components/schemas' \
		--package-name transaction \
		--name-tags csv \
		--output ./pkg/transaction/entity.go && \
		gofmt -w ./pkg/transaction/entity.go

.PHONY: generate-api
generate-api: $(JSON_CLI) $(SWAC)
	@rm -rf ./internal/api && \
		mkdir -p ./internal/api

	@$(SWAC) go-client $(OPENAPI) \
		--patches patch-client.json \
		--operations post/oauth/token,post/api/mfa/challenge,get/api/smrt/transactions \
		--skip-default-additional-properties \
		--out ./internal/api \
		--pkg-name api && \
		gofmt -w ./internal/api

.PHONY: generate
generate: generate-transaction generate-api

$(GOLANGCI_LINT):
	@echo "$(OK_COLOR)==> Installing golangci-lint $(GOLANGCI_LINT_VERSION)$(NO_COLOR)"; \
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin "$(GOLANGCI_LINT_VERSION)"
	@mv ./bin/golangci-lint $(GOLANGCI_LINT)

$(JSON_CLI): $(BIN_DIR)
	@curl -s -L 'https://github.com/swaggest/json-cli/releases/download/v$(JSON_CLI_VERSION)/json-cli' > $@
	@chmod +x $@

$(SWAC): $(BIN_DIR)
	@curl -s -L 'https://github.com/swaggest/swac/releases/download/v$(SWAC_VERSION)/swac' > $@
	@chmod +x $@
