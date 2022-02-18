GOROOT:=$(shell go env GOROOT)
GO=$(GOROOT)/bin/go
GORUN=$(GO) run

.PHONY: main
main:
	$(GORUN) ./