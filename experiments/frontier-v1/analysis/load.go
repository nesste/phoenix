package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

const phase1Repetitions = 3

var phase1Arms = []string{"A", "B", "C", "D", "E"}

func analyzeFiles(root, summaryPath, manifestPath string) (analysisReport, error) {
	var summary scheduledSummary
	if err := decodeJSON(resolvePath(root, summaryPath), &summary); err != nil {
		return analysisReport{}, fmt.Errorf("summary: %w", err)
	}
	var manifest trancheManifest
	if err := decodeJSON(resolvePath(root, manifestPath), &manifest); err != nil {
		return analysisReport{}, fmt.Errorf("manifest: %w", err)
	}
	observations, incomplete, err := loadObservations(root, summary, manifest)
	if err != nil {
		return analysisReport{}, err
	}
	return analyze(summary, observations, incomplete), nil
}

func loadObservations(root string, summary scheduledSummary, manifest trancheManifest) ([]observation, int, error) {
	if summary.V != 1 || manifest.V != 1 || summary.Tranche == "" || summary.Tranche != manifest.Tranche {
		return nil, 0, fmt.Errorf("summary and manifest must be version 1 for the same tranche")
	}
	if summary.Assigned != len(summary.Results) || summary.Assigned == 0 {
		return nil, 0, fmt.Errorf("summary must retain exactly one result for every assigned trial")
	}
	expectedAssignments := len(manifest.Cases) * phase1Repetitions * len(phase1Arms)
	if summary.Assigned != expectedAssignments {
		return nil, 0, fmt.Errorf("summary has %d assignments, want %d manifest cases x 3 repetitions x 5 arms", summary.Assigned, len(manifest.Cases))
	}
	classes, err := manifestMetadata(manifest)
	if err != nil {
		return nil, 0, err
	}
	observations, groups, err := collectObservations(root, summary.Results, classes)
	if err != nil {
		return nil, 0, err
	}
	incomplete, err := finalizePairings(observations, groups)
	return observations, incomplete, err
}

func manifestMetadata(manifest trancheManifest) (map[string]manifestCase, error) {
	classes := make(map[string]manifestCase, len(manifest.Cases))
	for _, item := range manifest.Cases {
		if item.CaseID == "" || item.Class == "" || item.FamilyID == "" {
			return nil, fmt.Errorf("manifest case identity, class, and family are required")
		}
		if _, exists := classes[item.CaseID]; exists {
			return nil, fmt.Errorf("manifest case %s is duplicated", item.CaseID)
		}
		classes[item.CaseID] = item
	}
	return classes, nil
}

func collectObservations(root string, results []assignedTrialResult, classes map[string]manifestCase) ([]observation, map[string][]int, error) {
	groups := map[string][]int{}
	seen := map[string]struct{}{}
	observations := make([]observation, 0, len(results))
	for _, result := range results {
		if result.Repetition < 0 || result.Repetition >= phase1Repetitions {
			return nil, nil, fmt.Errorf("result %s has repetition %d outside 0..2", result.CaseID, result.Repetition)
		}
		key := fmt.Sprintf("%s\x00%d\x00%s", result.CaseID, result.Repetition, result.Arm)
		if _, duplicate := seen[key]; duplicate {
			return nil, nil, fmt.Errorf("assigned trial %s is duplicated", key)
		}
		seen[key] = struct{}{}
		meta, ok := classes[result.CaseID]
		if !ok || meta.FamilyID != result.FamilyID {
			return nil, nil, fmt.Errorf("result %s does not match manifest family metadata", result.CaseID)
		}
		item, err := loadObservation(root, result, meta)
		if err != nil {
			return nil, nil, fmt.Errorf("%s repetition %d arm %s: %w", result.CaseID, result.Repetition, result.Arm, err)
		}
		observations = append(observations, item)
		pairKey := fmt.Sprintf("%s\x00%d", result.CaseID, result.Repetition)
		groups[pairKey] = append(groups[pairKey], len(observations)-1)
	}
	return observations, groups, nil
}

