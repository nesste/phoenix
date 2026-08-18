GO ?= go
VERSION ?= dev
TARGET_GOOS ?= linux
TARGET_GOARCH ?= amd64

export CGO_ENABLED := 0
export GOTOOLCHAIN := go1.26.6

ifeq ($(OS),Windows_NT)
TARGET_ENV := set GOOS=$(TARGET_GOOS)&& set GOARCH=$(TARGET_GOARCH)&&
CREATE_BIN := if not exist bin mkdir bin
else
TARGET_ENV := GOOS=$(TARGET_GOOS) GOARCH=$(TARGET_GOARCH)
CREATE_BIN := mkdir -p bin
endif

STATICCHECK_VERSION := v0.7.0
GOVULNCHECK_VERSION := v1.7.0
GOCYCLO_VERSION := v0.6.0
DUPL_VERSION := v1.1.0

.PHONY: test lint build validate-spec validate-authoring quality format-check dependency-check complexity manifest

test:
	$(GO) test -count=1 ./...

format-check:
	$(GO) run ./cmd/quality-check format cmd internal verbs experiments/frontier-v1/runner

lint: format-check
	$(GO) vet ./...
	$(GO) run honnef.co/go/tools/cmd/staticcheck@$(STATICCHECK_VERSION) ./...

dependency-check:
	$(GO) mod verify
	$(GO) run golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION) ./...

complexity:
	$(GO) run github.com/fzipp/gocyclo/cmd/gocyclo@$(GOCYCLO_VERSION) -over 15 cmd internal verbs experiments/frontier-v1/runner
	$(GO) run ./cmd/quality-check no-output -- $(GO) run github.com/mibk/dupl@$(DUPL_VERSION) -plumbing -t 100 cmd internal verbs experiments/frontier-v1/runner

validate-spec:
	cd experiments/surface-spike && $(GO) run ./cmd/validate-spec --repo-root ../..

validate-authoring:
	cd experiments/frontier-v1/corpusctl && $(GO) run ./cmd/corpusctl validate --repo-root ../../.. --tranche authoring --world-source experiments/frontier-v1/worlds/authoring.dev_repo.json

build:
	$(CREATE_BIN)
	$(TARGET_ENV) $(GO) build -trimpath -buildvcs=false -ldflags="-s -w -buildid= -X main.version=$(VERSION)" -o bin/phoenix ./cmd/phoenix

manifest: build
	$(GO) run ./cmd/build-manifest \
		--repo-root . \
		--output build/manifest.json \
		--executable bin/phoenix \
		--verb internal/verb/exec.go \
		--verb internal/verb/external.go \
		--verb internal/verb/registry.go \
		--verb verbs/dev-repo/commands.go \
		--verb verbs/dev-repo/episodes.go \
		--verb verbs/dev-repo/register.go \
		--verb verbs/dev-repo/repository.go \
		--verb verbs/dev-repo/schemas.go \
		--rule verbs/dev-repo/refusals.go \
		--rule verbs/dev-repo/transitions.go \
		--world-definition worlds/dev-repo/world.json \
		--schema spec/result.schema.json \
		--schema spec/episode.schema.json \
		--schema spec/world.schema.json \
		--goos $(TARGET_GOOS) \
		--goarch $(TARGET_GOARCH) \
		--cgo-enabled=false

quality: test lint dependency-check complexity validate-spec validate-authoring manifest
