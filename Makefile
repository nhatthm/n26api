JSON_CLI_VERSION = 1.8.3
SWAC_VERSION = 0.1.19

BINDIR = bin
VENDORDIR = vendor
OPENAPI = openapi.yaml

GO ?= go
GOLANGCI_LINT ?= golangci-lint
JSON_CLI ?= ${BINDIR}/json-cli
SWAC ?= ${BINDIR}/swac

$(VENDORDIR):
	@mkdir -p $VENDORDIR
	@$(GO) mod vendor
	@$(GO) mod tidy

$(BINDIR):
	@mkdir -p $@

$(JSON_CLI): $(BINDIR)
	@curl -s -L 'https://github.com/swaggest/json-cli/releases/download/v$(JSON_CLI_VERSION)/json-cli' > $@
	@chmod +x $@

$(SWAC): $(BINDIR)
	@curl -s -L 'https://github.com/swaggest/swac/releases/download/v$(SWAC_VERSION)/swac' > $@
	@chmod +x $@

.PHONY: $(VENDORDIR) deps lint test generate

deps: $(VENDORDIR)

clean:
	@rm -rf "${BINDIR}" "${VENDORDIR}"

lint:
	@$(GOLANGCI_LINT) run

test: test-unit

## Run unit tests
test-unit:
	@echo ">> unit test"
	@$(GO) test -gcflags=-l -coverprofile=unit.coverprofile -covermode=atomic -race ./...

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

generate: generate-transaction generate-api
