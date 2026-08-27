package main

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestCommittedValidationScheduleIsRetiredUnderV5Arms(t *testing.T) {
	// Protocol v5 removes arm E, so the committed v4 five-arm schedule can no
	// longer be reproduced by the deterministic A-D construction. The retired
	// bytes stay committed as custody evidence; a new disjoint sealed tranche
	// must produce the next schedule.
	root, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	cases, manifestDigest, err := loadValidationCases(root)
	if err != nil {
		t.Fatal(err)
	}
	if manifestDigest != "sha256:39acbad5e45ad65302659cd0875bdfe589165ede9b60b6448ac4b09ccfb1e0c6" {
		t.Fatalf("manifest digest = %s", manifestDigest)
	}
	regenerated, err := generateSchedule("validation", cases, phase1ScheduleSeed, phase1Repetitions, phase1Arms)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "experiments", "frontier-v1", "schedules", "validation.json")
	var committed launchSchedule
	if err := decodeStrict(path, &committed); err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(committed, regenerated) {
		t.Fatal("retired five-arm schedule unexpectedly matches the A-D construction")
	}
	if len(regenerated.Arms) != 4 || len(regenerated.Entries) != 120*phase1Repetitions*4 {
		t.Fatalf("A-D construction = %d arms, %d entries", len(regenerated.Arms), len(regenerated.Entries))
	}
	digest, err := digestJSONFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if digest != "sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813" {
		t.Fatalf("schedule digest = %s", digest)
	}
}
