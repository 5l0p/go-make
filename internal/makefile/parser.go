package makefile

import (
	"bufio"
	"os"
	"strings"

	"go-make/pkg/types"
)

// ParseMakefile parses a Makefile from the given filename and returns a Makefile struct
func ParseMakefile(filename string) (*types.Makefile, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	makefile := types.NewMakefile()
	scanner := bufio.NewScanner(file)
	var currentRule *types.Rule

	for scanner.Scan() {
		line := scanner.Text()
		
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		if strings.HasPrefix(line, "\t") {
			if currentRule != nil {
				command := strings.TrimPrefix(line, "\t")
				currentRule.Commands = append(currentRule.Commands, command)
			}
		} else if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			target := strings.TrimSpace(parts[0])
			deps := strings.Fields(strings.TrimSpace(parts[1]))
			
			rule := &types.Rule{
				Target:       target,
				Dependencies: deps,
				Commands:     []string{},
			}
			
			if makefile.FirstRule == "" {
				makefile.FirstRule = target
			}
			
			makefile.Rules[target] = rule
			currentRule = rule
		}
	}

	return makefile, scanner.Err()
}