export PATH := "$(PWD)/bin:$(PATH)"
export GOBIN := $(PWD)/bin

lint:
	GOROOT=$(GOROOT) golangci-lint run --verbose -c .golangci.yml --timeout 300s

test:
	gotestsum -- -race ./... -coverprofile coverage.out

mutation-test:
	go-mutesting --config=.ci/go-mutesting/config.yml --blacklist=.ci/go-mutesting/blacklist.txt ./...

update-dependencies:
	go get -u ./...
	go mod tidy
	go mod vendor

install-tools:
	@echo Installing tools from tools.go
	@cat tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

# Run docker build locally for debuging
docker-build:
	docker build --file build/$(DOCKERFILE) --tag $(APP):latest --tag $(APP):$(GIT_SHA) . --build-arg APP=$(APP) --build-arg VERSION=$(GIT_SHA)

docker-push:
	docker push eu.gcr.io/intrepid-honor-109516/$(APP):$(GIT_SHA)
