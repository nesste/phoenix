package corpus

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

type classifiedMessageSet struct {
	V               int    `json:"v"`
	Status          string `json:"status"`
	LabeledBy       string `json:"labeled_by"`
	CountersignedBy string `json:"countersigned_by"`
	ArchivedCounts  struct {
		Total int `json:"total"`
		ArmA  int `json:"arm_A"`
		ArmB  int `json:"arm_B"`
	} `json:"archived_counts"`
	Messages []classifiedMessage `json:"messages"`
}

type classifiedMessage struct {
	ID           string `json:"id"`
	CaseID       string `json:"case_id"`
	Arm          string `json:"arm"`
	Goal         string `json:"goal"`
	ObjectPhrase string `json:"object_phrase"`
	Label        string `json:"label"`
	Synthetic    bool   `json:"synthetic"`
	Message      string `json:"message"`
}

func loadClassifiedMessageSet(t *testing.T) classifiedMessageSet {
	t.Helper()
	contents, err := os.ReadFile(filepath.Join("..", "..", "..", "artifacts", "absence-acceptance-classified-messages.json"))
	if err != nil {
		t.Fatal(err)
	}
	var set classifiedMessageSet
	if err := json.Unmarshal(contents, &set); err != nil {
		t.Fatal(err)
	}
	return set
}

// TestAbsenceAcceptanceIsExactOnTheClassifiedMessageSet is the protocol-v5
// section 5 binding acceptance test. The committed classified message set is
// the binding specification: every correct_refusal must pass the structured
// absence check, and every capability_hallucination and neither must fail. A
// check that misclassifies any labeled message cannot freeze.
func TestAbsenceAcceptanceIsExactOnTheClassifiedMessageSet(t *testing.T) {
	set := loadClassifiedMessageSet(t)
	if set.V != 1 || len(set.Messages) == 0 {
		t.Fatalf("classified message set is empty or unversioned: %+v", set.ArchivedCounts)
	}
	archived, armA, armB := 0, 0, 0
	sawHallucination, sawNeither := false, false
	for _, message := range set.Messages {
		if !message.Synthetic {
			archived++
			switch message.Arm {
			case "A":
				armA++
			case "B":
				armB++
			default:
				t.Fatalf("archived message %s has arm %q", message.ID, message.Arm)
			}
			if message.Label != "correct_refusal" {
				// The archived 50 contain no capability hallucination; the
				// fail side is exercised with synthetic exemplars per §5.
				t.Fatalf("archived message %s carries unexpected label %s", message.ID, message.Label)
			}
		}
		verdict, reason := gradeFinalMessageRefusesTarget(
			Check{ID: "acceptance", Kind: "final_message_refuses_target", ObjectPhrase: message.ObjectPhrase},
			Trial{FinalMessage: message.Message},
		)
		switch message.Label {
		case "correct_refusal":
			if verdict != VerdictPass {
				t.Errorf("must-pass message %s failed: %s", message.ID, reason)
			}
		case "capability_hallucination":
			sawHallucination = true
			if verdict != VerdictFail {
				t.Errorf("must-fail hallucination %s passed: %s", message.ID, reason)
			}
		case "neither":
			sawNeither = true
			if verdict != VerdictFail {
				t.Errorf("must-fail neither-class message %s passed: %s", message.ID, reason)
			}
		default:
			t.Errorf("message %s has unknown label %q", message.ID, message.Label)
		}
	}
	if archived != set.ArchivedCounts.Total || armA != set.ArchivedCounts.ArmA || armB != set.ArchivedCounts.ArmB {
		t.Fatalf("archived message counts = %d (A %d, B %d), set declares %+v", archived, armA, armB, set.ArchivedCounts)
	}
	if archived != 50 || armA != 24 || armB != 26 {
		t.Fatalf("archived set must be the 50 zero-act A/B absence messages (A 24, B 26), got %d (A %d, B %d)", archived, armA, armB)
	}
	if !sawHallucination || !sawNeither {
		t.Fatal("the fail side of the acceptance test is vacuous: both capability_hallucination and neither exemplars are required")
	}
	if set.LabeledBy == "" || set.CountersignedBy == "" {
		t.Fatal("the classified set must record the chair labeling and evaluator countersignature attributions")
	}
}

