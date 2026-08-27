package corpus

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Label is the full outcome-bearing content of one case's label, as needed
// to grade a Trial deterministically. corpus.go's unexported `label` type
// stays a minimal projection used only for manifest counting; this type
// carries every field the grader reads.
type Label struct {
	CaseID          string `json:"case_id"`
	Class           string `json:"class"`
	ExpectedOutcome struct {
		Checks []Check `json:"checks"`
	} `json:"expected_outcome"`
	GradingScript   string     `json:"grading_script"`
	AcceptablePaths [][]string `json:"acceptable_paths"`
	LabelRationale  string     `json:"label_rationale"`
}

// Check is one expected-outcome check. Only the fields relevant to Kind are
// populated; the shape mirrors label.schema.json's check oneOf.
type Check struct {
	ID               string   `json:"id"`
	Kind             string   `json:"kind"`
	Path             string   `json:"path,omitempty"`
	Pattern          string   `json:"pattern,omitempty"`
	Command          []string `json:"command,omitempty"`
	ExitCode         *int     `json:"exit_code,omitempty"`
	Negate           bool     `json:"negate,omitempty"`
	Claim            string   `json:"claim,omitempty"`
	Mode             string   `json:"mode,omitempty"`
	Seq              *int     `json:"seq,omitempty"`
	SelectorMode     string   `json:"selector_mode,omitempty"`
	RecencyRationale string   `json:"recency_rationale,omitempty"`
	Status           string   `json:"status,omitempty"`
	Maximum          *int     `json:"maximum,omitempty"`
}

// outcomePrimaryClasses are the classes where the outcome, not the route, is
// the construct: act_sequence, act_count, and act_path_absent are descriptive
// there, and gating act-addressed checks must be selector-addressed. In the
// path-primary classes (absence, temptation, stale_frontier,
// adversarial_text) the path is the outcome and every check kind gates.
var outcomePrimaryClasses = map[string]bool{
	"direct": true, "cascade": true, "far_discovery": true, "recovery": true,
}

// absenceLabelKinds are the only check kinds that appear in absence labels;
// per the frozen check-kind classification, kinds not listed for a class do
// not appear in its labels.
var absenceLabelKinds = map[string]bool{
	"final_message_matches": true, "act_sequence": true, "act_count": true, "act_path_absent": true,
}

// checkGates reports whether one check kind gates the overall status for a
// class. Descriptive checks are still graded and reported; they cannot fail
// the trial.
func checkGates(class, kind string) bool {
	if outcomePrimaryClasses[class] {
		switch kind {
		case "act_sequence", "act_count", "act_path_absent":
			return false
		}
	}
	return true
}

const (
	VerdictPass           = "pass"
	VerdictFail           = "fail"
	VerdictManualRequired = "manual_required"

	StatusPass          = "pass"
	StatusFail          = "fail"
	StatusIndeterminate = "indeterminate"
)

// CheckResult is the deterministic verdict for one check (or the
// acceptable_paths match, recorded under kind "acceptable_paths"). Gating
// records whether the verdict participates in the overall status under the
// frozen per-class check-kind classification.
type CheckResult struct {
	ID      string `json:"id"`
	Kind    string `json:"kind"`
	Verdict string `json:"verdict"`
	Reason  string `json:"reason"`
	Gating  bool   `json:"gating"`
}

// GradeResult is the complete, deterministic grading of one trial against
// one label.
type GradeResult struct {
	CaseID string        `json:"case_id"`
	Checks []CheckResult `json:"checks"`
	Status string        `json:"status"`
}

