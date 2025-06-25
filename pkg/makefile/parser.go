// Package makefile provides functionality for parsing and working with Makefiles.
package makefile

import (
	"bufio"
	"os"
	"strings"

	"go-make/pkg/types"
)

// ParseMakefile parses a Makefile from the given filename and returns a Makefile struct.
// It supports:
//   - Target definitions with dependencies (target: dep1 dep2)
//   - Commands indented with tabs
//   - Comments (lines starting with #)
//   - Empty lines (ignored)
//
// Example usage:
//   makefile, err := ParseMakefile("Makefile")
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Printf("First target: %s\n", makefile.FirstRule)
func ParseMakefile(filename string) (*types.Makefile, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ParseMakefileFromReader(file)
}

// ParseMakefileFromReader parses a Makefile from an io.Reader.
// This is useful for testing or when the Makefile content comes from a source
// other than a file on disk.
func ParseMakefileFromReader(reader *os.File) (*types.Makefile, error) {
	makefile := types.NewMakefile()
	scanner := bufio.NewScanner(reader)
	var currentRule *types.Rule

	for scanner.Scan() {
		line := scanner.Text()
		
		// Skip empty lines and comments
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		// Commands start with a tab
		if strings.HasPrefix(line, "\t") {
			if currentRule != nil {
				command := strings.TrimPrefix(line, "\t")
				currentRule.Commands = append(currentRule.Commands, command)
			}
		} else if strings.Contains(line, ":") {
			// Target definition: target: dependency1 dependency2
			parts := strings.SplitN(line, ":", 2)
			target := strings.TrimSpace(parts[0])
			deps := strings.Fields(strings.TrimSpace(parts[1]))
			
			rule := &types.Rule{
				Target:       target,
				Dependencies: deps,
				Commands:     []string{},
			}
			
			// Set the first rule as the default target
			if makefile.FirstRule == "" {
				makefile.FirstRule = target
			}
			
			makefile.Rules[target] = rule
			currentRule = rule
		}
	}

	return makefile, scanner.Err()
}