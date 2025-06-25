package builder

import (
	"fmt"
	"os"
	"os/exec"

	"go-make/pkg/types"
)

// Builder handles the build process for targets
type Builder struct {
	makefile *types.Makefile
	built    map[string]bool
	building map[string]bool
}

// NewBuilder creates a new Builder instance
func NewBuilder(makefile *types.Makefile) *Builder {
	return &Builder{
		makefile: makefile,
		built:    make(map[string]bool),
		building: make(map[string]bool),
	}
}

// Build builds the specified target and all its dependencies
func (b *Builder) Build(target string) error {
	if b.built[target] {
		return nil
	}

	if b.building[target] {
		return fmt.Errorf("circular dependency detected involving target '%s'", target)
	}

	rule, exists := b.makefile.Rules[target]
	if !exists {
		if b.fileExists(target) {
			return nil
		}
		return fmt.Errorf("no rule to make target '%s'", target)
	}

	b.building[target] = true

	for _, dep := range rule.Dependencies {
		if err := b.Build(dep); err != nil {
			return err
		}
	}

	if b.needsRebuild(target, rule.Dependencies) {
		fmt.Printf("Building target: %s\n", target)
		for _, command := range rule.Commands {
			if err := b.executeCommand(command); err != nil {
				return fmt.Errorf("command failed: %s", err)
			}
		}
	}

	b.building[target] = false
	b.built[target] = true
	return nil
}

// needsRebuild determines if a target needs to be rebuilt based on dependency timestamps
func (b *Builder) needsRebuild(target string, dependencies []string) bool {
	targetStat, err := os.Stat(target)
	if err != nil {
		return true
	}

	targetTime := targetStat.ModTime()

	for _, dep := range dependencies {
		depStat, err := os.Stat(dep)
		if err != nil {
			continue
		}
		if depStat.ModTime().After(targetTime) {
			return true
		}
	}

	return false
}

// fileExists checks if a file exists on the filesystem
func (b *Builder) fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// executeCommand executes a shell command and prints it
func (b *Builder) executeCommand(command string) error {
	fmt.Printf("\t%s\n", command)
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}