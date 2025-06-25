package builder

import (
	"os"
	"strings"
	"testing"
	"time"

	"go-make/pkg/types"
)

func TestBuilderBuild(t *testing.T) {
	makefile := &types.Makefile{
		Rules: map[string]*types.Rule{
			"all": {
				Target:       "all",
				Dependencies: []string{"hello.o"},
				Commands:     []string{"echo 'Linking all'"},
			},
			"hello.o": {
				Target:       "hello.o",
				Dependencies: []string{"hello.c"},
				Commands:     []string{"echo 'Compiling hello.c'"},
			},
		},
	}

	tmpdir := t.TempDir()
	oldwd, _ := os.Getwd()
	defer os.Chdir(oldwd)
	os.Chdir(tmpdir)

	os.WriteFile("hello.c", []byte("int main() { return 0; }"), 0644)

	builder := NewBuilder(makefile)
	err := builder.Build("all")
	if err != nil {
		t.Errorf("Build failed: %v", err)
	}

	if !builder.IsBuilt("all") {
		t.Error("Target 'all' should be marked as built")
	}

	if !builder.IsBuilt("hello.o") {
		t.Error("Target 'hello.o' should be marked as built")
	}
}

func TestBuilderNeedsRebuild(t *testing.T) {
	makefile := &types.Makefile{Rules: map[string]*types.Rule{}}
	builder := NewBuilder(makefile)

	tmpdir := t.TempDir()
	oldwd, _ := os.Getwd()
	defer os.Chdir(oldwd)
	os.Chdir(tmpdir)

	sourceFile := "source.txt"
	targetFile := "target.txt"

	os.WriteFile(sourceFile, []byte("source content"), 0644)
	time.Sleep(10 * time.Millisecond)
	os.WriteFile(targetFile, []byte("target content"), 0644)

	if builder.needsRebuild(targetFile, []string{sourceFile}) {
		t.Error("Target should not need rebuild when newer than dependencies")
	}

	time.Sleep(10 * time.Millisecond)
	os.WriteFile(sourceFile, []byte("updated source"), 0644)

	if !builder.needsRebuild(targetFile, []string{sourceFile}) {
		t.Error("Target should need rebuild when dependencies are newer")
	}

	if !builder.needsRebuild("nonexistent", []string{sourceFile}) {
		t.Error("Nonexistent target should always need rebuild")
	}
}

func TestBuilderFileExists(t *testing.T) {
	makefile := &types.Makefile{Rules: map[string]*types.Rule{}}
	builder := NewBuilder(makefile)

	tmpdir := t.TempDir()
	oldwd, _ := os.Getwd()
	defer os.Chdir(oldwd)
	os.Chdir(tmpdir)

	testFile := "test.txt"
	os.WriteFile(testFile, []byte("test"), 0644)

	if !builder.fileExists(testFile) {
		t.Error("fileExists should return true for existing file")
	}

	if builder.fileExists("nonexistent.txt") {
		t.Error("fileExists should return false for nonexistent file")
	}
}

func TestBuilderMissingTarget(t *testing.T) {
	makefile := &types.Makefile{Rules: map[string]*types.Rule{}}
	builder := NewBuilder(makefile)

	tmpdir := t.TempDir()
	oldwd, _ := os.Getwd()
	defer os.Chdir(oldwd)
	os.Chdir(tmpdir)

	err := builder.Build("nonexistent")
	if err == nil {
		t.Error("Build should fail for nonexistent target that doesn't exist as file")
	}

	os.WriteFile("existing_file", []byte("content"), 0644)
	err = builder.Build("existing_file")
	if err != nil {
		t.Errorf("Build should succeed for existing file: %v", err)
	}
}

func TestBuilderCircularDependency(t *testing.T) {
	makefile := &types.Makefile{
		Rules: map[string]*types.Rule{
			"a": {
				Target:       "a",
				Dependencies: []string{"b"},
				Commands:     []string{"echo 'building a'"},
			},
			"b": {
				Target:       "b",
				Dependencies: []string{"a"},
				Commands:     []string{"echo 'building b'"},
			},
		},
	}

	builder := NewBuilder(makefile)
	err := builder.Build("a")
	if err == nil {
		t.Error("Expected circular dependency error, but build succeeded")
	}
	
	expectedError := "circular dependency detected"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', got: %v", expectedError, err)
	}
}

func TestBuilderReset(t *testing.T) {
	makefile := &types.Makefile{
		Rules: map[string]*types.Rule{
			"test": {
				Target:   "test",
				Commands: []string{"echo 'test'"},
			},
		},
	}

	builder := NewBuilder(makefile)
	err := builder.Build("test")
	if err != nil {
		t.Errorf("Build failed: %v", err)
	}

	if !builder.IsBuilt("test") {
		t.Error("Target should be marked as built")
	}

	builder.Reset()

	if builder.IsBuilt("test") {
		t.Error("Target should not be marked as built after reset")
	}
}