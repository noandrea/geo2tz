GOPACKAGES = $(shell go list ./...  | grep -v /vendor/)
APP=geo2tz

.PHONY: list
list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | xargs

default: build

workdir:
	mkdir -p dist

build: build-dist

build-dist: $(GOFILES)
	@echo build binary
	goreleaser build --single-target --config .github/.goreleaser.yaml --snapshot --clean
	@echo done


_check_version:
ifndef APP_VERSION
	$(error APP_VERSION is not set, please specifiy the version you want to tag)
endif


test: test-all

test-all:
	@go test -v $(GOPACKAGES) -race -covermode=atomic -coverprofile coverage.txt

test-coverage:
	go test -mod=readonly -coverprofile=coverage.out -covermode=atomic -timeout 30s $(GOPACKAGES) && \
	go tool cover -html=coverage.out

test-ci:
	go run main.go update current
	go test -coverprofile=coverage.txt -covermode=atomic -race -mod=readonly $(GOPACKAGES)

bench: bench-all

bench-all:
	@go test -bench -v $(GOPACKAGES)

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	GO111MODULE=on go mod verify

lint:
	@echo "--> Running linter"
	golangci-lint run --config .github/.golangci.yaml
	@go mod verify

debug-start:
	@go run main.go start

k8s-deploy: _check_version
	@echo deploy k8s
	kubectl -n $(K8S_NAMESPACE) set image deployment/$(K8S_DEPLOYMENT) $(DOCKER_IMAGE)=$(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(APP_VERSION)
	@echo done

k8s-rollback:
	@echo deploy k8s
	kubectl -n $(K8S_NAMESPACE) rollout undo deployment/$(K8S_DEPLOYMENT)
	@echo done


update-tzdata:
	@echo "--> Updating timzaone data"
	@echo build binary
	goreleaser build --single-target --config .github/.goreleaser.yaml --snapshot --clean -o geo2tz
	./geo2tz update latest
	@echo done
