package types

import (
	"os"
	"regexp"
	"strings"
)

// Variable reference patterns: $(VAR), ${VAR}, and automatic variables
var (
	varPattern1 = regexp.MustCompile(`\$\(([^)]+)\)`)  // $(VAR)
	varPattern2 = regexp.MustCompile(`\$\{([^}]+)\}`)  // ${VAR}
	autoVarPattern = regexp.MustCompile(`\$[@<^?]`)    // $@, $<, $^, $?
)

// AutomaticVariables holds the context for automatic variables in a build rule.
type AutomaticVariables struct {
	Target         string   // $@ - the target name
	FirstPrereq    string   // $< - the first prerequisite
	AllPrereqs     []string // $^ - all prerequisites (space-separated)
	NewerPrereqs   []string // $? - prerequisites newer than target
}

// ToString converts automatic variable lists to space-separated strings.
func (av *AutomaticVariables) AllPrereqsString() string {
	return strings.Join(av.AllPrereqs, " ")
}

func (av *AutomaticVariables) NewerPrereqsString() string {
	return strings.Join(av.NewerPrereqs, " ")
}

// expandVariables expands variable references in text using the provided variable map.
// It supports both $(VAR) and ${VAR} syntax and falls back to environment variables.
func expandVariables(text string, variables map[string]string) string {
	return expandVariablesWithContext(text, variables, nil)
}

// expandVariablesWithContext expands variable references including automatic variables.
func expandVariablesWithContext(text string, variables map[string]string, autoVars *AutomaticVariables) string {
	// Replace automatic variables first ($@, $<, $^, $?)
	if autoVars != nil {
		text = autoVarPattern.ReplaceAllStringFunc(text, func(match string) string {
			switch match {
			case "$@":
				return autoVars.Target
			case "$<":
				return autoVars.FirstPrereq
			case "$^":
				return autoVars.AllPrereqsString()
			case "$?":
				return autoVars.NewerPrereqsString()
			default:
				return match // shouldn't happen with our regex
			}
		})
	}

	// Replace $(VAR) patterns
	text = varPattern1.ReplaceAllStringFunc(text, func(match string) string {
		// Extract variable name from $(VAR)
		varName := match[2 : len(match)-1] // Remove $( and )
		return getVariableValue(varName, variables)
	})

	// Replace ${VAR} patterns
	text = varPattern2.ReplaceAllStringFunc(text, func(match string) string {
		// Extract variable name from ${VAR}
		varName := match[2 : len(match)-1] // Remove ${ and }
		return getVariableValue(varName, variables)
	})

	return text
}

// getVariableValue looks up a variable value, first in the provided map,
// then in environment variables.
func getVariableValue(name string, variables map[string]string) string {
	// First check Makefile variables
	if value, exists := variables[name]; exists {
		return value
	}

	// Fall back to environment variables
	return os.Getenv(name)
}

// ParseVariableAssignment parses a variable assignment line like "VAR = value"
// Returns the variable name, value, and whether it was a valid assignment.
func ParseVariableAssignment(line string) (name, value string, isAssignment bool) {
	// Look for = sign (supporting spaces around it)
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", false
	}

	name = strings.TrimSpace(parts[0])
	value = strings.TrimSpace(parts[1])

	// Variable names should be valid identifiers (letters, digits, underscore)
	if name == "" || strings.ContainsAny(name, " \t:") {
		return "", "", false
	}

	return name, value, true
}

// IsVariableAssignment returns true if the line looks like a variable assignment.
func IsVariableAssignment(line string) bool {
	_, _, isAssignment := ParseVariableAssignment(line)
	return isAssignment
}