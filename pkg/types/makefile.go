// Package types defines the core data structures used throughout the go-make project.
package types

// Rule represents a single target rule in a Makefile.
// A rule consists of a target name, its dependencies, and the commands to build it.
//
// Example Makefile rule:
//   hello: hello.c
//   	gcc -o hello hello.c
//
// This would be represented as:
//   Rule{
//       Target: "hello",
//       Dependencies: []string{"hello.c"},
//       Commands: []string{"gcc -o hello hello.c"},
//   }
type Rule struct {
	// Target is the name of the target being built
	Target string
	
	// Dependencies are the files or targets that this target depends on
	Dependencies []string
	
	// Commands are the shell commands to execute when building this target
	Commands []string
}

// Makefile represents a parsed Makefile with all its rules.
// It contains a map of all rules indexed by target name, tracks
// the first rule encountered (used as the default target), and
// stores variable definitions.
type Makefile struct {
	// Rules maps target names to their corresponding Rule definitions
	Rules map[string]*Rule
	
	// FirstRule is the name of the first target encountered in the Makefile.
	// This is used as the default target when none is specified.
	FirstRule string
	
	// Variables stores variable definitions from the Makefile (VAR = value)
	Variables map[string]string
}

// NewMakefile creates a new empty Makefile with initialized maps.
func NewMakefile() *Makefile {
	return &Makefile{
		Rules:     make(map[string]*Rule),
		Variables: make(map[string]string),
	}
}

// HasTarget returns true if the Makefile contains a rule for the given target.
func (m *Makefile) HasTarget(target string) bool {
	_, exists := m.Rules[target]
	return exists
}

// GetTarget returns the Rule for the given target, or nil if it doesn't exist.
func (m *Makefile) GetTarget(target string) *Rule {
	return m.Rules[target]
}

// Targets returns a slice of all target names in the Makefile.
func (m *Makefile) Targets() []string {
	targets := make([]string, 0, len(m.Rules))
	for target := range m.Rules {
		targets = append(targets, target)
	}
	return targets
}

// SetVariable sets a variable in the Makefile.
func (m *Makefile) SetVariable(name, value string) {
	m.Variables[name] = value
}

// GetVariable returns the value of a variable, or empty string if not found.
func (m *Makefile) GetVariable(name string) string {
	return m.Variables[name]
}

// HasVariable returns true if the variable is defined.
func (m *Makefile) HasVariable(name string) bool {
	_, exists := m.Variables[name]
	return exists
}

// ExpandVariables expands all variable references in the given string.
// Supports both $(VAR) and ${VAR} syntax.
func (m *Makefile) ExpandVariables(text string) string {
	return expandVariables(text, m.Variables)
}

// ExpandVariablesWithContext expands variables including automatic variables.
// Used during command execution when we know the target context.
func (m *Makefile) ExpandVariablesWithContext(text string, autoVars *AutomaticVariables) string {
	return expandVariablesWithContext(text, m.Variables, autoVars)
}