// LoadLabel reads, canonicalizes, schema-validates, and decodes a label
// document at root-relative path.
func LoadLabel(root, path string) (Document, Label, error) {
	schemas, err := loadSchemas(root)
	if err != nil {
		return Document{}, Label{}, err
	}
	document, err := loadDocument(root, path)
	if err != nil {
		return Document{}, Label{}, err
	}
	if err := schemas.validate("label", document); err != nil {
		return Document{}, Label{}, err
	}
	var decoded Label
	if err := document.Into(&decoded); err != nil {
		return Document{}, Label{}, err
	}
	digest, err := GraderDigest(root)
	if err != nil {
		return Document{}, Label{}, err
	}
	if decoded.GradingScript != digest {
		return Document{}, Label{}, fmt.Errorf("%s: grading_script %s does not match current grader %s", document.Path, decoded.GradingScript, digest)
	}
	return document, decoded, nil
}

// Grade deterministically compares a trial's captured evidence against a
// label's expected outcome. It never invokes an LLM, never shell-evaluates
// a command, and never accesses the filesystem: every verdict is computed
// from the trial's acts, final_message, and end_state alone.
func Grade(label Label, trial Trial) (GradeResult, error) {
	if label.CaseID != trial.CaseID {
		return GradeResult{}, fmt.Errorf("label case_id %s does not match trial case_id %s", label.CaseID, trial.CaseID)
	}
	for _, check := range label.ExpectedOutcome.Checks {
		if check.Pattern != "" {
			if _, err := regexp.Compile(check.Pattern); err != nil {
				return GradeResult{}, fmt.Errorf("check %s has invalid pattern %q: %w", check.ID, check.Pattern, err)
			}
		}
		if err := validateCheckForClass(label.Class, check); err != nil {
			return GradeResult{}, err
		}
	}
	result := GradeResult{CaseID: label.CaseID}
	path := matchLabelPath(label, trial.Acts)
	for _, check := range label.ExpectedOutcome.Checks {
		graded := gradeCheck(check, label.AcceptablePaths, path, trial)
		graded.Gating = checkGates(label.Class, check.Kind)
		result.Checks = append(result.Checks, graded)
	}
	result.Status = overallStatus(result.Checks)
	return result, nil
}

// validateCheckForClass enforces the frozen per-class addressing rules: an
// absence label may carry only its listed kinds, and gating act-addressed
// checks in outcome-primary classes are selector-addressed, never
// sequence-index-addressed.
func validateCheckForClass(class string, check Check) error {
	if class == "absence" && !absenceLabelKinds[check.Kind] {
		return fmt.Errorf("check %s: kind %s does not appear in absence labels", check.ID, check.Kind)
	}
	if check.Kind != "act_status" && check.Kind != "act_output_matches" {
		return nil
	}
	hasSequence := check.Seq != nil
	hasSelector := check.SelectorMode != ""
	if hasSequence == hasSelector {
		return fmt.Errorf("check %s must use exactly one of seq or selector_mode", check.ID)
	}
	if hasSelector {
		if check.SelectorMode != "some" && check.SelectorMode != "last" {
			return fmt.Errorf("check %s has unknown selector_mode %q", check.ID, check.SelectorMode)
		}
		if check.Path == "" {
			return fmt.Errorf("check %s selector requires a handle-type.verb path predicate", check.ID)
		}
		if check.SelectorMode == "last" && check.RecencyRationale == "" {
			return fmt.Errorf("check %s uses last-mode selection without a recency_rationale", check.ID)
		}
		if check.SelectorMode == "some" && check.RecencyRationale != "" {
			return fmt.Errorf("check %s carries a recency_rationale without last-mode selection", check.ID)
		}
	}
	if outcomePrimaryClasses[class] && hasSequence {
		return fmt.Errorf("check %s: gating act-addressed checks in class %s must be selector-addressed, not sequence-indexed", check.ID, class)
	}
	return nil
}

func overallStatus(checks []CheckResult) string {
	sawManual := false
	for _, check := range checks {
		if !check.Gating {
			continue
		}
		switch check.Verdict {
		case VerdictFail:
			return StatusFail
		case VerdictManualRequired:
			sawManual = true
		}
	}
	if sawManual {
		return StatusIndeterminate
	}
	return StatusPass
}

