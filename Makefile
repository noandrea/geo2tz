GOFILES = $(shell find . -name '*.go' -not -path './vendor/*')
GOPACKAGES = $(shell go list ./...  | grep -v /vendor/)
APP_VERSION = $(shell git describe --tags --always)
APP=geo2tz
# build output folder
OUTPUTFOLDER = dist
RELEASEFOLDER = release
# docker image
DOCKER_REGISTRY = docker.pkg.github.com/noandrea/geo2tz
DOCKER_IMAGE = geo2tz
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
	GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0 go build -ldflags '-s -w -extldflags "-static" -X main.Version=$(APP_VERSION)' -o $(OUTPUTFOLDER)/$(APP) .
	@echo copy resources
	cp -r README.md LICENSE $(OUTPUTFOLDER)
	@echo done


_check_version:
ifndef APP_VERSION
	$(error APP_VERSION is not set, please specifiy the version you want to tag)
endif


test: test-all

test-all:
	@go test -v $(GOPACKAGES) -coverprofile .testCoverage.txt

test-ci:
	go test -coverprofile=coverage.txt -covermode=atomic -race -mod=readonly $(GOPACKAGES)

bench: bench-all

bench-all:
	@go test -bench -v $(GOPACKAGES)

lint: lint-all

lint-all:
	@echo running linters
	staticcheck $(GOPACKAGES)
	golint -set_exit_status $(GOPACKAGES)
	@echo done

clean:
	@echo remove $(OUTPUTFOLDER) folder
	rm -rf $(OUTPUTFOLDER)
	@echo remove $(RELEASEFOLDER) folder
	rm -rf $(RELEASEFOLDER)
	@echo done

docker: docker-build

docker-build: _check_version
	@echo copy resources
	docker build --build-arg DOCKER_TAG='$(APP_VERSION)' -t $(DOCKER_IMAGE)  .
	@echo done

docker-push: _check_version
	@echo push image
	docker tag $(DOCKER_IMAGE):latest $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(APP_VERSION)
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(APP_VERSION)
	@echo done

docker-run: 
	docker run -p 2004:2004 $(DOCKER_IMAGE):latest

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

changelog:
	git-chglog --sort semver --output CHANGELOG.md


release-prepare: _check_version
	@echo making release $(APP_VERSION)
	git tag $(APP_VERSION)
	git-chglog --sort semver --output CHANGELOG.md
	git tag $(APP_VERSION) --delete
	git add CHANGELOG.md && git commit -m "chore: update changelog for $(APP_VERSION)"
	@echo release complete

git-tag: _check_version
ifneq ($(shell git rev-parse --abbrev-ref HEAD),main)
	$(error you are not on the main branch. aborting)
endif
	git tag -s -a "$(APP_VERSION)" -m "Changelog: https://github.com/noandrea/geo2tz/blob/main/CHANGELOG.md"

gh-publish-release: _check_version clean build
	@echo publish release
	mkdir -p $(RELEASEFOLDER)
	zip -rmT $(RELEASEFOLDER)/$(APP)-$(APP_VERSION).zip $(OUTPUTFOLDER)/
	sha256sum $(RELEASEFOLDER)/$(APP)-$(APP_VERSION).zip | tee $(RELEASEFOLDER)/$(APP)-$(APP_VERSION).zip.checksum
	gh release create $(APP_VERSION) $(RELEASEFOLDER)/* -t $(APP_VERSION) -F CHANGELOG.md
	@echo done
