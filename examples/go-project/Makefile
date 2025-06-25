# Go project example
all: build

build: myapp

myapp: ./cmd/main.go go.mod
	go build -o myapp ./cmd/main.go
	echo "Built Go binary: myapp"

test:
	go test ./...
	echo "Go tests completed"

fmt:
	go fmt ./...
	echo "Go code formatted"

clean:
	rm -f myapp
	go clean
	echo "Cleaned Go project"

run: myapp
	./myapp

deps:
	go mod tidy
	go mod download
	echo "Dependencies updated"

.PHONY: all build test fmt clean run deps