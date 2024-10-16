COLOR_GREEN := \033[0;32m
COLOR_YELLOW := \033[0;33m
COLOR_CYAN := \033[0;36m
COLOR_OFF := \033[0m

GIT_USER_NAME := $(shell git config --get user.name)
GIT_USER_EMAIL := $(shell git config --get user.email)
HOST_NAME := $(shell hostname)


## help:print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## apt-update:apt update autoremove upgrade and install
.PHONY: apt-update
apt-update:
	@echo "🐧 ${COLOR_CYAN}apt clean all${COLOR_OFF}"
	@sudo apt -qq clean all
	@echo "🐧 ${COLOR_CYAN}apt update${COLOR_OFF}"
	@sudo apt -qq update
	@echo "🐧 ${COLOR_CYAN}apt autoremove${COLOR_OFF}"
	@sudo apt -qq -y autoremove
	@echo "🐧 ${COLOR_CYAN}apt upgrade${COLOR_OFF}"
	@sudo apt -qq -y upgrade
	@echo "🐧 ${COLOR_CYAN}apt install${COLOR_OFF}"
	@sudo apt -qq -y install bash-completion

## golang-vscode:install golang dependencies for vscode
.PHONY: golang-vscode
golang-vscode:
	@echo "🐿️  ${COLOR_CYAN}tools-vscode gopls${COLOR_OFF}"
	@go install golang.org/x/tools/gopls@latest
	@echo "🐿️  ${COLOR_CYAN}tools-vscode gotests${COLOR_OFF}"
	@go install github.com/cweill/gotests/gotests@latest
	@echo "🐿️  ${COLOR_CYAN}tools-vscode gomodifytags${COLOR_OFF}"
	@go install github.com/fatih/gomodifytags@latest
	@echo "🐿️  ${COLOR_CYAN}tools-vscode impl${COLOR_OFF}"
	@go install github.com/josharian/impl@latest
	@echo "🐿️  ${COLOR_CYAN}tools-vscode goplay${COLOR_OFF}"
	@go install github.com/haya14busa/goplay/cmd/goplay@latest
	@echo "🐿️  ${COLOR_CYAN}tools-vscode dlv${COLOR_OFF}"
	@go install github.com/go-delve/delve/cmd/dlv@latest

## golang-vscode:install golang tools
.PHONY: golang-tools
golang-tools:
	@echo "🐿️  ${COLOR_CYAN}tools-vscode staticcheck${COLOR_OFF}"
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@echo "🐿️  ${COLOR_CYAN}tools-vscode golangci-lint${COLOR_OFF}"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "🐿️  ${COLOR_CYAN}tools mockgen${COLOR_OFF}"
	@go install go.uber.org/mock/mockgen@latest

## requirements:everything you need for development
.PHONY: requirements
requirements: apt-update golang-vscode golang-tools
	@echo ""
ifeq ($(GIT_USER_NAME),)
	@echo "📣 Welcome ${COLOR_GREEN}$$USER@${HOST_NAME}${COLOR_OFF}"
else
	@echo "📣 Welcome ${COLOR_GREEN}${GIT_USER_NAME}${COLOR_OFF} ${COLOR_YELLOW}(${GIT_USER_EMAIL})${COLOR_OFF}"
endif
	@echo ""

## pre-push:tidy modfiles, go generate and format .go files
.PHONY: pre-push
pre-push:
	@go mod tidy
	@go generate ./...
	@go fmt ./...

## update-lib:update all go modules
.PHONY: update-lib
update-lib:
	@go get -u ./...
	@go mod tidy

## lint:run quality control checks
.PHONY: lint
lint:
	@go vet ./...
	@staticcheck ./...
	@golangci-lint run ./...

## test:run all tests
.PHONY: test
test:
	@go test -v -race ./...
	@go test -cover ./...
