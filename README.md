# Go-Make

A basic implementation of GNU Make written in Go. This project provides a simplified version of the classic `make` build automation tool, capable of parsing Makefiles and executing build targets with dependency resolution.

## Features

- **Makefile Parsing**: Supports standard Makefile syntax with targets, dependencies, and commands
- **Dependency Resolution**: Automatically resolves and builds dependencies in the correct order
- **File Timestamp Checking**: Only rebuilds targets when dependencies are newer
- **Circular Dependency Detection**: Detects and reports circular dependencies
- **Shell Command Execution**: Executes shell commands for each target
- **Default Target Selection**: Automatically selects the first target when none specified
- **PHONY Target Support**: Supports targets that don't correspond to files

## Project Structure

```
go-make/
├── cmd/go-make/          # Main application entry point
│   └── main.go
├── pkg/                  # Public packages (importable by other projects)
│   ├── cmd/              # High-level convenience API
│   │   └── gomake.go
│   ├── types/            # Core data structures
│   │   ├── makefile.go
│   │   └── makefile_test.go
│   ├── makefile/         # Makefile parsing functionality
│   │   ├── parser.go
│   │   └── parser_test.go
│   └── builder/          # Build execution and dependency resolution
│       ├── builder.go
│       └── builder_test.go
├── examples/             # Example Makefiles and projects
│   ├── simple/           # Basic C program example
│   ├── complex/          # Multi-file C project example
│   ├── go-project/       # Go project example
│   └── README.md
├── integration_test.go   # Integration tests
├── Makefile             # Project build configuration
├── go.mod               # Go module definition
└── README.md            # This file
```

## Installation

### From Source

```bash
git clone <repository-url>
cd go-make
make build
```

The binary will be built in `./bin/go-make`.

### Using Go Install

```bash
go install ./cmd/go-make
```

## Usage

### Command Line Usage

```bash
# Build the default target (first target in Makefile)
go-make

# Build a specific target
go-make target-name

# Examples
go-make all
go-make clean
go-make test
```

### Library Usage

You can use go-make as a library in your Go programs. There are two approaches:

#### High-Level API (Recommended)

The easiest way is to use the `pkg/cmd` package:

```go
package main

import (
    "log"
    
    "go-make/pkg/cmd"
)

func main() {
    // Simple one-liner to build a target
    err := cmd.Build("Makefile", "all")
    if err != nil {
        log.Fatal(err)
    }
    
    // Or build the default target
    err = cmd.BuildDefault("Makefile")
    if err != nil {
        log.Fatal(err)
    }
}
```

#### Advanced Usage with High-Level API

```go
package main

import (
    "fmt"
    "log"
    
    "go-make/pkg/cmd"
)

func main() {
    // Create a Make instance for multiple operations
    make, err := cmd.New("Makefile")
    if err != nil {
        log.Fatal(err)
    }
    
    // List available targets
    fmt.Println("Available targets:", make.Targets())
    
    // Build multiple targets
    err = make.BuildMultiple("clean", "build", "test")
    if err != nil {
        log.Fatal(err)
    }
    
    // Build only if target exists
    err = make.BuildIfExists("optional-target")
    if err != nil {
        log.Fatal(err)
    }
    
    // Check build status
    if make.IsBuilt("test") {
        fmt.Println("Tests completed successfully")
    }
}
```

#### Low-Level API

For fine-grained control, use the individual packages:

```go
package main

import (
    "log"
    
    "go-make/pkg/makefile"
    "go-make/pkg/builder"
)

func main() {
    // Parse a Makefile
    mf, err := makefile.ParseMakefile("Makefile")
    if err != nil {
        log.Fatal(err)
    }
    
    // Create a builder
    b := builder.NewBuilder(mf)
    
    // Build a specific target
    err = b.Build("all")
    if err != nil {
        log.Fatal(err)
    }
    
    // Check if a target was built
    if b.IsBuilt("all") {
        log.Println("Target 'all' was successfully built")
    }
}
```

#### Available Packages

