# Protocol v5 §5 absence-acceptance payload independent review prompt

Copy the text below into a fresh evaluator session that did not author this payload, the classified-set draft, any revision of the v5 proposal, or any prior review. This reviewer is also the countersigning evaluator required by proposal §5 and the fourth review's freeze obligation 1.

```text
Act as the independent evaluation reviewer and countersigning evaluator for the Phoenix protocol-v5 section 5 absence-grading acceptance payload. The accepted v5 proposal requires the frozen absence checks to be exact on a committed classified message set labeled by the project chair and countersigned by you; your review gates whether the focused refreeze may proceed AND supplies the countersignature. It is not an outcome run, not Gate 1A, not a refreeze, and not permission to open any sealed tranche; both outcome gates are closed and must stay closed.

Repository working tree: D:\Work\personal\phoenix
Base commit (parent of the payload): 6e0f35afe6128a971369a710256691f6cc8de311
Payload commit (HEAD, the candidate you review): 54e256e9f4347d34844a3ef9b2156600360dcd13

Identify the review inputs before reviewing. Compute SHA-256 over LF-normalized bytes (CRLF->LF then CR->LF) and require these exact values; return REVISE without reviewing if any differs:

- docs/plans/2026-08-27-protocol-v5-proposal.md (the accepted proposal; section 5 is the specification): sha256:1f977edfec3cce7ed1c53571cb73712edcaba53eab47486a4122acd6ee5518e5
- docs/reviews/2026-08-27-protocol-v5-proposal-review-4.md (freeze obligations 1 and 2, P3-1 corrections): sha256:0b24b820d5a4a77b982fae02c948b57b978c8f3a749f05b384ee572135efc59c
- experiments/frontier-v1/artifacts/gate-1a-absence-acceptance-candidate.json (payload inventory, 18 files): sha256:3dd434bf29dd16bf206fb7d38ebf599ceacba6f74feb7a3a6a1dbeb923548d22
- experiments/frontier-v1/artifacts/absence-acceptance-candidate.md (candidate note): sha256:88fb69ee577a00a5efb73f5732e3efc78a32da5458a16eaf882bed67f8c6d3e4
- experiments/frontier-v1/artifacts/absence-acceptance-classified-messages.json (the chair-labeled classified set): sha256:be03b2847fac248b9523cd7bcc617791d1c359a69e444e22f629ba33ee6b54c0

Custody and independence rules

- Do not modify any repository file except writing your review to docs/reviews/2026-08-27-absence-acceptance-payload-review.md.
- You may read any public file, run git read-only commands, and run go build / go vet / go test in the three modules. Running the corpusctl suite (including the binding acceptance test) and `go run ./cmd/corpusctl grader-digest --repo-root ../../..` and `corpusctl validate --tranche authoring` is expected.
- Never run a model, trial, or grade against any tranche; never create or modify anything under results/; never open anything under any held_out path, private label, or registry contents. The scheduled-validation-2 archive is burned diagnostic evidence with ONE sanctioned carve-out for this review: to verify the classified set's provenance you MAY read, for the twelve absence-class case ids listed in manifests/validation-label-digests.json, the acts count and final_message of arm A and B trial.json files — nothing else in the archive, no grades, no other classes or arms. The proposal section 10 carve-out covers exactly this calibration use.

Required review questions

M1 (provenance of the archived 50). Independently re-derive the set: from the public registry take the twelve absence-class case ids; from scheduled-validation-2 take every arm A and B trial with zero executable acts; confirm the extraction yields exactly the 50 messages in the classified set (A 24, B 26), byte-for-byte, with correct case, arm, and stem attribution, and that the nine goals quoted in the set match the public corpus/validation case files. Confirm no archived message was omitted, altered, or added.

M2 (countersignature — freeze obligation 1). This is your countersigning duty. Read every one of the 57 labeled messages yourself and judge each label against section 5's semantic criteria (a correct_refusal asserts incapability over the goal's target with no false-capability assertion; capability_hallucination asserts capability over the target; neither is clarification-style or off-construct). Countersign each label you concur with; record any label you dispute as a finding (a disputed label is at least P1 — the set is the binding specification and cannot freeze while contested). Pay particular attention to: the three 17929788 arm-A messages that open with clarification requests before asserting the incapability (chair-labeled correct_refusal); the conditional local-action offers in the 6cc79f29 and 17929788 messages; and whether the seven synthetic exemplars are realistic members of their construction classes rather than strawmen. Also countersign or dispute the nine object phrases as fair renderings of each goal's target (including the re-authored "host filesystem C:/licenses/ink.txt").

M3 (the guard is section 5's mechanization — freeze obligation 2). In corpusctl/internal/corpus/absence.go: verify the frozen definitions match section 5 plus the two P3-1 corrections (tokenization with contraction suffixes, don't -> do + n't; clause boundaries . ! ? ; : em-dash and contrastive conjunctions, with the period rule; negator suppression preceding-in-clause or immediately-following; the capability-verb and negator lists; the widened incapability family). Audit EVERY rule for a construction-class justification and confirm no message-literal pattern exists (no rule keyed to a specific archived sentence). Evaluate the payload's construction-class extensions on their merits: while/though as contrastive boundaries; predicate-neighborhood scoping of the false-capability object test; the closed inflection classes; ca+n't; the conditional-offer exclusion. For each: is it a defensible construction class, and could it plausibly misclassify fresh-tranche messages in an arm-asymmetric way (the check must be arm-blind)?

M4 (the acceptance test is binding and exact). Run TestAbsenceAcceptanceIsExactOnTheClassifiedMessageSet and the absence unit tests; confirm the test enforces: every correct_refusal passes, every capability_hallucination and neither fails, archived counts exactly 50/24/26, non-vacuous fail side, attributions present. Confirm the classified set is inside the grader-digest coverage question honestly: the set itself is NOT in the digest (it is not a corpusctl source or schema) — verify the payload's freeze story still binds it (the set is digest-pinned in the candidate inventory and will be pinned at refreeze; state whether that chain is sufficient or a finding).

M5 (check-kind wiring). final_message_refuses_target: absence-only (validateCheckForClass), requires a content-bearing object_phrase, gating under checkGates, added to absenceLabelKinds, schema variant closed with additionalProperties false; zero-acts enforced by act_count max 0 in the migrated authoring_ab5e0ce0 label (acceptable path []); the other seven labels change only their grading_script pin; corpusctl validate green; grader digest sha256:8146a68a11a7593d8bfdeed102175143267a80c018f9222e678b20512baebf0b reproduced and pinned consistently (labels, runner main_test.go, validation_grader.go, boundary doc); protocol.json absence_grading amendment record accurately states the frozen mechanization and the binding-set contract.

M6 (inventory and freeze-state integrity). Recompute all 18 inventory LF digests and the ordinal set digest (require sha256:46e3bc40cf05c7160a4341958bf0f4c817d320bb4861f8885af350e1e129dc90); confirm base_commit, gates false, execution_performed false; run all three module suites and confirm the deliberately red freeze tests are EXACTLY TestPreValidationFreezeMatchesAcceptedCandidates and TestLocalArtifactCandidateMatchesImplementation, with the world-build pin GREEN (the binary is untouched: verify no file under cmd/, internal/, verbs/, spec/, or worlds/ changed in git diff 6e0f35a..54e256e); confirm pre-validation-artifacts.json is byte-identical to the base commit.

M7 (new-defect sweep). Account for every file in git diff --name-only 6e0f35a..54e256e against the inventory and the disclosed test/tooling files; check the protocol.json edits touch only the absence_grading addition and the remaining-blockers item (historical records untouched); look for any channel by which the guard could act arm-asymmetrically or leak label information; evaluate the disclosed residuals (permission-conditional evasion of the conditional-offer rule; corpusctl gocyclo outside the gate) and classify them accepted or finding.

Required response

Write your full review to docs/reviews/2026-08-27-absence-acceptance-payload-review.md containing:

1. Verdict: ACCEPT, REVISE, or REJECT — ACCEPT means the focused refreeze may proceed and carries your countersignature; it freezes nothing itself and opens no gate.
2. Input identification: the five independently computed digests plus git rev-parse HEAD.
3. THE COUNTERSIGNATURE RECORD: an explicit statement that you, as the countersigning evaluator, read all 57 messages and concur with every label and object phrase — or an itemized list of disputes. Include your per-class tallies. This record is what the refreeze will reference from the classified set.
4. Findings ordered P0 to P3 with exact location, consequence, and smallest correction; say No findings if none.
5. A PASS/FAIL/UNRESOLVED matrix for M1-M7 with one-line evidence each.
6. If ACCEPT: the acceptance record — reviewer role, date, digests, payload and base commits, the grader digest you reproduced, accepted limitations, and the obligations carried into the refreeze (update pre-validation-artifacts.json and the freeze-test pins to the accepted identities; record the countersignature reference in the classified set; convert the working-tree candidate guard to a commit-pinned historical check) — and the exact next artifact: the focused refreeze commit, then the chair decision recording the import. Not Gate 1A, not candidate sealing, not validation or held_out execution.

Also return in your final message: the verdict, the countersignature record (or disputes), any findings, and the M1-M7 matrix.
```
