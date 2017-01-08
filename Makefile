# メタ情報
NAME     := fod
VERSION  := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS  := \
	-X 'main.name=$(NAME)' \
	-X 'main.version=$(VERSION)' \
	-X 'main.revision=$(REVISION)'

# 必要なツール類をセットアップする
## Setup
setup:
	go get github.com/Masterminds/glide
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/goimports
	go get github.com/Songmu/make2help/cmd/make2help
  
# テストを実行する
## Run tests
test: deps
	go test $${ARGS} $$(glide novendor)

# glideを使って依存パッケージをインストールする  
## Install dependencies
deps: setup
	glide install
 
## Update dependencies
update: setup
	glide update

## Lint
lint: setup
	go vet -n -x $$(glide novendor)
	for pkg in $$(glide novendor -x); do \
		golint -set_exit_status $$pkg || exit $$?; \
	done
  
## Format source codes
fmt: setup
	goimports -w $$(glide novendor -x)

## build binaries ex. make bin/myproj
bin/%: cmd/%/main.go deps
	go build -ldflags "$(LDFLAGS)" -o $@ $<
  
## Show help
help:
	@make2help $(MAKEFILE_LIST)
    
.PHONY: setup deps update test lint help