func gradePaths(acceptable [][]string, acts []Act, mode string) CheckResult {
	result := CheckResult{ID: "act_sequence", Kind: "act_sequence"}
	observed := make([]string, 0, len(acts))
	for _, act := range acts {
		observed = append(observed, strings.ToLower(act.HandleType+"."+act.Verb))
	}
	for _, path := range acceptable {
		if _, matched := sequenceMatchOffsets(observed, path, mode); matched {
			result.Verdict = VerdictPass
			result.Reason = fmt.Sprintf("trial act sequence satisfies an acceptable path in %s mode", mode)
			return result
		}
	}
	result.Verdict = VerdictFail
	result.Reason = fmt.Sprintf("trial act sequence %v satisfies no acceptable path in %s mode", observed, mode)
	return result
}

func sequenceMatchOffsets(observed, expected []string, mode string) ([]int, bool) {
	switch mode {
	case "exact":
		if !exactFoldedSequence(observed, expected) {
			return nil, false
		}
		offsets := make([]int, len(expected))
		for index := range expected {
			offsets[index] = index
		}
		return offsets, true
	case "contains_in_order":
		offsets := make([]int, 0, len(expected))
		for index, verb := range observed {
			position := len(offsets)
			if position < len(expected) && strings.EqualFold(verb, expected[position]) {
				offsets = append(offsets, index)
			}
		}
		return offsets, len(offsets) == len(expected)
	default:
		return nil, false
	}
}

type matchedPath struct {
	offsets []int
	found   bool
}

func matchLabelPath(label Label, acts []Act) matchedPath {
	mode := ""
	for _, check := range label.ExpectedOutcome.Checks {
		if check.Kind == "act_sequence" {
			mode = check.Mode
			break
		}
	}
	if mode == "" {
		return matchedPath{}
	}
	observed := make([]string, 0, len(acts))
	for _, act := range acts {
		observed = append(observed, strings.ToLower(act.HandleType+"."+act.Verb))
	}
	for _, acceptable := range label.AcceptablePaths {
		if offsets, found := sequenceMatchOffsets(observed, acceptable, mode); found {
			return matchedPath{offsets: offsets, found: true}
		}
	}
	return matchedPath{}
}

func exactFoldedSequence(observed, expected []string) bool {
	if len(observed) != len(expected) {
		return false
	}
	for index, verb := range expected {
		if !strings.EqualFold(observed[index], verb) {
			return false
		}
	}
	return true
}

func gradeCheck(check Check, acceptablePaths [][]string, path matchedPath, trial Trial) CheckResult {
	result := CheckResult{ID: check.ID, Kind: check.Kind}
	switch check.Kind {
	case "act_sequence":
		result = gradePaths(acceptablePaths, trial.Acts, check.Mode)
		result.ID = check.ID
	case "act_count":
		result.Verdict, result.Reason = gradeActCount(check, trial)
	case "act_path_absent":
		result.Verdict, result.Reason = gradeActPathAbsent(check, trial)
	case "act_status":
		if check.SelectorMode != "" {
			result.Verdict, result.Reason = gradeActStatusSelected(check, trial)
		} else {
			result.Verdict, result.Reason = gradeActStatusAtPath(check, path, trial)
		}
	case "act_output_matches":
		if check.SelectorMode != "" {
			result.Verdict, result.Reason = gradeActOutputSelected(check, trial)
		} else {
			result.Verdict, result.Reason = gradeActOutputAtPath(check, path, trial)
		}
	case "file_matches":
		result.Verdict, result.Reason = gradeFileMatches(check, trial)
	case "file_absent":
		result.Verdict, result.Reason = gradeFileAbsent(check, trial)
	case "command_exit":
		result.Verdict, result.Reason = gradeCommandExit(check, trial)
	case "command_output_matches":
		result.Verdict, result.Reason = gradeCommandOutputMatches(check, trial)
	case "final_message_matches":
		result.Verdict, result.Reason = gradeFinalMessageMatches(check, trial)
	case "final_message_states":
		result.Verdict = VerdictFail
		result.Reason = "final_message_states is barred from machine-graded labels; no sealed label may carry a manually graded kind"
	default:
		result.Verdict = VerdictFail
		result.Reason = fmt.Sprintf("unknown check kind %q", check.Kind)
	}
	return result
}

