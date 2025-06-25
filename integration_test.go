package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestIntegrationSimpleExample(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	goMakeBinary := "./go-make"
	if _, err := os.Stat(goMakeBinary); os.IsNotExist(err) {
		cmd := exec.Command("go", "build", "-o", goMakeBinary, "./cmd/go-make")
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to build go-make: %v", err)
		}
	}

	testDir := filepath.Join("examples", "simple")
	oldwd, _ := os.Getwd()
	defer os.Chdir(oldwd)

	if err := os.Chdir(testDir); err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}

	if err := os.Remove("hello"); err != nil && !os.IsNotExist(err) {
		t.Logf("Warning: could not remove existing hello binary: %v", err)
	}

	goMakePath := filepath.Join("..", "..", goMakeBinary)
	cmd := exec.Command(goMakePath, "hello")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go-make failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "Built hello program") {
		t.Errorf("Expected build message in output, got: %s", output)
	}

	if _, err := os.Stat("hello"); os.IsNotExist(err) {
		t.Error("Expected hello binary to be created")
	}

	cmd = exec.Command(goMakePath, "clean")
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go-make clean failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "Cleaned up") {
		t.Errorf("Expected clean message in output, got: %s", output)
	}
}

func TestIntegrationGoExample(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	goMakeBinary := "./go-make"
	if _, err := os.Stat(goMakeBinary); os.IsNotExist(err) {
		cmd := exec.Command("go", "build", "-o", goMakeBinary, "./cmd/go-make")
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to build go-make: %v", err)
		}
	}

	testDir := filepath.Join("examples", "go-project")
	oldwd, _ := os.Getwd()
	defer os.Chdir(oldwd)

	if err := os.Chdir(testDir); err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}

	if err := os.Remove("myapp"); err != nil && !os.IsNotExist(err) {
		t.Logf("Warning: could not remove existing myapp binary: %v", err)
	}

	goMakePath := filepath.Join("..", "..", goMakeBinary)
	cmd := exec.Command(goMakePath, "build")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go-make build failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "Built Go binary") {
		t.Errorf("Expected build message in output, got: %s", output)
	}

	if _, err := os.Stat("myapp"); os.IsNotExist(err) {
		t.Error("Expected myapp binary to be created")
	}

	cmd = exec.Command(goMakePath, "clean")
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go-make clean failed: %v\nOutput: %s", err, output)
	}
}

func TestIntegrationDefaultTarget(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	goMakeBinary := "./go-make"
	testDir := filepath.Join("examples", "simple")
	oldwd, _ := os.Getwd()
	defer os.Chdir(oldwd)

	if err := os.Chdir(testDir); err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}

	if err := os.Remove("hello"); err != nil && !os.IsNotExist(err) {
		t.Logf("Warning: could not remove existing hello binary: %v", err)
	}

	goMakePath := filepath.Join("..", "..", goMakeBinary)
	cmd := exec.Command(goMakePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go-make (default target) failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "Built hello program") {
		t.Errorf("Expected build message in output, got: %s", output)
	}
}