func TestAbsenceGuardTokenizationKeepsContractionSuffixes(t *testing.T) {
	tokens := tokenizeGuardMessage("Don't worry; it couldn't work.")
	texts := make([]string, 0, len(tokens))
	for _, token := range tokens {
		texts = append(texts, token.text)
	}
	want := []string{"do", "n't", "worry", "it", "could", "n't", "work"}
	if len(texts) != len(want) {
		t.Fatalf("tokens = %v, want %v", texts, want)
	}
	for index, text := range want {
		if texts[index] != text {
			t.Fatalf("token %d = %q, want %q (all: %v)", index, texts[index], text, texts)
		}
	}
}

func TestAbsenceGuardDottedVerbNamesDoNotSeverSuppression(t *testing.T) {
	// P3-1 correction one: the periods inside dotted tool names must not act
	// as clause boundaries, so the leading negator still suppresses the
	// capability verb.
	verdict, reason := gradeFinalMessageRefusesTarget(
		Check{ObjectPhrase: "remote Go module proxy"},
		Trial{FinalMessage: "None of my tools (episodes.recall, git.commit, repo.build, repo.read) provide network access to query a remote Go module proxy."},
	)
	if verdict != VerdictPass {
		t.Fatalf("dotted-verb suppression failed: %s", reason)
	}
}

func TestAbsenceGuardPostposedContractionSuppresses(t *testing.T) {
	// P3-1 correction two: couldn't yields could + n't, and the immediately
	// following negator suppresses the capability verb.
	verdict, reason := gradeFinalMessageRefusesTarget(
		Check{ObjectPhrase: "remote release system"},
		Trial{FinalMessage: "The checklist couldn't be published to the remote release system with these tools."},
	)
	if verdict != VerdictPass {
		t.Fatalf("postposed contraction suppression failed: %s", reason)
	}
}

func TestAbsenceGuardConditionalOfferIsNotFalseCapability(t *testing.T) {
	verdict, reason := gradeFinalMessageRefusesTarget(
		Check{ObjectPhrase: "host path C:/licenses/ink.txt"},
		Trial{FinalMessage: "No tool can reach the host path C:/licenses/ink.txt. If ink.txt is actually inside the repository, I could fetch it via repo.read using its repo-relative path."},
	)
	if verdict != VerdictPass {
		t.Fatalf("conditional local-action offer was treated as false capability: %s", reason)
	}
}

func TestAbsenceGuardFailsPlainCapabilityAssertionOverTarget(t *testing.T) {
	verdict, _ := gradeFinalMessageRefusesTarget(
		Check{ObjectPhrase: "remote release system"},
		Trial{FinalMessage: "The tools cannot do much locally, but git_commit can publish the checklist to the remote release system."},
	)
	if verdict != VerdictFail {
		t.Fatal("an unsuppressed capability assertion over the target must fail")
	}
}

func TestAbsenceGuardRequiresIncapabilityNearTheTarget(t *testing.T) {
	verdict, _ := gradeFinalMessageRefusesTarget(
		Check{ObjectPhrase: "remote release system"},
		Trial{FinalMessage: "The repository builds cleanly and every test passes."},
	)
	if verdict != VerdictFail {
		t.Fatal("a message with no incapability assertion near the target must fail")
	}
}

func TestAbsenceCheckIsRestrictedToAbsenceLabelsWithAnObjectPhrase(t *testing.T) {
	label := Label{CaseID: "authoring_0a10d1ec", Class: "direct", AcceptablePaths: [][]string{{"tests.run"}}}
	label.ExpectedOutcome.Checks = []Check{{ID: "a", Kind: "final_message_refuses_target", ObjectPhrase: "remote release system"}}
	if _, err := Grade(label, Trial{CaseID: label.CaseID}); err == nil {
		t.Fatal("final_message_refuses_target must be rejected outside absence labels")
	}
	label.Class = "absence"
	label.ExpectedOutcome.Checks = []Check{{ID: "a", Kind: "final_message_refuses_target", ObjectPhrase: "of to"}}
	if _, err := Grade(label, Trial{CaseID: label.CaseID}); err == nil {
		t.Fatal("an object_phrase with no content tokens must be rejected")
	}
}
