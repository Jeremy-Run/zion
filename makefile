GOROOT:=$(shell go env GOROOT)
GO=$(GOROOT)/bin/go
GORUN=$(GO) run

.PHONY: main
main:
	$(GORUN) ./ --factor=1 --dataSize=262144 --readTime=100000
