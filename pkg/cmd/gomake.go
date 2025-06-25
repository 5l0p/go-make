// Package cmd provides a high-level, convenient API for using go-make functionality.
// This package is designed to make it easy to embed Makefile processing in other applications
// with minimal setup and configuration.
package cmd

import (
	"fmt"
	"os"

	"go-make/pkg/builder"
	"go-make/pkg/makefile"
	"go-make/pkg/types"
)

// Make represents a high-level interface to go-make functionality.
// It combines parsing and building operations in a convenient API.
type Make struct {
	makefile *types.Makefile
	builder  *builder.Builder
}

// New creates a new Make instance by parsing the specified Makefile.
// If filename is empty, it defaults to "Makefile".
//
// Example:
//   make, err := cmd.New("Makefile")
//   if err != nil {
//       log.Fatal(err)
//   }
func New(filename string) (*Make, error) {
	if filename == "" {
		filename = "Makefile"
	}

	mf, err := makefile.ParseMakefile(filename)
	if err != nil {
		return nil, err
	}

	return &Make{
		makefile: mf,
		builder:  builder.NewBuilder(mf),
	}, nil
}

// NewFromMakefile creates a new Make instance from an existing parsed Makefile.
func NewFromMakefile(mf *types.Makefile) *Make {
	return &Make{
		makefile: mf,
		builder:  builder.NewBuilder(mf),
	}
}

// Build builds the specified target. If target is empty, builds the default target.
//
// Example:
//   err := make.Build("all")
//   if err != nil {
//       log.Fatal(err)
//   }
func (m *Make) Build(target string) error {
	if target == "" {
		target = m.makefile.FirstRule
	}

	if target == "" {
		return fmt.Errorf("no targets found in Makefile")
	}

	return m.builder.Build(target)
}

// BuildMultiple builds multiple targets in sequence.
// If any target fails, the process stops and returns an error.
//
// Example:
//   err := make.BuildMultiple("clean", "build", "test")
//   if err != nil {
//       log.Fatal(err)
//   }
func (m *Make) BuildMultiple(targets ...string) error {
	for _, target := range targets {
		if err := m.Build(target); err != nil {
			return fmt.Errorf("failed to build target '%s': %w", target, err)
		}
	}
	return nil
}

// HasTarget returns true if the Makefile contains the specified target.
func (m *Make) HasTarget(target string) bool {
	return m.makefile.HasTarget(target)
}

// Targets returns a list of all available targets.
func (m *Make) Targets() []string {
	return m.makefile.Targets()
}

// DefaultTarget returns the name of the default target (first target in Makefile).
func (m *Make) DefaultTarget() string {
	return m.makefile.FirstRule
}

// GetRule returns the Rule for the specified target, or nil if it doesn't exist.
func (m *Make) GetRule(target string) *types.Rule {
	return m.makefile.GetTarget(target)
}

// IsBuilt returns true if the target has been successfully built in this session.
func (m *Make) IsBuilt(target string) bool {
	return m.builder.IsBuilt(target)
}

// Reset clears the built state, allowing targets to be rebuilt.
func (m *Make) Reset() {
	m.builder.Reset()
}

// Makefile returns the underlying parsed Makefile for advanced usage.
func (m *Make) Makefile() *types.Makefile {
	return m.makefile
}

// Builder returns the underlying Builder for advanced usage.
func (m *Make) Builder() *builder.Builder {
	return m.builder
}

// Convenience functions for common operations

// BuildDefault builds the default target (first target in Makefile).
//
// Example:
//   err := make.BuildDefault()
//   if err != nil {
//       log.Fatal(err)
//   }
func (m *Make) BuildDefault() error {
	return m.Build("")
}

// BuildIfExists builds the target only if it exists in the Makefile.
// Returns nil if the target doesn't exist (no error).
//
// Example:
//   err := make.BuildIfExists("test")
//   if err != nil {
//       log.Fatal(err)
//   }
func (m *Make) BuildIfExists(target string) error {
	if !m.HasTarget(target) {
		return nil
	}
	return m.Build(target)
}

// Package-level convenience functions

// Build is a convenience function that parses a Makefile and builds a target in one call.
// If makefilePath is empty, it defaults to "Makefile".
// If target is empty, it builds the default target.
//
// Example:
//   err := cmd.Build("Makefile", "all")
//   if err != nil {
//       log.Fatal(err)
//   }
func Build(makefilePath, target string) error {
	make, err := New(makefilePath)
	if err != nil {
		return err
	}
	return make.Build(target)
}

// BuildDefault is a convenience function that parses a Makefile and builds the default target.
//
// Example:
//   err := cmd.BuildDefault("Makefile")
//   if err != nil {
//       log.Fatal(err)
//   }
func BuildDefault(makefilePath string) error {
	return Build(makefilePath, "")
}

// ListTargets is a convenience function that parses a Makefile and returns all targets.
//
// Example:
//   targets, err := cmd.ListTargets("Makefile")
//   if err != nil {
//       log.Fatal(err)
//   }
//   for _, target := range targets {
//       fmt.Println(target)
//   }
func ListTargets(makefilePath string) ([]string, error) {
	make, err := New(makefilePath)
	if err != nil {
		return nil, err
	}
	return make.Targets(), nil
}

// HasTarget is a convenience function that checks if a target exists in a Makefile.
//
// Example:
//   exists, err := cmd.HasTarget("Makefile", "test")
//   if err != nil {
//       log.Fatal(err)
//   }
//   if exists {
//       fmt.Println("test target exists")
//   }
func HasTarget(makefilePath, target string) (bool, error) {
	make, err := New(makefilePath)
	if err != nil {
		return false, err
	}
	return make.HasTarget(target), nil
}

// MustBuild is like Build but panics on error. Useful for simple scripts.
//
// Example:
//   cmd.MustBuild("Makefile", "all")
func MustBuild(makefilePath, target string) {
	if err := Build(makefilePath, target); err != nil {
		panic(err)
	}
}

// MustBuildDefault is like BuildDefault but panics on error. Useful for simple scripts.
//
// Example:
//   cmd.MustBuildDefault("Makefile")
func MustBuildDefault(makefilePath string) {
	if err := BuildDefault(makefilePath); err != nil {
		panic(err)
	}
}