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

.PHONY: deps test $(PKGS) build clean

test: $(PKGS)

$(GODEP):
	go get github.com/tools/godep

$(GOLINT):
	go get github.com/golang/lint/golint


deps: $(GODEP)
	go get $(PKGS)
	$(GODEP) save -r

build:
	go build $(PKG)

$(PKGS): $(GOLINT) version.go deps
	gofmt -w=true $(GOPATH)/src/$@/*.go
	$(GOLINT) $(GOPATH)/src/$@/*.go
ifeq ($(COVERAGE),1)
	go test -cover -coverprofile=$(GOPATH)/src/$@/c.out $@ -test.v
	go tool cover -html=$(GOPATH)/src/$@/c.out
else
	go test -v $@
endif

build/*: version.go
version.go: VERSION
	echo 'package main' > version.go
	echo '' >> version.go
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
