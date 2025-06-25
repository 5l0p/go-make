// Package builder provides functionality for building targets from parsed Makefiles.
package builder

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/5l0p/go-make/pkg/types"
)

// Builder handles the build process for Makefile targets.
// It manages dependency resolution, file timestamp checking, and command execution.
type Builder struct {
	makefile *types.Makefile
	built    map[string]bool
	building map[string]bool
}

// NewBuilder creates a new Builder instance for the given Makefile.
//
// Example usage:
//   makefile, err := makefile.ParseMakefile("Makefile")
//   if err != nil {
//       log.Fatal(err)
//   }
//   
//   builder := NewBuilder(makefile)
//   err = builder.Build("all")
//   if err != nil {
//       log.Fatal(err)
//   }
func NewBuilder(makefile *types.Makefile) *Builder {
	return &Builder{
		makefile: makefile,
		built:    make(map[string]bool),
		building: make(map[string]bool),
	}
}

// Build builds the specified target and all its dependencies.
// It will:
//   - Resolve dependencies recursively
//   - Check file timestamps to determine if rebuilding is needed
//   - Execute commands for targets that need rebuilding
//   - Detect circular dependencies
//
// Returns an error if:
//   - A circular dependency is detected
//   - A target has no rule and doesn't exist as a file
//   - A command execution fails
func (b *Builder) Build(target string) error {
	// If already built, skip
	if b.built[target] {
		return nil
	}

	// Detect circular dependencies
	if b.building[target] {
		return fmt.Errorf("circular dependency detected involving target '%s'", target)
	}

	rule, exists := b.makefile.Rules[target]
	if !exists {
		// If no rule exists, check if it's a file
		if b.fileExists(target) {
			return nil
		}
		return fmt.Errorf("no rule to make target '%s'", target)
	}

	// Mark as currently building
	b.building[target] = true

	// Build all dependencies first
	for _, dep := range rule.Dependencies {
		if err := b.Build(dep); err != nil {
			return err
		}
	}

	// Check if target needs rebuilding
	if b.needsRebuild(target, rule.Dependencies) {
		fmt.Printf("Building target: %s\n", target)
		
		// Create automatic variables context
		autoVars := b.createAutomaticVariables(target, rule.Dependencies)
		
		for _, command := range rule.Commands {
			if err := b.executeCommandWithContext(command, autoVars); err != nil {
				return fmt.Errorf("command failed: %s", err)
			}
		}
	}

	// Mark as no longer building and as built
	b.building[target] = false
	b.built[target] = true
	return nil
}

// IsBuilt returns true if the target has been successfully built in this session.
func (b *Builder) IsBuilt(target string) bool {
	return b.built[target]
}

// Reset clears the built state, allowing targets to be rebuilt.
func (b *Builder) Reset() {
	b.built = make(map[string]bool)
	b.building = make(map[string]bool)
}

// needsRebuild determines if a target needs to be rebuilt based on dependency timestamps.
// A target needs rebuilding if:
//   - The target file doesn't exist
//   - Any dependency is newer than the target
func (b *Builder) needsRebuild(target string, dependencies []string) bool {
	targetStat, err := os.Stat(target)
	if err != nil {
		// Target doesn't exist, needs rebuild
		return true
	}

	targetTime := targetStat.ModTime()

	// Check if any dependency is newer than the target
	for _, dep := range dependencies {
		depStat, err := os.Stat(dep)
		if err != nil {
			// Dependency doesn't exist as file, skip timestamp check
			continue
		}
		if depStat.ModTime().After(targetTime) {
			return true
		}
	}

	return false
}

// fileExists checks if a file exists on the filesystem.
func (b *Builder) fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// executeCommand executes a shell command and prints it for visibility.
func (b *Builder) executeCommand(command string) error {
	fmt.Printf("\t%s\n", command)
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// executeCommandWithContext executes a shell command with automatic variable expansion.
func (b *Builder) executeCommandWithContext(command string, autoVars *types.AutomaticVariables) error {
	// Expand automatic variables in the command
	expandedCommand := b.makefile.ExpandVariablesWithContext(command, autoVars)
	fmt.Printf("\t%s\n", expandedCommand)
	cmd := exec.Command("sh", "-c", expandedCommand)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// createAutomaticVariables creates automatic variables context for a target.
func (b *Builder) createAutomaticVariables(target string, dependencies []string) *types.AutomaticVariables {
	autoVars := &types.AutomaticVariables{
		Target:      target,
		AllPrereqs:  dependencies,
	}
	
	// Set first prerequisite
	if len(dependencies) > 0 {
		autoVars.FirstPrereq = dependencies[0]
	}
	
	// Determine newer prerequisites ($?)
	autoVars.NewerPrereqs = b.getNewerPrerequisites(target, dependencies)
	
	return autoVars
}

// getNewerPrerequisites returns prerequisites that are newer than the target.
func (b *Builder) getNewerPrerequisites(target string, dependencies []string) []string {
	targetStat, err := os.Stat(target)
	if err != nil {
		// If target doesn't exist, all dependencies are "newer"
		return dependencies
	}
	
	targetTime := targetStat.ModTime()
	var newerDeps []string
	
	for _, dep := range dependencies {
		depStat, err := os.Stat(dep)
		if err != nil {
			// If dependency doesn't exist as file, skip it
			continue
		}
		if depStat.ModTime().After(targetTime) {
			newerDeps = append(newerDeps, dep)
		}
	}
	
	return newerDeps
}