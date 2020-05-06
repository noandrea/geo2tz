GOFILES = $(shell find . -name '*.go' -not -path './vendor/*')
GOPACKAGES = $(shell go list ./...  | grep -v /vendor/)
GIT_DESCR = $(shell git describe --tags --always) 
APP=geo2tz
# build output folder
OUTPUTFOLDER = dist
# docker image
DOCKER_REGISTRY = apeunit
DOCKER_IMAGE = geo2tz
DOCKER_TAG = $(GIT_DESCR)
# build paramters
OS = linux
ARCH = amd64
# K8S
K8S_NAMESPACE = geo
K8S_DEPLOYMENT = geo2tz

.PHONY: list
list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | xargs

default: build

workdir:
	mkdir -p dist

build: build-dist

build-dist: $(GOFILES)
	@echo build binary to $(OUTPUTFOLDER)
	GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0 go build -ldflags '-s -w -extldflags "-static" -X main.Version=$(GIT_DESCR)' -o $(OUTPUTFOLDER)/$(APP) .
	@echo copy resources
	cp -r README.md LICENSE $(OUTPUTFOLDER)
	@echo done

test: test-all

test-all:
	@go test -v $(GOPACKAGES) -coverprofile .testCoverage.txt

bench: bench-all

bench-all:
	@go test -bench -v $(GOPACKAGES)

lint: lint-all

lint-all:
	@golint -set_exit_status $(GOPACKAGES)

clean:
	@echo remove $(OUTPUTFOLDER) folder
	@rm -rf dist
	@echo done

docker: docker-build

docker-build:
	@echo copy resources
	docker build --build-arg DOCKER_TAG='$(GIT_DESCR)' -t $(DOCKER_IMAGE)  .
	@echo done

docker-push:
	@echo push image
	docker tag $(DOCKER_IMAGE):latest $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG)
	@echo done

docker-run: 
	docker run -p 2004:2004 $(DOCKER_IMAGE):latest

debug-start:
	@go run main.go start

k8s-deploy:
	@echo deploy k8s
	kubectl -n $(K8S_NAMESPACE) set image deployment/$(K8S_DEPLOYMENT) $(DOCKER_IMAGE)=$(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG)
	@echo done

k8s-rollback:
	@echo deploy k8s
	kubectl -n $(K8S_NAMESPACE) rollout undo deployment/$(K8S_DEPLOYMENT)
	@echo done

