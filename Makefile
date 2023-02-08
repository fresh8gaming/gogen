export PATH := $(PWD)/bin:$(PATH)
export GOBIN := $(PWD)/bin

lint:
	golangci-lint run -c .golangci.yml -v

test:
	gotestsum -- -race ./...

update-dependencies:
	go get -u ./...
	go mod tidy
	go mod vendor

install-tools:
	@echo Installing tools from tools.go
	@cat tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %
