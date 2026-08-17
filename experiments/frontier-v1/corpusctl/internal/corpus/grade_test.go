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

func TestGradeActStatusUsesStructuredEvidence(t *testing.T) {
	trial := Trial{Acts: []Act{{Status: "refused"}}}
	if verdict, _ := gradeActStatus(Check{Seq: sequence(0), Status: "refused"}, trial); verdict != VerdictPass {
		t.Fatalf("expected pass, got %s", verdict)
	}
	if verdict, _ := gradeActStatus(Check{Seq: sequence(0), Status: "ok"}, trial); verdict != VerdictFail {
		t.Fatalf("expected fail, got %s", verdict)
	}
}

func TestGradeActOutputUsesStructuredEvidence(t *testing.T) {
	trial := Trial{Acts: []Act{{Output: "FAIL TestOverdraftFloor"}}}
	check := Check{Seq: sequence(0), Pattern: "TestOverdraftFloor"}
	if verdict, _ := gradeActOutputMatches(check, trial); verdict != VerdictPass {
		t.Fatalf("expected pass, got %s", verdict)
	}
}

func TestOverallStatusPassFailIndeterminate(t *testing.T) {
	if status := overallStatus([]CheckResult{{Verdict: VerdictPass}, {Verdict: VerdictPass}}); status != StatusPass {
		t.Fatalf("expected pass, got %s", status)
	}
	if status := overallStatus([]CheckResult{{Verdict: VerdictPass}, {Verdict: VerdictFail}}); status != StatusFail {
		t.Fatalf("expected fail, got %s", status)
	}
	if status := overallStatus([]CheckResult{{Verdict: VerdictPass}, {Verdict: VerdictManualRequired}}); status != StatusIndeterminate {
		t.Fatalf("expected indeterminate, got %s", status)
	}
}

func TestGradeMixedMachineAndManualIsIndeterminate(t *testing.T) {
	label := Label{
		CaseID:          "authoring_direct_001",
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
	if result.Status != StatusIndeterminate {
		t.Fatalf("expected indeterminate, got %s: %+v", result.Status, result.Checks)
	}
	var sawManual bool
	for _, check := range result.Checks {
		if check.ID == "b" {
			if check.Verdict != VerdictManualRequired {
				t.Fatalf("expected manual_required, got %s", check.Verdict)
			}
			sawManual = true
		}
	}
	if !sawManual {
		t.Fatal("expected a manual_required check in results")
	}
}

func TestGradeAllMachineChecksPass(t *testing.T) {
	label := Label{CaseID: "authoring_direct_001", AcceptablePaths: [][]string{{"tests.run"}}}
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
	label := Label{CaseID: "authoring_direct_001", AcceptablePaths: [][]string{{}}}
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
	label := Label{CaseID: "authoring_direct_001", AcceptablePaths: [][]string{{}}}
	label.ExpectedOutcome.Checks = []Check{{ID: "message", Kind: "final_message_matches", Pattern: "["}}
	if _, err := Grade(label, Trial{CaseID: "authoring_direct_002"}); err == nil {
		t.Fatal("expected case mismatch rejection")
	}
	if _, err := Grade(label, Trial{CaseID: label.CaseID}); err == nil {
		t.Fatal("expected invalid label regex rejection")
	}
}
