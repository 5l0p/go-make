# Go-Make Examples

This directory contains example Makefiles that demonstrate the capabilities of the go-make tool.

## Simple Example

The `simple/` directory contains a basic C program with a simple Makefile:

```bash
cd simple
../../go-make hello    # Build the hello program
../../go-make clean    # Clean up built files
../../go-make          # Build default target (hello)
```

## Complex Example

The `complex/` directory demonstrates a more sophisticated build system with:
- Multiple source files
- Object file generation
- Directory creation
- Variables and pattern rules

```bash
cd complex
../../go-make all      # Build the complete program
../../go-make clean    # Clean all build artifacts
../../go-make test     # Build and run tests
```

## Go Project Example

The `go-project/` directory shows how to use go-make with Go projects:

```bash
cd go-project
../../go-make build    # Build the Go binary
../../go-make test     # Run Go tests
../../go-make clean    # Clean build artifacts
../../go-make run      # Build and run the program
```

## Features Demonstrated

- **Basic target dependencies**: Simple file-to-file relationships
- **Pattern rules**: Using wildcards and variables
- **Directory targets**: Creating directories as dependencies
- **PHONY targets**: Targets that don't create files
- **Variable substitution**: Using Make-style variables
- **Command execution**: Running shell commands
- **File modification time checking**: Rebuilding only when necessary

## Testing Examples

You can test all examples by running the integration tests from the project root:

```bash
go test -v ./... -run Integration
```