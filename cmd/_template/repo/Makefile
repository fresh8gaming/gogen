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

generate-proto:
	buf generate

check-proto:
	buf breaking --against '.git#branch=trunk'

# Run docker build locally for debuging
docker-build:
	docker build --file build/$(DOCKERFILE) --tag $(APP):latest --tag $(APP):$(GIT_SHA) . --build-arg APP=$(APP) --build-arg VERSION=$(GIT_SHA)

docker-push:
	docker push eu.gcr.io/intrepid-honor-109516/$(APP):$(GIT_SHA)