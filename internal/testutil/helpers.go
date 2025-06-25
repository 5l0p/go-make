package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

// CreateTempMakefile creates a temporary file with the given makefile content
func CreateTempMakefile(t *testing.T, content string) string {
	tmpfile, err := os.CreateTemp("", "test-makefile-*")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	return tmpfile.Name()
}

// CreateTempDir creates a temporary directory and changes to it, returning a cleanup function
func CreateTempDir(t *testing.T) (string, func()) {
	tmpdir, err := os.MkdirTemp("", "go-make-test-*")
	if err != nil {
		t.Fatal(err)
	}

	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(tmpdir); err != nil {
		t.Fatal(err)
	}

	cleanup := func() {
		os.Chdir(oldwd)
		os.RemoveAll(tmpdir)
	}

	return tmpdir, cleanup
}

// CreateTestFile creates a test file with the given content in the specified path
func CreateTestFile(t *testing.T, path, content string) {
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}