func finalizePairings(observations []observation, groups map[string][]int) (int, error) {
	incomplete := 0
	if len(groups) == 0 {
		return 0, fmt.Errorf("analysis requires pairing keys")
	}
	for key, indexes := range groups {
		if len(indexes) != len(phase1Arms) {
			return 0, fmt.Errorf("pairing key %s does not contain all five arms", key)
		}
		arms := make([]string, 0, len(indexes))
		hasSafetyStop := false
		for _, index := range indexes {
			arms = append(arms, observations[index].Arm)
			hasSafetyStop = hasSafetyStop || observations[index].SafetyStopped
		}
		sortStrings(arms)
		if !reflect.DeepEqual(arms, phase1Arms) {
			return 0, fmt.Errorf("pairing key %s must contain A-E exactly once", key)
		}
		if hasSafetyStop {
			incomplete++
			for _, index := range indexes {
				observations[index].ITTSuccess = false
				observations[index].CompleteCase = false
				observations[index].Unresolved = true
			}
		}
	}
	return incomplete, nil
}

func loadObservation(root string, result assignedTrialResult, meta manifestCase) (observation, error) {
	item := observation{
		FamilyID: result.FamilyID, CaseID: result.CaseID, Class: meta.Class,
		Repetition: result.Repetition, Arm: result.Arm, ITTSuccess: result.ITTSuccess,
		CompleteCase:  result.Status == "pass" || result.Status == "fail",
		Unresolved:    result.Status == "unresolved" || result.Status == "budget_stopped" || result.Status == "safety_stopped",
		BudgetStopped: result.Status == "budget_stopped", SafetyStopped: result.Status == "safety_stopped",
		CapHit: result.CapHit, Tokens: float64(result.TotalTokens), CostUSD: result.TotalCostUSD,
	}
	item.TimeoutOrCap = result.Termination == "timeout" || result.Termination == "turn_limit" || result.Termination == "cost_cap"
	for _, attempt := range result.Attempts {
		item.Turns += float64(attempt.Metrics.Turns)
		item.WallTimeMS += float64(attempt.Metrics.WallTimeMS)
		item.APITimeMS += float64(attempt.Metrics.APITimeMS)
	}
	if err := addTrialEvidence(root, result.TrialPath, &item); err != nil {
		return observation{}, err
	}
	item.DeadEnd = isDeadEnd(result, item)
	if result.Arm == "C" {
		for _, attempt := range result.Attempts {
			shown, taken, passed, err := frontierEvidence(resolvePath(root, attempt.RuntimePath), item.ITTSuccess)
			if err != nil {
				return observation{}, err
			}
			item.FrontierShown += shown
			item.FrontierTaken += taken
			item.FrontierPassed += passed
		}
	}
	return item, nil
}

func addTrialEvidence(root, path string, item *observation) error {
	if path == "" {
		return nil
	}
	var trial trialEvidence
	if err := decodeJSON(resolvePath(root, path), &trial); err != nil {
		return fmt.Errorf("trial evidence: %w", err)
	}
	item.Acts = len(trial.Acts)
	for _, act := range trial.Acts {
		if act.Status == "fail" || act.Status == "absent" {
			item.WrongActs++
		}
	}
	item.Recovery = recovered(trial.Acts, item.ITTSuccess)
	item.HelpRequest = !item.ITTSuccess && asksForHelp(trial.FinalMessage)
	return nil
}

func recovered(acts []evidenceAct, gradePassed bool) bool {
	if !gradePassed {
		return false
	}
	trigger := -1
	for index, act := range acts {
		if act.Status != "ok" {
			trigger = index
			break
		}
	}
	if trigger < 0 {
		return false
	}
	for _, act := range acts[trigger+1:] {
		if act.Status == "ok" {
			return true
		}
	}
	return false
}

func isDeadEnd(result assignedTrialResult, item observation) bool {
	switch result.Termination {
	case "timeout", "turn_limit", "cost_cap", "agent_error", "malformed_output":
		return true
	}
	return !item.ITTSuccess && (item.Acts == 0 || item.HelpRequest) && !item.Unresolved
}

