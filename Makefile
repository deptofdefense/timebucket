# =================================================================
#
# Work of the U.S. Department of Defense, Defense Digital Service.
# Released as open source under the MIT License.  See LICENSE file.
#
# =================================================================
.PHONY: help
help:  ## Print the help documentation
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

#
# Go building, formatting, testing, and installing
#

fmt:  ## Format Go source code
	go fmt $$(go list ./... )

.PHONY: imports
imports: bin/goimports ## Update imports in Go source code
	# If missing, install goimports with: go get golang.org/x/tools/cmd/goimports
	bin/goimports -w -local github.com/deptofdefense/timebucket,github.com/deptofdefense $$(find . -iname '*.go')

vet: ## Vet Go source code
	go vet $$(go list ./...)

tidy: ## Tidy Go source code
	go mod tidy

.PHONY: test_go
test_go: bin/errcheck bin/ineffassign bin/staticcheck bin/shadow ## Run Go tests
	bash scripts/test.sh

.PHONY: test_cli
test_cli: bin/timebucket ## Run CLI tests
	bash scripts/test-cli.sh

install:  ## Install timebucket CLI on current platform
	go install github.com/deptofdefense/timebucket/cmd/timebucket

#
# Command line Programs
#

bin/errcheck:
	go build -o bin/errcheck github.com/kisielk/errcheck

bin/goimports:
	go build -o bin/goimports golang.org/x/tools/cmd/goimports

bin/gox:
	go build -o bin/gox github.com/mitchellh/gox

bin/ineffassign:
	go build -o bin/ineffassign github.com/gordonklaus/ineffassign

bin/staticcheck:
	go build -o bin/staticcheck honnef.co/go/tools/cmd/staticcheck

bin/shadow:
	go build -o bin/shadow golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow

bin/timebucket: ## Build timebucket CLI for Darwin / amd64
	go build -o bin/timebucket github.com/deptofdefense/timebucket/cmd/timebucket

bin/timebucket_linux_amd64: bin/gox ## Build timebucket CLI for Darwin / amd64
	scripts/build-release linux amd64

.PHONY: build
build: bin/timebucket

.PHONY: build_release
build_release: bin/gox
	scripts/build-release

#
# Docker
#

docker_build: bin/timebucket_linux_amd64 ## Build docker image
	docker build -f Dockerfile --tag timebucket:latest .

## Clean

clean:  ## Clean artifacts
	rm -fr bin
