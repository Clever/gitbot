SHELL := /bin/bash
PKG := github.com/Clever/gitbot
SUBPKGS := 
PKGS := $(PKG) $(SUBPKGS)
GODEP := $(GOPATH)/bin/godep
GOLINT := $(GOPATH)/bin/golint
.PHONY: deps test $(PKGS) build

test: $(PKGS)

$(GODEP):
	go get github.com/tools/godep

$(GOLINT):
	go get github.com/golang/lint/golint


deps: $(GODEP)
	go get ./...
	$(GODEP) save -r

build:
	go build $(PKG)

$(PKGS): $(GOLINT) deps
	gofmt -w=true $(GOPATH)/src/$@/*.go
	$(GOLINT) $(GOPATH)/src/$@/*.go
ifeq ($(COVERAGE),1)
	go test -cover -coverprofile=$(GOPATH)/src/$@/c.out $@ -test.v
	go tool cover -html=$(GOPATH)/src/$@/c.out
else
	go test -v $@
endif
