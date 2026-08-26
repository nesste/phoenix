package main

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestCommittedValidationScheduleMatchesPublicSealedManifest(t *testing.T) {
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
	expected, err := generateSchedule("validation", cases, phase1ScheduleSeed, phase1Repetitions, phase1Arms)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "experiments", "frontier-v1", "schedules", "validation.json")
	var actual launchSchedule
	if err := decodeStrict(path, &actual); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatal("committed validation schedule differs from deterministic construction")
	}
	digest, err := digestJSONFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if digest != "sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813" {
		t.Fatalf("schedule digest = %s", digest)
	}
}
