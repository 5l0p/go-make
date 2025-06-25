package makefile

import (
	"os"
	"reflect"
	"testing"

	"go-make/pkg/types"
)

func TestParseMakefile(t *testing.T) {
	testMakefile := `# This is a comment
all: hello world
	echo "Building all"
	gcc -o hello hello.c

hello: hello.c
	gcc -c hello.c

world: world.c
	gcc -c world.c

clean:
	rm -f *.o hello

.PHONY: clean all
`

	tmpfile, err := os.CreateTemp("", "test-makefile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(testMakefile)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	makefile, err := ParseMakefile(tmpfile.Name())
	if err != nil {
		t.Fatalf("ParseMakefile failed: %v", err)
	}

	expectedRules := map[string]*types.Rule{
		"all": {
			Target:       "all",
			Dependencies: []string{"hello", "world"},
			Commands:     []string{"echo \"Building all\"", "gcc -o hello hello.c"},
		},
		"hello": {
			Target:       "hello",
			Dependencies: []string{"hello.c"},
			Commands:     []string{"gcc -c hello.c"},
		},
		"world": {
			Target:       "world",
			Dependencies: []string{"world.c"},
			Commands:     []string{"gcc -c world.c"},
		},
		"clean": {
			Target:       "clean",
			Dependencies: []string{},
			Commands:     []string{"rm -f *.o hello"},
		},
		".PHONY": {
			Target:       ".PHONY",
			Dependencies: []string{"clean", "all"},
			Commands:     []string{},
		},
	}

	if len(makefile.Rules) != len(expectedRules) {
		t.Errorf("Expected %d rules, got %d", len(expectedRules), len(makefile.Rules))
	}

	for target, expectedRule := range expectedRules {
		rule, exists := makefile.Rules[target]
		if !exists {
			t.Errorf("Rule %s not found", target)
			continue
		}

		if rule.Target != expectedRule.Target {
			t.Errorf("Rule %s: expected target %s, got %s", target, expectedRule.Target, rule.Target)
		}

		if !reflect.DeepEqual(rule.Dependencies, expectedRule.Dependencies) {
			t.Errorf("Rule %s: expected dependencies %v, got %v", target, expectedRule.Dependencies, rule.Dependencies)
		}

		if !reflect.DeepEqual(rule.Commands, expectedRule.Commands) {
			t.Errorf("Rule %s: expected commands %v, got %v", target, expectedRule.Commands, rule.Commands)
		}
	}
}

func TestParseSimpleMakefile(t *testing.T) {
	testMakefile := `target: dependency
	echo "simple test"
`

	tmpfile, err := os.CreateTemp("", "simple-makefile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(testMakefile)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	makefile, err := ParseMakefile(tmpfile.Name())
	if err != nil {
		t.Fatalf("ParseMakefile failed: %v", err)
	}

	if len(makefile.Rules) != 1 {
		t.Errorf("Expected 1 rule, got %d", len(makefile.Rules))
	}

	rule := makefile.Rules["target"]
	if rule == nil {
		t.Fatal("Rule 'target' not found")
	}

	if rule.Target != "target" {
		t.Errorf("Expected target 'target', got '%s'", rule.Target)
	}

	if len(rule.Dependencies) != 1 || rule.Dependencies[0] != "dependency" {
		t.Errorf("Expected dependencies ['dependency'], got %v", rule.Dependencies)
	}

	if len(rule.Commands) != 1 || rule.Commands[0] != "echo \"simple test\"" {
		t.Errorf("Expected commands ['echo \"simple test\"'], got %v", rule.Commands)
	}
}

func TestParseEmptyMakefile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "empty-makefile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	makefile, err := ParseMakefile(tmpfile.Name())
	if err != nil {
		t.Fatalf("ParseMakefile failed: %v", err)
	}

	if len(makefile.Rules) != 0 {
		t.Errorf("Expected 0 rules for empty makefile, got %d", len(makefile.Rules))
	}
}