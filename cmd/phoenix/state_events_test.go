package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/nesste/phoenix/internal/surface"
)

func TestStateEventsApplyAtConfiguredExecutableAct(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	target := filepath.Join(root, "state.txt")
	if err := os.WriteFile(target, []byte("before\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	planPath := filepath.Join(t.TempDir(), "events.json")
	if err := os.WriteFile(planPath, []byte(`{"v":1,"events":[{"after_act":0,"path":"state.txt","content":"after\n"}]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	hook, err := loadStateEvents(planPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := hook(context.Background(), 0, surface.Input{}, surface.Envelope{}); err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != "after\n" {
		t.Fatalf("state contents = %q, want applied event", contents)
	}
}

func TestStateEventsRejectPathsOutsideSandbox(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	planPath := filepath.Join(t.TempDir(), "events.json")
	if err := os.WriteFile(planPath, []byte(`{"v":1,"events":[{"after_act":0,"path":"../outside.txt","content":"x"}]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := loadStateEvents(planPath); err == nil {
		t.Fatal("loadStateEvents accepted an escaping path")
	}
}
