PROJECT=pkgind
ORGANIZATION=hectorj2f

BINARY_SERVER=bin/$(PROJECT)
BINARY_CTL=bin/$(PROJECT)ctl

SOURCE := $(shell find . -name '*.go')
VERSION := $(shell cat VERSION)
COMMIT := $(shell git rev-parse --short HEAD)
GOPATH := $(shell pwd)/.gobuild
PROJECT_PATH := $(GOPATH)/src/github.com/$(ORGANIZATION)

ifndef GOOS
  GOOS := $(shell go env GOOS)
endif
ifndef GOARCH
  GOARCH := $(shell go env GOARCH)
endif

.PHONY: all clean vendor-clean vendor-update

all: .gobuild $(BINARY_SERVER) $(BINARY_CTL)

.gobuild:
	mkdir -p $(PROJECT_PATH)
	mkdir -p $(GOPATH)/doc
	cd $(PROJECT_PATH) && ln -s ../../../.. $(PROJECT)

	#
	# Fetch public dependencies via `go get`
	# All of the dependencies are listed here
	@GOPATH=$(GOPATH) go get github.com/spf13/cobra
	@GOPATH=$(GOPATH) go get github.com/golang/glog
	@GOPATH=$(GOPATH) go get github.com/op/go-logging

test: .gobuild
	docker run \
	    --rm \
	    -v $(shell pwd):/usr/code \
	    -e GOPATH=/usr/code/.gobuild \
	    -e GOOS=$(GOOS) \
	    -e GOARCH=$(GOARCH) \
	    -e GO15VENDOREXPERIMENT=1 \
	    -w /usr/code/ \
		golang:1.5 \
	    bash -c 'cd .gobuild/src/github.com/hectorj2f/pkgind && go test $$(go list ./... | grep -v vendor)'

$(BINARY_SERVER): $(SOURCE) VERSION .gobuild
	@echo Building for $(GOOS)/$(GOARCH)
	docker run \
	    --rm \
	    -v $(shell pwd):/usr/code \
	    -e GOPATH=/usr/code/.gobuild \
	    -e GOOS=$(GOOS) \
	    -e GOARCH=$(GOARCH) \
	    -e GO15VENDOREXPERIMENT=1 \
	    -w /usr/code \
      golang:1.5 \
	    go build -a -ldflags "-X github.com/hectorj2f/pkgind/cmd.projectVersion=$(VERSION) -X github.com/hectorj2f/pkgind/cmd.projectBuild=$(COMMIT)" -o $(BINARY_SERVER) github.com/$(ORGANIZATION)/$(PROJECT)

$(BINARY_CTL): $(SOURCE) VERSION .gobuild
	docker run \
	    --rm \
	    -v $(shell pwd):/usr/code \
	    -e GOPATH=/usr/code/.gobuild \
	    -e GOOS=$(GOOS) \
	    -e GOARCH=$(GOARCH) \
	    -e GO15VENDOREXPERIMENT=1 \
	    -w /usr/code \
      golang:1.5 \
	    go build -a -ldflags "-X github.com/hectorj2f/pkgind/client.projectVersion=$(VERSION) -X github.com/hectorj2f/pkgind/client.projectBuild=$(COMMIT)" -o $(BINARY_CTL) github.com/$(ORGANIZATION)/$(PROJECT)/pkgindctl

distclean: clean clean-bin-dist

clean:
	rm -rf .gobuild bin

vendor-clean:
	rm -rf vendor/

vendor-update: vendor-clean
	rm -rf glide.lock
	GO15VENDOREXPERIMENT=1 glide install
	find vendor/ -name .git -type d -prune | xargs rm -rf

install: $(BINARY_SERVER) $(BINARY_CTL)
	cp $(BINARY_SERVER) $(BINARY_CTL) /usr/local/bin/

ci: clean all test

fmt:
	gofmt -l -w .

lint:
	GOPATH=$(GOPATH) go vet $(go list ./... | grep -v "gopath")
	GOPATH=$(GOPATH) golint $(go list ./... | grep -v "gopath")


godoc: all
	@echo Opening godoc server at http://localhost:6060/pkg/github.com/$(ORGANIZATION)/$(PROJECT)/
	docker run \
	    --rm \
	    -v $(shell pwd):/usr/code \
	    -e GOPATH=/usr/code/.gobuild \
	    -e GOROOT=/usr/code/.gobuild \
	    -e GOOS=$(GOOS) \
	    -e GOARCH=$(GOARCH) \
	    -e GO15VENDOREXPERIMENT=1 \
	    -w /usr/code \
      -p 6060:6060 \
		golang:1.5 \
		godoc -http=:6060
