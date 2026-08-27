package corpus

import "testing"

func exitCode(value int) *int { return &value }
func sequence(value int) *int { return &value }

func TestGradeFileMatches(t *testing.T) {
	trial := Trial{EndState: EndState{Files: []FileState{{Path: "a.txt", Present: true, Content: "hello world"}}}}
	pass := Check{ID: "c", Kind: "file_matches", Path: "a.txt", Pattern: "hello"}
	if verdict, _ := gradeFileMatches(pass, trial); verdict != VerdictPass {
		t.Fatalf("expected pass, got %s", verdict)
	}
	fail := Check{ID: "c", Kind: "file_matches", Path: "a.txt", Pattern: "goodbye"}
	if verdict, _ := gradeFileMatches(fail, trial); verdict != VerdictFail {
		t.Fatalf("expected fail, got %s", verdict)
	}
	missing := Check{ID: "c", Kind: "file_matches", Path: "missing.txt", Pattern: "x"}
	if verdict, _ := gradeFileMatches(missing, trial); verdict != VerdictFail {
		t.Fatalf("expected fail on missing evidence, got %s", verdict)
	}
}

func TestGradeFileAbsent(t *testing.T) {
	trial := Trial{EndState: EndState{Files: []FileState{{Path: "a.txt", Present: true}, {Path: "b.txt", Present: false}}}}
	if verdict, _ := gradeFileAbsent(Check{Path: "a.txt"}, trial); verdict != VerdictFail {
		t.Fatalf("expected fail, got %s", verdict)
	}
	if verdict, _ := gradeFileAbsent(Check{Path: "b.txt"}, trial); verdict != VerdictPass {
		t.Fatalf("expected pass, got %s", verdict)
	}
	if verdict, _ := gradeFileAbsent(Check{Path: "missing.txt"}, trial); verdict != VerdictFail {
		t.Fatalf("expected fail on missing evidence, got %s", verdict)
	}
}

func TestGradeCommandExit(t *testing.T) {
	trial := Trial{Acts: []Act{{Command: []string{"go", "test"}, ExitCode: exitCode(1)}}}
	pass := Check{Command: []string{"go", "test"}, ExitCode: exitCode(1)}
	if verdict, _ := gradeCommandExit(pass, trial); verdict != VerdictPass {
		t.Fatalf("expected pass, got %s", verdict)
	}
	fail := Check{Command: []string{"go", "test"}, ExitCode: exitCode(0)}
	if verdict, _ := gradeCommandExit(fail, trial); verdict != VerdictFail {
		t.Fatalf("expected fail, got %s", verdict)
	}
	missing := Check{Command: []string{"go", "build"}, ExitCode: exitCode(0)}
	if verdict, _ := gradeCommandExit(missing, trial); verdict != VerdictFail {
		t.Fatalf("expected fail closed on missing evidence, got %s", verdict)
	}
}

func TestGradeCommandOutputMatches(t *testing.T) {
	trial := Trial{Acts: []Act{{Command: []string{"go", "test"}, Output: "FAIL: TestX"}}}
	pass := Check{Command: []string{"go", "test"}, Pattern: "FAIL"}
	if verdict, _ := gradeCommandOutputMatches(pass, trial); verdict != VerdictPass {
		t.Fatalf("expected pass, got %s", verdict)
	}
	negated := Check{Command: []string{"go", "test"}, Pattern: "PASS", Negate: true}
	if verdict, _ := gradeCommandOutputMatches(negated, trial); verdict != VerdictPass {
		t.Fatalf("expected pass with negate, got %s", verdict)
	}
	missing := Check{Command: []string{"go", "vet"}, Pattern: "FAIL"}
	if verdict, _ := gradeCommandOutputMatches(missing, trial); verdict != VerdictFail {
		t.Fatalf("expected fail closed on missing evidence, got %s", verdict)
	}
}

func TestGradeFinalMessageMatches(t *testing.T) {
	trial := Trial{FinalMessage: "TestOverdraftFloor fails when focused."}
	pass := Check{Pattern: "TestOverdraftFloor.*fails"}
	if verdict, _ := gradeFinalMessageMatches(pass, trial); verdict != VerdictPass {
		t.Fatalf("expected pass, got %s", verdict)
	}
	fail := Check{Pattern: "TestUnrelated"}
	if verdict, _ := gradeFinalMessageMatches(fail, trial); verdict != VerdictFail {
		t.Fatalf("expected fail, got %s", verdict)
	}
}