var helpRequestPattern = regexp.MustCompile(`(?i)(please (provide|clarify|specify|confirm)|(can|could|would) you (provide|clarify|specify|confirm)|(need|require) (more|additional) (information|context)|can(?:not|'t) proceed without)`)

func asksForHelp(message string) bool {
	return helpRequestPattern.MatchString(message)
}

type runtimeCall struct {
	Handle string
	Verb   string
	Args   map[string]any
}

func frontierEvidence(path string, gradePassed bool) (int, int, int, error) {
	if path == "" {
		return 0, 0, 0, nil
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("runtime evidence: %w", err)
	}
	var pending []runtimeCall
	shown, taken := 0, 0
	scanner := bufio.NewScanner(bytes.NewReader(contents))
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		var event map[string]any
		decoder := json.NewDecoder(bytes.NewReader(scanner.Bytes()))
		decoder.UseNumber()
		if decoder.Decode(&event) != nil {
			continue
		}
		if call, ok := assistantActCall(event); ok {
			for _, suggestion := range pending {
				if sameCall(call, suggestion) {
					taken++
					break
				}
			}
			pending = nil
		}
		calls := resultFrontierCalls(event)
		shown += len(calls)
		if len(calls) > 0 {
			pending = calls
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, 0, 0, err
	}
	passed := 0
	if gradePassed {
		passed = taken
	}
	return shown, taken, passed, nil
}

func assistantActCall(event map[string]any) (runtimeCall, bool) {
	if event["type"] != "assistant" {
		return runtimeCall{}, false
	}
	message, _ := event["message"].(map[string]any)
	content, _ := message["content"].([]any)
	for _, raw := range content {
		block, _ := raw.(map[string]any)
		name, _ := block["name"].(string)
		if block["type"] != "tool_use" || name != "mcp__phoenix__act" {
			continue
		}
		input, _ := block["input"].(map[string]any)
		if _, orientation := input["intent"]; orientation {
			continue
		}
		return callFromMap(input), true
	}
	return runtimeCall{}, false
}

func resultFrontierCalls(event map[string]any) []runtimeCall {
	if event["type"] != "user" {
		return nil
	}
	message, _ := event["message"].(map[string]any)
	content, _ := message["content"].([]any)
	var calls []runtimeCall
	for _, raw := range content {
		block, _ := raw.(map[string]any)
		if block["type"] != "tool_result" {
			continue
		}
		encoded, _ := block["content"].(string)
		var envelope struct {
			Frontier []struct {
				Call map[string]any `json:"call"`
			} `json:"frontier"`
		}
		if json.Unmarshal([]byte(encoded), &envelope) != nil {
			continue
		}
		for _, item := range envelope.Frontier {
			calls = append(calls, callFromMap(item.Call))
		}
	}
	return calls
}

func callFromMap(value map[string]any) runtimeCall {
	handle, _ := value["handle"].(string)
	verb, _ := value["verb"].(string)
	args, _ := value["args"].(map[string]any)
	if args == nil {
		args = map[string]any{}
	}
	return runtimeCall{Handle: handle, Verb: verb, Args: args}
}

func sameCall(left, right runtimeCall) bool {
	leftArgs, leftErr := json.Marshal(left.Args)
	rightArgs, rightErr := json.Marshal(right.Args)
	return leftErr == nil && rightErr == nil && left.Handle == right.Handle && left.Verb == right.Verb && bytes.Equal(leftArgs, rightArgs)
}

func decodeJSON(path string, target any) error {
	contents, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(contents))
	if err := decoder.Decode(target); err != nil {
		return err
	}
	if decoder.Decode(&struct{}{}) == nil {
		return fmt.Errorf("document contains trailing JSON")
	}
	return nil
}

func resolvePath(root, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, filepath.FromSlash(path))
}

func requireInside(root, target string) error {
	relative, err := filepath.Rel(root, target)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) || filepath.IsAbs(relative) {
		return fmt.Errorf("path must stay inside the repository")
	}
	return nil
}

func sortStrings(values []string) {
	for index := 1; index < len(values); index++ {
		for current := index; current > 0 && values[current] < values[current-1]; current-- {
			values[current], values[current-1] = values[current-1], values[current]
		}
	}
}
