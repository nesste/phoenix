package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestClosureVerdictParserAgreesWithTheCommittedReviewCorpus runs the
// post-closure verdict parser over every review record this project has
// committed. The parser exists because matching the bare token ACCEPT
// anywhere in a file accepts a REVISE record, and its first implementation
// then refused the house style and accepted evaluator prompts. The corpus is
// the arbiter for both directions: an evaluator prompt listing the available
// verdicts is never a verdict statement, and a record whose stated verdict is
// REVISE is never accepted.
func TestClosureVerdictParserAgreesWithTheCommittedReviewCorpus(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	records, err := filepath.Glob(filepath.Join(root, "docs", "reviews", "*.md"))
	if err != nil {
		t.Fatal(err)
	}
	if len(records) < 60 {
		t.Fatalf("review corpus = %d records, expected the committed corpus", len(records))
	}
	accepted, refused := 0, 0
	for _, path := range records {
		contents, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		name := filepath.Base(path)
		err = verifyClosureReviewVerdictLine(string(contents))
		if err == nil {
			accepted++
			if strings.Contains(name, "-prompt") {
				t.Errorf("%s is an evaluator prompt and must not read as an accepting record", name)
			}
			continue
		}
		refused++
		if strings.Contains(name, "-prompt") {
			continue
		}
		// A non-prompt record is refused only when it genuinely states a
		// non-ACCEPT verdict or states none at all.
		if verdict, found := closureReviewVerdict(string(contents)); found && verdict == "ACCEPT" {
			t.Errorf("%s states ACCEPT but was refused: %v", name, err)
		}
	}
	if accepted == 0 {
		t.Fatalf("no committed record was accepted; the parser refuses the whole corpus")
	}
	t.Logf("review corpus: %d accepted, %d refused", accepted, refused)
}

// TestClosureVerdictParserReadsEveryCommittedVerdictStyle pins the concrete
// forms this project has actually used, so a future tightening cannot quietly
// stop recognizing one of them.
func TestClosureVerdictParserReadsEveryCommittedVerdictStyle(t *testing.T) {
	tests := []struct {
		name     string
		contents string
		verdict  string
		found    bool
	}{
		{"list item with a code span", "- **Verdict:** `ACCEPT`\n", "ACCEPT", true},
		{"list item with a code span, revise", "- **Verdict:** `REVISE`\n", "REVISE", true},
		{"bold inline", "**Verdict: ACCEPT** - no findings\n", "ACCEPT", true},
		{"bold inline, revise", "**Verdict: REVISE** - two P1 findings\n", "REVISE", true},
		{"title case", "**Verdict: Accept**\n", "ACCEPT", true},
		{"quoted token", "Verdict: \"ACCEPT\"\n", "ACCEPT", true},
		{"numbered heading then token", "## 1. Verdict\n\n**ACCEPT.**\n", "ACCEPT", true},
		{"heading then token", "## Verdict\n\nACCEPT\n", "ACCEPT", true},
		{"final verdict label", "**Final verdict: ACCEPT**\n", "ACCEPT", true},
		{"first statement wins over a later citation",
			"**Verdict: ACCEPT**\n\nThe first round returned **Verdict: REVISE**; all findings are fixed.\n", "ACCEPT", true},
		{"evaluator prompt boilerplate", "1. Verdict: ACCEPT, REVISE, or REJECT.\n", "", false},
		{"prompt boilerplate without numbering", "Verdict: ACCEPT or REVISE.\n", "", false},
		{"mid-sentence use", "The record must state a verdict: ACCEPT is required.\n", "", false},
		{"bare token with no label", "I would ACCEPT these residuals.\n", "", false},
		{"no verdict at all", "a record with no verdict line\n", "", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			verdict, found := closureReviewVerdict(test.contents)
			if found != test.found || verdict != test.verdict {
				t.Fatalf("verdict = %q/%t, want %q/%t", verdict, found, test.verdict, test.found)
			}
			err := verifyClosureReviewVerdictLine(test.contents)
			if (err == nil) != (test.verdict == "ACCEPT") {
				t.Fatalf("gate verdict for %q = %v", test.name, err)
			}
		})
	}
}