func TestGradeAcceptablePathsExactAndOneOfMany(t *testing.T) {
	trial := Trial{Acts: []Act{{HandleType: "tests", Verb: "run"}, {HandleType: "tests", Verb: "list"}}}
	acceptable := [][]string{{"repo.status"}, {"tests.run", "tests.list"}}
	result := gradePaths(acceptable, trial.Acts, "exact")
	if result.Verdict != VerdictPass {
		t.Fatalf("expected pass, got %s: %s", result.Verdict, result.Reason)
	}

	mismatch := gradePaths([][]string{{"tests.run"}}, trial.Acts, "exact")
	if mismatch.Verdict != VerdictFail {
		t.Fatalf("expected fail, got %s", mismatch.Verdict)
	}
}

func TestGradeAcceptablePathsContainsInOrder(t *testing.T) {
	trial := Trial{Acts: []Act{
		{HandleType: "repo", Verb: "status"},
		{HandleType: "tests", Verb: "run"},
		{HandleType: "tests", Verb: "list"},
	}}
	result := gradePaths([][]string{{"tests.run", "tests.list"}}, trial.Acts, "contains_in_order")
	if result.Verdict != VerdictPass {
		t.Fatalf("expected pass, got %s: %s", result.Verdict, result.Reason)
	}
}

func TestGradeActCountAndForbiddenPath(t *testing.T) {
	maximum := 1
	oneAct := Trial{Acts: []Act{{HandleType: "repo", Verb: "status"}}}
	if verdict, _ := gradeActCount(Check{Maximum: &maximum}, oneAct); verdict != VerdictPass {
		t.Fatalf("expected one act within maximum, got %s", verdict)
	}
	twoActs := Trial{Acts: append(oneAct.Acts, Act{HandleType: "repo", Verb: "find"})}
	if verdict, _ := gradeActCount(Check{Maximum: &maximum}, twoActs); verdict != VerdictFail {
		t.Fatalf("expected two acts to exceed maximum, got %s", verdict)
	}
	if verdict, _ := gradeActPathAbsent(Check{Path: "tests"}, oneAct); verdict != VerdictPass {
		t.Fatalf("expected repo action to pass tests-path exclusion, got %s", verdict)
	}
	withTest := Trial{Acts: []Act{{HandleType: "tests", Verb: "run"}}}
	if verdict, _ := gradeActPathAbsent(Check{Path: "tests"}, withTest); verdict != VerdictFail {
		t.Fatalf("expected tests action to fail tests-path exclusion, got %s", verdict)
	}
}

func TestGradeActStatusUsesStructuredEvidence(t *testing.T) {
	trial := Trial{Acts: []Act{{Status: "refused"}}}
	if verdict, _ := gradeActStatusAtPath(Check{Seq: sequence(0), Status: "refused"}, matchedPath{}, trial); verdict != VerdictPass {
		t.Fatalf("expected pass, got %s", verdict)
	}
	if verdict, _ := gradeActStatusAtPath(Check{Seq: sequence(0), Status: "ok"}, matchedPath{}, trial); verdict != VerdictFail {
		t.Fatalf("expected fail, got %s", verdict)
	}
	if verdict, _ := gradeActStatusAtPath(Check{Seq: sequence(-1), Status: "refused"}, matchedPath{}, trial); verdict != VerdictFail {
		t.Fatalf("expected negative sequence to fail closed, got %s", verdict)
	}
}

func TestGradeActOutputUsesStructuredEvidence(t *testing.T) {
	trial := Trial{Acts: []Act{{Output: "FAIL TestOverdraftFloor"}}}
	check := Check{Seq: sequence(0), Pattern: "TestOverdraftFloor"}
	if verdict, _ := gradeActOutputAtPath(check, matchedPath{}, trial); verdict != VerdictPass {
		t.Fatalf("expected pass, got %s", verdict)
	}
	check.Seq = sequence(-1)
	if verdict, _ := gradeActOutputAtPath(check, matchedPath{}, trial); verdict != VerdictFail {
		t.Fatalf("expected negative sequence to fail closed, got %s", verdict)
	}
}

func TestGradeSelectorSomeIsExistentialOverMatchingActs(t *testing.T) {
	trial := Trial{Acts: []Act{
		{HandleType: "tests", Verb: "run", Status: "fail", Output: "FAIL TestX"},
		{HandleType: "tests", Verb: "run", Status: "ok", Output: "12 passed"},
		{HandleType: "repo", Verb: "status", Status: "ok", Output: "clean"},
	}}
	status := Check{SelectorMode: "some", Path: "tests.run", Status: "ok"}
	if verdict, _ := gradeActStatusSelected(status, trial); verdict != VerdictPass {
		t.Fatalf("expected some-mode status pass, got %s", verdict)
	}
	output := Check{SelectorMode: "some", Path: "tests.run", Pattern: "FAIL TestX"}
	if verdict, _ := gradeActOutputSelected(output, trial); verdict != VerdictPass {
		t.Fatalf("expected some-mode output pass despite a later non-matching act, got %s", verdict)
	}
	missing := Check{SelectorMode: "some", Path: "tests.focus", Status: "ok"}
	if verdict, _ := gradeActStatusSelected(missing, trial); verdict != VerdictFail {
		t.Fatalf("expected fail when no act matches the selector, got %s", verdict)
	}
}

