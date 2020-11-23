PROJECT_NAME=raida_go

GOBASE=$(shell pwd)
GOPATH=$(GOBASE)/vendor:$(GOBASE):/home/alexander/go
GOFILES=$(wildcard *.go)

all: build buildwin

build:
	GOPATH=$(GOPATH) go build -o $(PROJECT_NAME) -v $(GOFILES)  

buildwin:
	GOPATH=$(GOPATH) GOOS=windows GOARCH=amd64 go build  -o $(PROJECT_NAME).exe -v $(GOFILES) 


clean:
	go clean
	rm -f $(PROJECT_NAME)
	rm -f $(PROJECT_NAME).exe


deps:
	
