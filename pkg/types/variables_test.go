package types

import (
	"os"
	"testing"
)

func TestExpandVariables(t *testing.T) {
	variables := map[string]string{
		"CC":     "gcc",
		"CFLAGS": "-Wall -O2",
		"TARGET": "hello",
	}

	tests := []struct {
		input    string
		expected string
		name     string
	}{
		{
			input:    "$(CC) $(CFLAGS) -o $(TARGET)",
			expected: "gcc -Wall -O2 -o hello",
			name:     "basic substitution with $()",
		},
		{
			input:    "${CC} ${CFLAGS} -o ${TARGET}",
			expected: "gcc -Wall -O2 -o hello",
			name:     "basic substitution with {}",
		},
		{
			input:    "$(CC) $(UNKNOWN) $(TARGET)",
			expected: "gcc  hello",
			name:     "unknown variable becomes empty",
		},
		{
			input:    "no variables here",
			expected: "no variables here",
			name:     "no variables to expand",
		},
		{
			input:    "$(TARGET).o depends on $(TARGET).c",
			expected: "hello.o depends on hello.c",
			name:     "multiple same variable",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := expandVariables(test.input, variables)
			if result != test.expected {
				t.Errorf("expandVariables(%q) = %q, want %q", test.input, result, test.expected)
			}
		})
	}
}

func TestExpandVariablesWithEnvironment(t *testing.T) {
	// Set an environment variable for testing
	os.Setenv("TEST_ENV_VAR", "env_value")
	defer os.Unsetenv("TEST_ENV_VAR")

	variables := map[string]string{
		"MAKEFILE_VAR": "makefile_value",
	}

	tests := []struct {
		input    string
		expected string
		name     string
	}{
		{
			input:    "$(MAKEFILE_VAR)",
			expected: "makefile_value",
			name:     "makefile variable takes precedence",
		},
		{
			input:    "$(TEST_ENV_VAR)",
			expected: "env_value",
			name:     "environment variable used when not in makefile",
		},
		{
			input:    "$(USER)",
			expected: os.Getenv("USER"),
			name:     "real environment variable",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := expandVariables(test.input, variables)
			if result != test.expected {
				t.Errorf("expandVariables(%q) = %q, want %q", test.input, result, test.expected)
			}
		})
	}
}

func TestParseVariableAssignment(t *testing.T) {
	tests := []struct {
		input        string
		expectedName string
		expectedValue string
		expectedValid bool
		name         string
	}{
		{
			input:        "CC = gcc",
			expectedName: "CC",
			expectedValue: "gcc",
			expectedValid: true,
			name:         "simple assignment",
		},
		{
			input:        "CFLAGS=-Wall -O2",
			expectedName: "CFLAGS",
			expectedValue: "-Wall -O2",
			expectedValid: true,
			name:         "assignment without spaces",
		},
		{
			input:        "   VAR   =   value   ",
			expectedName: "VAR",
			expectedValue: "value",
			expectedValid: true,
			name:         "assignment with extra spaces",
		},
		{
			input:        "target: dependency",
			expectedName: "",
			expectedValue: "",
			expectedValid: false,
			name:         "not an assignment (target rule)",
		},
		{
			input:        "just some text",
			expectedName: "",
			expectedValue: "",
			expectedValid: false,
			name:         "not an assignment (plain text)",
		},
		{
			input:        "EMPTY =",
			expectedName: "EMPTY",
			expectedValue: "",
			expectedValid: true,
			name:         "empty value assignment",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			name, value, isValid := ParseVariableAssignment(test.input)
			
			if isValid != test.expectedValid {
				t.Errorf("ParseVariableAssignment(%q) validity = %v, want %v", 
					test.input, isValid, test.expectedValid)
			}
			
			if isValid {
				if name != test.expectedName {
					t.Errorf("ParseVariableAssignment(%q) name = %q, want %q", 
						test.input, name, test.expectedName)
				}
				if value != test.expectedValue {
					t.Errorf("ParseVariableAssignment(%q) value = %q, want %q", 
						test.input, value, test.expectedValue)
				}
			}
		})
	}
}

func TestMakefileVariableMethods(t *testing.T) {
	mf := NewMakefile()
	
	// Test SetVariable and GetVariable
	mf.SetVariable("TEST_VAR", "test_value")
	
	if !mf.HasVariable("TEST_VAR") {
		t.Error("HasVariable should return true for set variable")
	}
	
	if value := mf.GetVariable("TEST_VAR"); value != "test_value" {
		t.Errorf("GetVariable returned %q, want %q", value, "test_value")
	}
	
	if mf.HasVariable("NONEXISTENT") {
		t.Error("HasVariable should return false for unset variable")
	}
	
	if value := mf.GetVariable("NONEXISTENT"); value != "" {
		t.Errorf("GetVariable for nonexistent var returned %q, want empty string", value)
	}
	
	// Test ExpandVariables
	mf.SetVariable("CC", "gcc")
	mf.SetVariable("FLAGS", "-Wall")
	
	expanded := mf.ExpandVariables("$(CC) $(FLAGS) -o target")
	expected := "gcc -Wall -o target"
	if expanded != expected {
		t.Errorf("ExpandVariables returned %q, want %q", expanded, expected)
	}
}