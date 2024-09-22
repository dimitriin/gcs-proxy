GO_VERSION ?= 1.22
GO_DOCKER_IMG ?= golang:${GO_VERSION}-bookworm
GOLANG_CI_LINT_VERSION ?= v1.61.0

SERVICE ?= gcs-proxy
PROJECT ?= github.com/dimitriin/gcs-proxy
PROJECT_DIR ?= $(shell pwd)
GCS_PROXY_DOCKER_IMG_REPO ?= gcs-proxy
GCS_PROXY_DOCKER_IMG_TAG ?= latest

COVER_PROFILE ?=

GOOS ?= linux
GOARCH ?= amd64
BIN_PATH ?= ./bin

RELEASE ?= v0.0.0
RELEASE_DATE=$(shell date +%FT%T%Z)
GIT_REVISION = $(shell git rev-parse HEAD)
GIT_REPO_INFO=$(shell git config --get remote.origin.url)
GIT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

LDFLAGS = "-s -w \
	-X $(PROJECT)/pkg/version.RELEASE=$(RELEASE) \
	-X $(PROJECT)/pkg/version.DATE=$(RELEASE_DATE) \
	-X $(PROJECT)/pkg/version.REPO=$(GIT_REPO_INFO) \
	-X $(PROJECT)/pkg/version.COMMIT=$(GIT_REVISION) \
	-X $(PROJECT)/pkg/version.BRANCH=$(GIT_BRANCH)"

all: lintC testC
.PHONY: all

gocacheC:
	@docker volume create golang-${GO_VERSION}-mod-cache
	@docker volume create golang-${GO_VERSION}-build-cache
.PHONY: gocacheC

test:
	@echo "+ $@"
	@go test ./...
.PHONY: test

testC: gocacheC
	@echo "+ $@"
	docker run --rm -i  \
		-v ${PROJECT_DIR}:/go/src/${PROJECT} \
		--mount source=golang-${GO_VERSION}-mod-cache,target=/go/pkg/mod \
		--mount source=golang-${GO_VERSION}-build-cache,target=/root/.cache/go-build \
		-w /go/src/${PROJECT} ${GO_DOCKER_IMG} make test
.PHONY: testC

lintcacheC:
	@docker volume create golangci-lint-${GOLANG_CI_LINT_VERSION}-cache

lintC: gocacheC lintcacheC
	@echo "+ $@"
	docker run --rm \
		-v ${PROJECT_DIR}:/app \
		--mount source=golangci-lint-${GOLANG_CI_LINT_VERSION}-cache,target=/root/.cache \
		--mount source=golang-${GO_VERSION}-mod-cache,target=/go/pkg/mod \
		-w /app golangci/golangci-lint:${GOLANG_CI_LINT_VERSION} \
		golangci-lint run \
			--enable-all \
			--timeout=10m \
			--max-issues-per-linter 0 \
			--max-same-issues 0 \
			./...
.PHONY: lintC


build:
	@echo "+ $@ ${GOOS}/${GOARCH}"
	@GOOS=${GOOS} GOARCH=${GOARCH} go build -a -ldflags ${LDFLAGS} -o ${BIN_PATH}/${SERVICE}-${GOOS}-${GOARCH} ./cmd/${SERVICE}
	@echo "+ binary ${BIN_PATH}/${SERVICE}"
.PHONY: build

buildC:
	@echo "+ build image ${GCS_PROXY_DOCKER_IMG_REPO}:${GCS_PROXY_DOCKER_IMG_TAG}"
	@docker build -t ${GCS_PROXY_DOCKER_IMG_REPO}:${GCS_PROXY_DOCKER_IMG_TAG} .
.PHONY: buildC