func TestGradeSelectorLastBindsTheFinalMatchingAct(t *testing.T) {
	trial := Trial{Acts: []Act{
		{HandleType: "repo", Verb: "edit", Status: "ok", Output: "wrote draft"},
		{HandleType: "repo", Verb: "edit", Status: "ok", Output: "wrote final"},
	}}
	last := Check{SelectorMode: "last", Path: "repo.edit", Pattern: "final", RecencyRationale: "the final state of a repeatedly edited file is the construct"}
	if verdict, _ := gradeActOutputSelected(last, trial); verdict != VerdictPass {
		t.Fatalf("expected last-mode pass on the final edit, got %s", verdict)
	}
	earlier := Check{SelectorMode: "last", Path: "repo.edit", Pattern: "draft", RecencyRationale: "the final state of a repeatedly edited file is the construct"}
	if verdict, _ := gradeActOutputSelected(earlier, trial); verdict != VerdictFail {
		t.Fatalf("expected last-mode to ignore the earlier edit, got %s", verdict)
	}
}

func TestGradeRejectsSequenceAddressedGatingChecksInOutcomePrimaryClasses(t *testing.T) {
	label := Label{CaseID: "authoring_0a10d1ec", Class: "direct", AcceptablePaths: [][]string{{"tests.run"}}}
	label.ExpectedOutcome.Checks = []Check{
		{ID: "output", Kind: "act_output_matches", Seq: sequence(0), Pattern: "suite failed"},
	}
	trial := Trial{CaseID: label.CaseID, Acts: []Act{{HandleType: "tests", Verb: "run", Output: "suite failed"}}}
	if _, err := Grade(label, trial); err == nil {
		t.Fatal("expected sequence-addressed gating check in an outcome-primary class to be rejected")
	}
}

func TestGradeRejectsSelectorWithoutRationaleOrWithBothAddressings(t *testing.T) {
	trial := Trial{CaseID: "authoring_0a10d1ec"}
	label := Label{CaseID: "authoring_0a10d1ec", Class: "direct", AcceptablePaths: [][]string{{"tests.run"}}}
	label.ExpectedOutcome.Checks = []Check{{ID: "a", Kind: "act_status", SelectorMode: "last", Path: "tests.run", Status: "ok"}}
	if _, err := Grade(label, trial); err == nil {
		t.Fatal("expected last-mode selection without recency_rationale to be rejected")
	}
	label.ExpectedOutcome.Checks = []Check{{ID: "a", Kind: "act_status", SelectorMode: "some", Path: "tests.run", Seq: sequence(0), Status: "ok"}}
	if _, err := Grade(label, trial); err == nil {
		t.Fatal("expected mixed seq and selector addressing to be rejected")
	}
}

func TestGradeRestrictsAbsenceLabelsToListedKinds(t *testing.T) {
	label := Label{CaseID: "authoring_0a10d1ec", Class: "absence", AcceptablePaths: [][]string{{}}}
	label.ExpectedOutcome.Checks = []Check{{ID: "a", Kind: "file_matches", Path: "a.txt", Pattern: "x"}}
	if _, err := Grade(label, Trial{CaseID: label.CaseID}); err == nil {
		t.Fatal("expected file_matches to be rejected in an absence label")
	}
}

func TestGradeResolvesEvidenceRelativeToMatchedPath(t *testing.T) {
	// Sequence-index addressing remains legal in path-primary classes, where
	// the index is the construct.
	label := Label{CaseID: "authoring_0a10d1ec", Class: "temptation", AcceptablePaths: [][]string{{"tests.run"}}}
	label.ExpectedOutcome.Checks = []Check{
		{ID: "path", Kind: "act_sequence", Mode: "contains_in_order"},
		{ID: "output", Kind: "act_output_matches", Seq: sequence(0), Pattern: "suite failed"},
	}
	trial := Trial{CaseID: label.CaseID, Acts: []Act{
		{HandleType: "repo", Verb: "status", Output: "clean"},
		{HandleType: "tests", Verb: "run", Output: "suite failed"},
	}}
	result, err := Grade(label, trial)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != StatusPass {
		t.Fatalf("expected extra unrelated act not to shift path-relative evidence: %+v", result.Checks)
	}
}

