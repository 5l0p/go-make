# Go-Make Project Makefile
.PHONY: all build test clean install lint fmt vet integration-test

all: test build

build: ./bin/go-make

./bin/go-make: ./cmd/go-make/main.go pkg/makefile/parser.go pkg/builder/builder.go pkg/types/makefile.go pkg/cmd/gomake.go
	mkdir -p ./bin
	go build -o ./bin/go-make ./cmd/go-make

test:
	go test -v -short ./...

integration-test: build
	go test -v ./... -run Integration

lint:
	go vet ./...
	gofmt -l .

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf ./bin
	rm -f go-make
	go clean

install: build
	cp ./bin/go-make $(GOPATH)/bin/

dev-build: 
	go build -o go-make ./cmd/go-make