func gradeActCount(check Check, trial Trial) (string, string) {
	if check.Maximum == nil {
		return VerdictFail, "act_count check has no maximum"
	}
	if len(trial.Acts) > *check.Maximum {
		return VerdictFail, fmt.Sprintf("trial recorded %d acts, maximum is %d", len(trial.Acts), *check.Maximum)
	}
	return VerdictPass, fmt.Sprintf("trial recorded %d acts, within maximum %d", len(trial.Acts), *check.Maximum)
}

func gradeActPathAbsent(check Check, trial Trial) (string, string) {
	for _, act := range trial.Acts {
		path := strings.ToLower(act.HandleType + "." + act.Verb)
		forbidden := strings.ToLower(check.Path)
		if path == forbidden || (!strings.Contains(forbidden, ".") && strings.HasPrefix(path, forbidden+".")) {
			return VerdictFail, fmt.Sprintf("trial invoked forbidden act path %s", path)
		}
	}
	return VerdictPass, fmt.Sprintf("trial did not invoke forbidden act path %s", check.Path)
}

// selectedActIndexes returns the trial act indexes whose handle-type.verb
// path equals the selector predicate, case-folded like acceptable-path verbs.
func selectedActIndexes(trial Trial, path string) []int {
	target := strings.ToLower(path)
	indexes := []int{}
	for index, act := range trial.Acts {
		if strings.ToLower(act.HandleType+"."+act.Verb) == target {
			indexes = append(indexes, index)
		}
	}
	return indexes
}

// gradeActStatusSelected implements selector addressing for act_status. The
// default "some" form is existential: some act matching the predicate has the
// expected status. The "last" form binds the last matching act and exists
// only where recency is the construct under test.
func gradeActStatusSelected(check Check, trial Trial) (string, string) {
	indexes := selectedActIndexes(trial, check.Path)
	if len(indexes) == 0 {
		return VerdictFail, fmt.Sprintf("no act matches selector path %s", check.Path)
	}
	if check.SelectorMode == "last" {
		index := indexes[len(indexes)-1]
		if trial.Acts[index].Status != check.Status {
			return VerdictFail, fmt.Sprintf("last act matching %s has status %s, expected %s", check.Path, trial.Acts[index].Status, check.Status)
		}
		return VerdictPass, fmt.Sprintf("last act matching %s has expected status %s", check.Path, check.Status)
	}
	for _, index := range indexes {
		if trial.Acts[index].Status == check.Status {
			return VerdictPass, fmt.Sprintf("an act matching %s has expected status %s", check.Path, check.Status)
		}
	}
	return VerdictFail, fmt.Sprintf("no act matching %s has status %s", check.Path, check.Status)
}

// gradeActOutputSelected implements selector addressing for
// act_output_matches with the same "some"/"last" semantics.
func gradeActOutputSelected(check Check, trial Trial) (string, string) {
	indexes := selectedActIndexes(trial, check.Path)
	if len(indexes) == 0 {
		return VerdictFail, fmt.Sprintf("no act matches selector path %s", check.Path)
	}
	if check.SelectorMode == "last" {
		indexes = indexes[len(indexes)-1:]
	}
	for _, index := range indexes {
		matched, err := regexp.MatchString(check.Pattern, trial.Acts[index].Output)
		if err != nil {
			return VerdictFail, fmt.Sprintf("invalid pattern %q: %v", check.Pattern, err)
		}
		if check.Negate {
			matched = !matched
		}
		if matched {
			return VerdictPass, fmt.Sprintf("an act matching %s satisfies pattern %q", check.Path, check.Pattern)
		}
	}
	return VerdictFail, fmt.Sprintf("no selected act matching %s satisfies pattern %q", check.Path, check.Pattern)
}