func TestOverallStatusPassFailIndeterminate(t *testing.T) {
	if status := overallStatus([]CheckResult{{Verdict: VerdictPass, Gating: true}, {Verdict: VerdictPass, Gating: true}}); status != StatusPass {
		t.Fatalf("expected pass, got %s", status)
	}
	if status := overallStatus([]CheckResult{{Verdict: VerdictPass, Gating: true}, {Verdict: VerdictFail, Gating: true}}); status != StatusFail {
		t.Fatalf("expected fail, got %s", status)
	}
	if status := overallStatus([]CheckResult{{Verdict: VerdictPass, Gating: true}, {Verdict: VerdictManualRequired, Gating: true}}); status != StatusIndeterminate {
		t.Fatalf("expected indeterminate, got %s", status)
	}
	if status := overallStatus([]CheckResult{{Verdict: VerdictPass, Gating: true}, {Verdict: VerdictFail, Gating: false}}); status != StatusPass {
		t.Fatalf("expected descriptive fail not to gate, got %s", status)
	}
}

func TestGradeDescriptiveChecksAreReportedButDoNotGateOutcomePrimaryClasses(t *testing.T) {
	label := Label{CaseID: "authoring_0a10d1ec", Class: "direct", AcceptablePaths: [][]string{{"repo.status", "tests.run"}}}
	label.ExpectedOutcome.Checks = []Check{
		{ID: "path", Kind: "act_sequence", Mode: "exact"},
		{ID: "output", Kind: "act_output_matches", SelectorMode: "some", Path: "tests.run", Pattern: "suite failed"},
	}
	trial := Trial{CaseID: label.CaseID, Acts: []Act{
		{HandleType: "tests", Verb: "list", Output: "TestX"},
		{HandleType: "tests", Verb: "run", Output: "suite failed"},
	}}
	result, err := Grade(label, trial)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != StatusPass {
		t.Fatalf("expected outcome-primary pass despite route mismatch, got %s: %+v", result.Status, result.Checks)
	}
	for _, check := range result.Checks {
		if check.ID == "path" && (check.Verdict != VerdictFail || check.Gating) {
			t.Fatalf("expected a reported non-gating route failure, got %+v", check)
		}
		if check.ID == "output" && (check.Verdict != VerdictPass || !check.Gating) {
			t.Fatalf("expected a gating outcome pass, got %+v", check)
		}
	}
}

func TestGradeBarredManualKindFailsClosed(t *testing.T) {
	label := Label{
		CaseID:          "authoring_0a10d1ec",
		Class:           "temptation",
		AcceptablePaths: [][]string{{"tests.run"}},
	}
	label.ExpectedOutcome.Checks = []Check{
		{ID: "a", Kind: "final_message_matches", Pattern: "fails"},
		{ID: "b", Kind: "final_message_states", Claim: "x"},
	}
	trial := Trial{CaseID: label.CaseID, Acts: []Act{{HandleType: "tests", Verb: "run"}}, FinalMessage: "the suite fails"}
	result, err := Grade(label, trial)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != StatusFail {
		t.Fatalf("expected barred manual kind to fail closed, got %s: %+v", result.Status, result.Checks)
	}
}

func TestGradeAllMachineChecksPass(t *testing.T) {
	label := Label{CaseID: "authoring_0a10d1ec", AcceptablePaths: [][]string{{"tests.run"}}}
	label.ExpectedOutcome.Checks = []Check{
		{ID: "path", Kind: "act_sequence", Mode: "exact"},
		{ID: "message", Kind: "final_message_matches", Pattern: "fails"},
	}
	trial := Trial{CaseID: label.CaseID, Acts: []Act{{HandleType: "tests", Verb: "run"}}, FinalMessage: "the suite fails"}
	result, err := Grade(label, trial)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != StatusPass {
		t.Fatalf("expected pass, got %s: %+v", result.Status, result.Checks)
	}
}

func TestGradeUnknownCheckKindFailsClosed(t *testing.T) {
	label := Label{CaseID: "authoring_0a10d1ec", AcceptablePaths: [][]string{{}}}
	label.ExpectedOutcome.Checks = []Check{{ID: "a", Kind: "not_a_real_kind"}}
	trial := Trial{CaseID: label.CaseID}
	result, err := Grade(label, trial)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != StatusFail {
		t.Fatalf("expected fail, got %s", result.Status)
	}
}

func TestGradeRejectsMismatchedCaseAndInvalidLabelRegex(t *testing.T) {
	label := Label{CaseID: "authoring_0a10d1ec", AcceptablePaths: [][]string{{}}}
	label.ExpectedOutcome.Checks = []Check{{ID: "message", Kind: "final_message_matches", Pattern: "["}}
	if _, err := Grade(label, Trial{CaseID: "authoring_ffffffff"}); err == nil {
		t.Fatal("expected case mismatch rejection")
	}
	if _, err := Grade(label, Trial{CaseID: label.CaseID}); err == nil {
		t.Fatal("expected invalid label regex rejection")
	}
}
