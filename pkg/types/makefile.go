package types

// Rule represents a single target rule in a Makefile
type Rule struct {
	Target       string
	Dependencies []string
	Commands     []string
}

// Makefile represents a parsed Makefile with all its rules
type Makefile struct {
	Rules     map[string]*Rule
	FirstRule string
}

// NewMakefile creates a new empty Makefile
func NewMakefile() *Makefile {
	return &Makefile{
		Rules: make(map[string]*Rule),
	}
}