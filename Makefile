SHELL := /bin/bash
PKG := github.com/Clever/gitbot
SUBPKGS := 
PKGS := $(PKG) $(SUBPKGS)
GODEP := $(GOPATH)/bin/godep
GOLINT := $(GOPATH)/bin/golint
EXECUTABLE := gitbot
VERSION := $(shell cat VERSION)
BUILDS := \
	build/$(EXECUTABLE)-v$(VERSION)-darwin-amd64 \
	build/$(EXECUTABLE)-v$(VERSION)-linux-amd64
COMPRESSED_BUILDS := $(BUILDS:%=%.tar.gz)
RELEASE_ARTIFACTS := $(COMPRESSED_BUILDS:build/%=release/%)

GOVERSION := $(shell go version | grep 1.5)
ifeq "$(GOVERSION)" ""
  $(error must be running Go version 1.5)
endif
export GO15VENDOREXPERIMENT = 1


.PHONY: test $(PKGS) build clean

test: $(PKGS)

$(GODEP):
	go get github.com/tools/godep

$(GOLINT):
	go get github.com/golang/lint/golint

build:
	go build $(PKG)

$(PKGS): $(GOLINT) version.go
	gofmt -w=true $(GOPATH)/src/$@/*.go
	$(GOLINT) $(GOPATH)/src/$@/*.go
	go vet $@
	go test -v $@

build/*: version.go
version.go: VERSION
	echo 'package main' > version.go
	echo '' >> version.go
	echo '// Version of gitbot' >> version.go
	echo 'const Version = "$(VERSION)"' >> version.go
build/$(EXECUTABLE)-v$(VERSION)-darwin-amd64:
	GOARCH=amd64 GOOS=darwin go build -o "$@/$(EXECUTABLE)"
build/$(EXECUTABLE)-v$(VERSION)-linux-amd64:
	GOARCH=amd64 GOOS=linux go build -o "$@/$(EXECUTABLE)"

%.tar.gz: %
	tar -C `dirname $<` -zcvf "$<.tar.gz" `basename $<`

$(RELEASE_ARTIFACTS): release/% : build/%
	mkdir -p release
	cp $< $@

release: $(RELEASE_ARTIFACTS)

clean:
	rm -rf build release


SHELL := /bin/bash
PKGS := $(shell go list ./... | grep -v /vendor)
GODEP := $(GOPATH)/bin/godep

$(GODEP):
	go get -u github.com/tools/godep

vendor: $(GODEP)
	$(GODEP) save $(PKGS)
	find vendor/ -path '*/vendor' -type d | xargs -IX rm -r X # remove any nested vendor directories