- **`pkg/cmd`**: High-level convenience API (recommended for most users)
- **`pkg/types`**: Core data structures (Makefile, Rule)
- **`pkg/makefile`**: Makefile parsing functionality
- **`pkg/builder`**: Build execution and dependency resolution

#### Example: Custom Build Logic

```go
package main

import (
    "fmt"
    "log"
    
    "go-make/pkg/makefile"
    "go-make/pkg/builder"
)

func main() {
    // Parse Makefile
    mf, err := makefile.ParseMakefile("Makefile")
    if err != nil {
        log.Fatal(err)
    }
    
    // List all available targets
    fmt.Println("Available targets:")
    for _, target := range mf.Targets() {
        rule := mf.GetTarget(target)
        fmt.Printf("  %s (deps: %v)\n", target, rule.Dependencies)
    }
    
    // Build with custom logic
    b := builder.NewBuilder(mf)
    
    // Build multiple targets
    targets := []string{"clean", "build", "test"}
    for _, target := range targets {
        if mf.HasTarget(target) {
            fmt.Printf("Building target: %s\n", target)
            if err := b.Build(target); err != nil {
                log.Printf("Failed to build %s: %v", target, err)
                continue
            }
        } else {
            fmt.Printf("Target %s not found, skipping\n", target)
        }
    }
    
    // Reset builder state if needed
    b.Reset()
}
```

### Example Makefile

```makefile
# Simple C program Makefile
hello: hello.c
	gcc -o hello hello.c
	echo "Built hello program"

clean:
	rm -f hello
	echo "Cleaned up"

.PHONY: clean
```

## Development

### Building

```bash
# Build the project
make build

# Build for development (puts binary in project root)
make dev-build

# Clean build artifacts
make clean
```

### Testing

```bash
# Run unit tests
make test

# Run all tests including integration tests
make integration-test

# Run specific test packages
go test ./internal/builder
go test ./internal/makefile
```

### Code Quality

```bash
# Format code
make fmt

# Lint code
make lint

# Run vet
make vet
```

## Examples

The `examples/` directory contains sample projects demonstrating various use cases:

- **simple/**: Basic C program with simple compilation
- **complex/**: Multi-file C project with object files and directories
- **go-project/**: Go application with typical Go build tasks

See [examples/README.md](examples/README.md) for detailed information about each example.

## Supported Makefile Features

### Implemented

- Target definitions with dependencies
- Shell command execution
- File timestamp-based rebuilding
- Comment parsing (lines starting with `#`)
- PHONY targets
- Default target selection
- Circular dependency detection

### Not Yet Implemented

- Variable substitution (`$(VAR)` or `${VAR}`)
- Pattern rules (`%.o: %.c`)
- Built-in functions
- Conditional statements (`ifeq`, `ifdef`, etc.)
- Include directives
- Automatic variables (`$@`, `$<`, `$^`, etc.)

## Architecture

The application is structured using standard Go project layout:

- **`cmd/go-make`**: Main application with CLI interface
- **`pkg/cmd`**: High-level convenience API for easy integration
- **`pkg/makefile`**: Public Makefile parsing API
- **`pkg/builder`**: Public build execution and dependency resolution API
- **`pkg/types`**: Public types and data structures
- **`examples/`**: Example projects demonstrating usage

### Key Components

1. **Parser** (`pkg/makefile`): Parses Makefile syntax into structured data
2. **Builder** (`pkg/builder`): Executes build process with dependency resolution
3. **Types** (`pkg/types`): Defines core data structures (Rule, Makefile)

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (`make test`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## Testing

The project includes comprehensive testing:

- **Unit Tests**: Test individual components in isolation
- **Integration Tests**: Test complete workflows with real examples
- **Example Tests**: Verify example projects work correctly

Run all tests with:
```bash
make integration-test
```

## License

This project is provided as-is for educational purposes. See individual file headers for any specific licensing information.

## Acknowledgments

- Inspired by GNU Make
- Built as a learning exercise in Go and build systems
- Thanks to the Go community for excellent tooling and documentation