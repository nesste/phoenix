package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const closureVerdictExpectationsPath = "closure-verdict-corpus-expectations.json"

type closureVerdictExpectations struct {
	V          int               `json:"v"`
	Purpose    string            `json:"purpose"`
	Derivation string            `json:"derivation"`
	Records    map[string]string `json:"records"`
	Counts     map[string]int    `json:"counts"`
}

// TestClosureVerdictParserMatchesTheCommittedCorpusExpectations checks the
// post-closure verdict parser against committed ground truth, in both
// directions. The labels were derived independently of the implementation, so
// a record labelled ACCEPT being refused fails here just as loudly as a
// prompt being accepted - the first payload's parser passed a test that could
// only ever catch the second.
//
// The corpus is pinned rather than swept from the live docs/reviews directory:
// a frozen test whose outcome depends on files added by unrelated later work
// can be turned red by that work, and going green again would cost a full
// freeze cycle.
func TestClosureVerdictParserMatchesTheCommittedCorpusExpectations(t *testing.T) {
	expectations := loadClosureVerdictExpectations(t)
	root := filepath.Join("..", "..", "..")
	seen := map[string]int{}
	for name, expected := range expectations.Records {
		contents, err := os.ReadFile(filepath.Join(root, "docs", "reviews", name))
		if err != nil {
			t.Fatalf("pinned corpus record %s: %v", name, err)
		}
		seen[expected]++
		verdict, found := closureReviewVerdict(string(contents))
		assertClosureVerdict(t, name, expected, verdict, found)

		gateErr := verifyClosureReviewVerdictLine(string(contents))
		if expected == "ACCEPT" && gateErr != nil {
			t.Errorf("%s states ACCEPT but the closure gate refused it: %v", name, gateErr)
		}
		if expected != "ACCEPT" && gateErr == nil {
			t.Errorf("%s states %q but the closure gate accepted it", name, expected)
		}
	}
	for verdict, want := range expectations.Counts {
		if seen[verdict] != want {
			t.Errorf("corpus holds %d records labelled %s, expectations record %d", seen[verdict], verdict, want)
		}
	}
}

func assertClosureVerdict(t *testing.T, name, expected, verdict string, found bool) {
	t.Helper()
	if expected == "none" {
		if found {
			t.Errorf("%s states no verdict of its own, parser read %q", name, verdict)
		}
		return
	}
	if !found {
		t.Errorf("%s states %s, parser found no verdict", name, expected)
		return
	}
	if verdict != expected {
		t.Errorf("%s states %s, parser read %q", name, expected, verdict)
	}
}

func loadClosureVerdictExpectations(t *testing.T) closureVerdictExpectations {
	t.Helper()
	var expectations closureVerdictExpectations
	readJSONForTest(t, filepath.Join("..", "artifacts", closureVerdictExpectationsPath), &expectations)
	if expectations.V != 1 || len(expectations.Records) < 60 {
		t.Fatalf("closure verdict expectations = v%d, %d records", expectations.V, len(expectations.Records))
	}
	for _, required := range []string{"BOTH directions", "pinned rather than swept"} {
		if !strings.Contains(expectations.Purpose, required) {
			t.Fatalf("expectations purpose is missing %q", required)
		}
	}
	return expectations
}

// TestClosureVerdictParserReadsEveryCommittedVerdictStyle pins the concrete
// forms this project has used, and the forms that must never read as a
// verdict statement, so a later tightening cannot quietly stop recognizing one
// or start accepting the other.
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
		{"title case", "**Verdict: Accept**\n", "ACCEPT", true},
		{"quoted token", "Verdict: \"ACCEPT\"\n", "ACCEPT", true},
		{"em dash separator", "**Verdict — ACCEPT**\n", "ACCEPT", true},
		{"en dash separator", "**Verdict – ACCEPT**\n", "ACCEPT", true},
		{"footnote marker", "- [^1] **Verdict:** `ACCEPT`\n", "ACCEPT", true},
		{"numbered heading then token", "## 1. Verdict\n\n**ACCEPT.**\n", "ACCEPT", true},
		{"heading then token", "## Verdict\n\nACCEPT\n", "ACCEPT", true},
		{"final verdict label", "**Final verdict: ACCEPT**\n", "ACCEPT", true},
		{"reject is a verdict, not silence", "**Verdict: REJECT**\n", "REJECT", true},
		{"first statement wins over a later citation",
			"**Verdict: ACCEPT**\n\nThe first round returned **Verdict: REVISE**; all findings are fixed.\n", "ACCEPT", true},
		{"prose citing another verdict on the same line",
			"**Verdict: REVISE** - it refuses 17 committed ACCEPT records\n", "REVISE", true},
		// Enumerations of the available verdicts, in either word order.
		{"prompt boilerplate", "1. Verdict: ACCEPT, REVISE, or REJECT.\n", "", false},
		{"prompt boilerplate reordered", "1. Verdict: ACCEPT, REJECT, or REVISE.\n", "", false},
		{"prompt boilerplate with an aside", "Verdict: ACCEPT, with findings, or REVISE.\n", "", false},
		{"prompt boilerplate under a heading", "## Verdict\n\nACCEPT, REJECT, or REVISE.\n", "", false},
		{"slash separated", "Verdict: ACCEPT/REVISE\n", "", false},
		// Quoted, fenced or indented verdicts belong to another record.
		{"fenced block", "```\nVerdict: ACCEPT\n```\n", "", false},
		{"fenced markdown block", "```markdown\n- **Verdict:** `ACCEPT`\n```\n", "", false},
		{"block quote", "> **Verdict:** `ACCEPT`\n", "", false},
		{"indented block quote", "  > Verdict: ACCEPT\n", "", false},
		{"indented code block", "    Verdict: ACCEPT\n", "", false},
		{"heading whose token line is quoted", "## Verdict\n\n> ACCEPT\n", "", false},
		// Prose that merely contains the word.
		{"no separator after the label", "Verdict ACCEPT was withheld for the reasons below.\n", "", false},
		{"mid-sentence use", "The record must state a verdict: ACCEPT is required.\n", "", false},
		{"bare token with no label", "I would ACCEPT these residuals.\n", "", false},
		{"accepted is not accept", "**Verdict: ACCEPTED**\n", "", false},
		{"no verdict at all", "a record with no verdict line\n", "", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			verdict, found := closureReviewVerdict(test.contents)
			if found != test.found || verdict != test.verdict {
				t.Fatalf("verdict = %q/%t, want %q/%t", verdict, found, test.verdict, test.found)
			}
			if err := verifyClosureReviewVerdictLine(test.contents); (err == nil) != (test.verdict == "ACCEPT") {
				t.Fatalf("gate verdict = %v", err)
			}
		})
	}
}

// TestClosureVerdictParserHandlesWindowsLineEndings guards the CRLF path,
// since review records are edited on both hosts.
func TestClosureVerdictParserHandlesWindowsLineEndings(t *testing.T) {
	for _, contents := range []string{
		"- **Verdict:** `ACCEPT`\r\n",
		"## 1. Verdict\r\n\r\n**ACCEPT.**\r\n",
	} {
		if verdict, found := closureReviewVerdict(contents); !found || verdict != "ACCEPT" {
			t.Fatalf("CRLF verdict = %q/%t", verdict, found)
		}
	}
}