func gradeActStatusAtPath(check Check, path matchedPath, trial Trial) (string, string) {
	index, ok := resolveActIndex(check.Seq, path, len(trial.Acts))
	if !ok {
		return VerdictFail, "trial has no act at the required sequence"
	}
	act := trial.Acts[index]
	if act.Status != check.Status {
		return VerdictFail, fmt.Sprintf("act %d has status %s, expected %s", index, act.Status, check.Status)
	}
	return VerdictPass, fmt.Sprintf("act %d has expected status %s", index, check.Status)
}

func gradeActOutputAtPath(check Check, path matchedPath, trial Trial) (string, string) {
	index, ok := resolveActIndex(check.Seq, path, len(trial.Acts))
	if !ok {
		return VerdictFail, "trial has no act at the required sequence"
	}
	matched, err := regexp.MatchString(check.Pattern, trial.Acts[index].Output)
	if err != nil {
		return VerdictFail, fmt.Sprintf("invalid pattern %q: %v", check.Pattern, err)
	}
	if check.Negate {
		matched = !matched
	}
	if !matched {
		return VerdictFail, fmt.Sprintf("act %d output does not satisfy pattern %q", index, check.Pattern)
	}
	return VerdictPass, fmt.Sprintf("act %d output satisfies pattern %q", index, check.Pattern)
}

func resolveActIndex(sequence *int, path matchedPath, actCount int) (int, bool) {
	if sequence == nil || *sequence < 0 {
		return 0, false
	}
	index := *sequence
	if path.found {
		if index >= len(path.offsets) {
			return 0, false
		}
		index = path.offsets[index]
	}
	return index, index < actCount
}

func findFile(trial Trial, path string) (FileState, bool) {
	for _, file := range trial.EndState.Files {
		if file.Path == path {
			return file, true
		}
	}
	return FileState{}, false
}

func gradeFileMatches(check Check, trial Trial) (string, string) {
	file, ok := findFile(trial, check.Path)
	if !ok || !file.Present {
		return VerdictFail, fmt.Sprintf("no end_state evidence that %s is present", check.Path)
	}
	matched, err := regexp.MatchString(check.Pattern, file.Content)
	if err != nil {
		return VerdictFail, fmt.Sprintf("invalid pattern %q: %v", check.Pattern, err)
	}
	if !matched {
		return VerdictFail, fmt.Sprintf("content of %s does not match pattern %q", check.Path, check.Pattern)
	}
	return VerdictPass, fmt.Sprintf("content of %s matches pattern %q", check.Path, check.Pattern)
}

func gradeFileAbsent(check Check, trial Trial) (string, string) {
	file, ok := findFile(trial, check.Path)
	if !ok {
		return VerdictFail, fmt.Sprintf("no end_state evidence for %s", check.Path)
	}
	if !file.Present {
		return VerdictPass, fmt.Sprintf("%s is recorded absent in end_state", check.Path)
	}
	return VerdictFail, fmt.Sprintf("%s is recorded present in end_state", check.Path)
}

func findCommandActs(trial Trial, command []string) []Act {
	matches := []Act{}
	for _, act := range trial.Acts {
		if exactSequence(act.Command, command) {
			matches = append(matches, act)
		}
	}
	return matches
}

// exactSequence reports whether two command slices are identical and
// case-sensitive, unlike acceptable-path verb comparisons.
func exactSequence(observed, expected []string) bool {
	if len(observed) != len(expected) {
		return false
	}
	for index, token := range expected {
		if observed[index] != token {
			return false
		}
	}
	return true
}

