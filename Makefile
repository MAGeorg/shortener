default: help

.PHONY: help
help: # Show help for each of the Makefile recipes.
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

.PHONY: linter
linter: # run golangci-lint with local file
	@golangci-lint run -c ./.golangci.yml || echo "Please fix lint warnings"

.PHONY: build
build: # build the server binary (stored in app/)
	go build -o app/server cmd/shortener/main.go

.PHONY: run
run: # run server with default parameters
	go run cmd/shortener/main.go