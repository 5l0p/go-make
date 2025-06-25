package types

import (
	"reflect"
	"testing"
)

func TestNewMakefile(t *testing.T) {
	mf := NewMakefile()
	
	if mf.Rules == nil {
		t.Error("NewMakefile should initialize Rules map")
	}
	
	if len(mf.Rules) != 0 {
		t.Error("NewMakefile should create empty Rules map")
	}
	
	if mf.FirstRule != "" {
		t.Error("NewMakefile should have empty FirstRule")
	}
}

func TestMakefileHasTarget(t *testing.T) {
	mf := NewMakefile()
	mf.Rules["test"] = &Rule{Target: "test"}
	
	if !mf.HasTarget("test") {
		t.Error("HasTarget should return true for existing target")
	}
	
	if mf.HasTarget("nonexistent") {
		t.Error("HasTarget should return false for non-existing target")
	}
}

func TestMakefileGetTarget(t *testing.T) {
	mf := NewMakefile()
	rule := &Rule{Target: "test", Commands: []string{"echo test"}}
	mf.Rules["test"] = rule
	
	result := mf.GetTarget("test")
	if result != rule {
		t.Error("GetTarget should return the correct rule")
	}
	
	result = mf.GetTarget("nonexistent")
	if result != nil {
		t.Error("GetTarget should return nil for non-existing target")
	}
}

func TestMakefileTargets(t *testing.T) {
	mf := NewMakefile()
	mf.Rules["a"] = &Rule{Target: "a"}
	mf.Rules["b"] = &Rule{Target: "b"}
	mf.Rules["c"] = &Rule{Target: "c"}
	
	targets := mf.Targets()
	expected := []string{"a", "b", "c"}
	
	if len(targets) != len(expected) {
		t.Errorf("Expected %d targets, got %d", len(expected), len(targets))
	}
	
	// Convert to map for easier comparison (order doesn't matter)
	targetMap := make(map[string]bool)
	for _, target := range targets {
		targetMap[target] = true
	}
	
	for _, expectedTarget := range expected {
		if !targetMap[expectedTarget] {
			t.Errorf("Expected target %s not found in result", expectedTarget)
		}
	}
}

func TestRule(t *testing.T) {
	rule := &Rule{
		Target:       "hello",
		Dependencies: []string{"hello.c"},
		Commands:     []string{"gcc -o hello hello.c"},
	}
	
	if rule.Target != "hello" {
		t.Errorf("Expected target 'hello', got '%s'", rule.Target)
	}
	
	expectedDeps := []string{"hello.c"}
	if !reflect.DeepEqual(rule.Dependencies, expectedDeps) {
		t.Errorf("Expected dependencies %v, got %v", expectedDeps, rule.Dependencies)
	}
	
	expectedCommands := []string{"gcc -o hello hello.c"}
	if !reflect.DeepEqual(rule.Commands, expectedCommands) {
		t.Errorf("Expected commands %v, got %v", expectedCommands, rule.Commands)
	}
}