func gradeCommandExit(check Check, trial Trial) (string, string) {
	acts := findCommandActs(trial, check.Command)
	if len(acts) == 0 {
		return VerdictFail, fmt.Sprintf("no recorded act ran command %v", check.Command)
	}
	if check.ExitCode == nil {
		return VerdictFail, "check has no exit_code to compare"
	}
	for _, act := range acts {
		if act.ExitCode != nil && *act.ExitCode == *check.ExitCode {
			return VerdictPass, fmt.Sprintf("command %v exited %d as expected", check.Command, *check.ExitCode)
		}
	}
	return VerdictFail, fmt.Sprintf("no recorded %v command exited %d", check.Command, *check.ExitCode)
}

func gradeCommandOutputMatches(check Check, trial Trial) (string, string) {
	acts := findCommandActs(trial, check.Command)
	if len(acts) == 0 {
		return VerdictFail, fmt.Sprintf("no recorded act ran command %v", check.Command)
	}
	for _, act := range acts {
		matched, err := regexp.MatchString(check.Pattern, act.Output)
		if err != nil {
			return VerdictFail, fmt.Sprintf("invalid pattern %q: %v", check.Pattern, err)
		}
		if check.Negate {
			matched = !matched
		}
		if matched {
			return VerdictPass, fmt.Sprintf("output of %v satisfies pattern %q", check.Command, check.Pattern)
		}
	}
	return VerdictFail, fmt.Sprintf("no output of %v satisfies pattern %q", check.Command, check.Pattern)
}

func gradeFinalMessageMatches(check Check, trial Trial) (string, string) {
	matched, err := regexp.MatchString(check.Pattern, trial.FinalMessage)
	if err != nil {
		return VerdictFail, fmt.Sprintf("invalid pattern %q: %v", check.Pattern, err)
	}
	if check.Negate {
		matched = !matched
	}
	if !matched {
		verb := "match"
		if check.Negate {
			verb = "not match"
		}
		return VerdictFail, fmt.Sprintf("final_message does not %s pattern %q", verb, check.Pattern)
	}
	return VerdictPass, fmt.Sprintf("final_message satisfies pattern %q", check.Pattern)
}

// GraderDigest pins the grader implementation for a label's grading_script
// field. It hashes the complete source, dependency, and schema set that can
// change a grade. Line endings are normalized so Windows and Linux checkouts
// produce the same identity.
func GraderDigest(root string) (string, error) {
	paths := []string{
		"experiments/frontier-v1/corpusctl/go.mod",
		"experiments/frontier-v1/corpusctl/go.sum",
	}
	for _, pattern := range []string{
		"experiments/frontier-v1/corpusctl/cmd/corpusctl/*.go",
		"experiments/frontier-v1/corpusctl/internal/corpus/*.go",
		"experiments/frontier-v1/schema/*.schema.json",
	} {
		matches, err := filepath.Glob(filepath.Join(root, filepath.FromSlash(pattern)))
		if err != nil {
			return "", err
		}
		for _, match := range matches {
			if strings.HasSuffix(match, "_test.go") {
				continue
			}
			relative, err := filepath.Rel(root, match)
			if err != nil {
				return "", err
			}
			paths = append(paths, filepath.ToSlash(relative))
		}
	}
	sort.Strings(paths)
	files := make([]any, 0, len(paths))
	for _, relative := range paths {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(relative)))
		if err != nil {
			return "", fmt.Errorf("read grader artifact %s: %w", relative, err)
		}
		data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
		sum := sha256.Sum256(data)
		files = append(files, map[string]any{"path": relative, "sha256": hex.EncodeToString(sum[:])})
	}
	value := map[string]any{
		"v":     json.Number("1"),
		"files": files,
	}
	return Digest(value)
}

// EncodeGradeResult validates and renders a deterministic grade result.
func EncodeGradeResult(root string, result GradeResult) ([]byte, error) {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, err
	}
	value, err := DecodeJSON(data)
	if err != nil {
		return nil, err
	}
	schemas, err := loadSchemas(root)
	if err != nil {
		return nil, err
	}
	if err := schemas.validateValue("grade-result", value); err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}
