package main

import (
	"reflect"
	"testing"
)

func TestWilliamsSequencesBalanceOddArmPredecessors(t *testing.T) {
	rows, err := williamsSequences(phase1Arms)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 10 {
		t.Fatalf("Williams rows = %d, want 10", len(rows))
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
			if left != right && predecessors[left+"->"+right] != 2 {
				t.Fatalf("predecessor %s->%s count = %d, want 2", left, right, predecessors[left+"->"+right])
			}
		}
	}
}

func TestPhase1ScheduleIsDeterministicFamilyBlockedAndPairingComplete(t *testing.T) {
	cases := scheduleTestCases()
	first, err := generatePhase1Schedule(cases)
	if err != nil {
		t.Fatal(err)
	}
	second, err := generatePhase1Schedule(cases)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, second) {
		t.Fatal("same schedule seed produced different schedules")
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

func TestScheduleValidationRejectsOrderMutation(t *testing.T) {
	cases := scheduleTestCases()
	schedule, err := generatePhase1Schedule(cases)
	if err != nil {
		t.Fatal(err)
	}
	schedule.Entries[0].Arm, schedule.Entries[1].Arm = schedule.Entries[1].Arm, schedule.Entries[0].Arm
	if err := validateSchedule(schedule, cases); err == nil {
		t.Fatal("schedule validation accepted a mutated arm order")
	}
}

func TestScheduleValidationRejectsNonProtocolDesign(t *testing.T) {
	cases := scheduleTestCases()
	schedule, err := generateSchedule(cases, phase1ScheduleSeed+1, phase1Repetitions, phase1Arms)
	if err != nil {
		t.Fatal(err)
	}
	if err := validateSchedule(schedule, cases); err == nil {
		t.Fatal("schedule validation accepted a non-protocol seed")
	}
}

func scheduleTestCases() []runnableCase {
	return []runnableCase{
		{CaseID: "authoring_a1", FamilyID: "family_a"},
		{CaseID: "authoring_a2", FamilyID: "family_a"},
		{CaseID: "authoring_b1", FamilyID: "family_b"},
		{CaseID: "authoring_b2", FamilyID: "family_b"},
	}
}
