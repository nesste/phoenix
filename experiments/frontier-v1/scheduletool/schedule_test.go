package main

import (
	"reflect"
	"testing"
)

func TestValidationScheduleIsDeterministicFamilyBlockedAndComplete(t *testing.T) {
	cases := scheduleTestCases()
	first, err := generateSchedule("validation", cases, phase1ScheduleSeed, phase1Repetitions, phase1Arms)
	if err != nil {
		t.Fatal(err)
	}
	second, err := generateSchedule("validation", cases, phase1ScheduleSeed, phase1Repetitions, phase1Arms)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, second) {
		t.Fatal("same schedule inputs produced different schedules")
	}
	if first.Tranche != "validation" {
		t.Fatalf("tranche = %q, want validation", first.Tranche)
	}
	wantEntries := len(cases) * phase1Repetitions * len(phase1Arms)
	if len(first.Entries) != wantEntries {
		t.Fatalf("schedule entries = %d, want %d", len(first.Entries), wantEntries)
	}
	closedFamilies := map[string]bool{}
	currentFamily := ""
	for start := 0; start < len(first.Entries); start += len(phase1Arms) {
		group := first.Entries[start : start+len(phase1Arms)]
		if group[0].FamilyID != currentFamily {
			if closedFamilies[group[0].FamilyID] {
				t.Fatalf("family %s appears in multiple blocks", group[0].FamilyID)
			}
			if currentFamily != "" {
				closedFamilies[currentFamily] = true
			}
			currentFamily = group[0].FamilyID
		}
		arms := map[string]bool{}
		for offset, entry := range group {
			if entry.LaunchIndex != start+offset || entry.PairingIndex != group[0].PairingIndex ||
				entry.CaseID != group[0].CaseID || entry.Repetition != group[0].Repetition {
				t.Fatalf("pairing group is not contiguous and complete: %#v", group)
			}
			arms[entry.Arm] = true
		}
		if len(arms) != len(phase1Arms) {
			t.Fatalf("pairing group arms = %v", arms)
		}
	}
}

func TestWilliamsSequencesBalanceEvenArmPredecessors(t *testing.T) {
	rows, err := williamsSequences(phase1Arms)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 4 {
		t.Fatalf("Williams rows = %d, want 4 for the four Phase 1 arms", len(rows))
	}
	predecessors := map[string]int{}
	for _, row := range rows {
		seen := map[string]bool{}
		for index, arm := range row {
			seen[arm] = true
			if index > 0 {
				predecessors[row[index-1]+"->"+arm]++
			}
		}
		if len(seen) != len(phase1Arms) {
			t.Fatalf("Williams row is not a permutation: %v", row)
		}
	}
	for _, left := range phase1Arms {
		for _, right := range phase1Arms {
			if left != right && predecessors[left+"->"+right] != 1 {
				t.Fatalf("predecessor %s->%s count = %d, want 1", left, right, predecessors[left+"->"+right])
			}
		}
	}
}

func scheduleTestCases() []scheduleCase {
	return []scheduleCase{
		{CaseID: "validation_a1", FamilyID: "family_a"},
		{CaseID: "validation_a2", FamilyID: "family_a"},
		{CaseID: "validation_b1", FamilyID: "family_b"},
		{CaseID: "validation_b2", FamilyID: "family_b"},
	}
}
