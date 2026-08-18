package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/nesste/phoenix/internal/surface"
)

type stateEvent struct {
	AfterAct int    `json:"after_act"`
	Path     string `json:"path"`
	Content  string `json:"content"`
}

type stateEventFile struct {
	V      int          `json:"v"`
	Events []stateEvent `json:"events"`
}

type stateEventPlan struct {
	root    string
	events  map[int][]stateEvent
	applied map[int]bool
	mu      sync.Mutex
}

func loadStateEvents(path string) (func(context.Context, int, surface.Input, surface.Envelope) error, error) {
	if path == "" {
		return nil, nil
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read state events: %w", err)
	}
	decoder := json.NewDecoder(bytes.NewReader(contents))
	decoder.DisallowUnknownFields()
	var document stateEventFile
	if err := decoder.Decode(&document); err != nil {
		return nil, fmt.Errorf("decode state events: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return nil, fmt.Errorf("decode state events: trailing JSON data")
	}
	if document.V != 1 {
		return nil, fmt.Errorf("state events version must be 1")
	}
	root, err := filepath.Abs(".")
	if err != nil {
		return nil, fmt.Errorf("resolve state-event root: %w", err)
	}
	plan := &stateEventPlan{root: root, events: make(map[int][]stateEvent), applied: make(map[int]bool)}
	for _, event := range document.Events {
		if event.AfterAct < 0 || strings.TrimSpace(event.Path) == "" {
			return nil, fmt.Errorf("state event has invalid act index or path")
		}
		if _, err := plan.target(event.Path); err != nil {
			return nil, err
		}
		plan.events[event.AfterAct] = append(plan.events[event.AfterAct], event)
	}
	return plan.apply, nil
}

func (plan *stateEventPlan) apply(_ context.Context, index int, _ surface.Input, _ surface.Envelope) error {
	plan.mu.Lock()
	defer plan.mu.Unlock()
	if plan.applied[index] {
		return nil
	}
	for _, event := range plan.events[index] {
		target, err := plan.target(event.Path)
		if err != nil {
			return err
		}
		info, err := os.Lstat(target)
		if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("state event target %q is not a regular file", event.Path)
		}
		if err := os.WriteFile(target, []byte(event.Content), info.Mode().Perm()); err != nil {
			return fmt.Errorf("apply state event %q: %w", event.Path, err)
		}
	}
	plan.applied[index] = true
	return nil
}

func (plan *stateEventPlan) target(relative string) (string, error) {
	if filepath.IsAbs(relative) {
		return "", fmt.Errorf("state event path %q must be relative", relative)
	}
	target := filepath.Join(plan.root, filepath.Clean(filepath.FromSlash(relative)))
	rel, err := filepath.Rel(plan.root, target)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
		return "", fmt.Errorf("state event path %q escapes the sandbox", relative)
	}
	return target, nil
}
