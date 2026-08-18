package main

import "time"

type runnableCase struct {
	CaseID         string        `json:"case_id"`
	Goal           string        `json:"goal"`
	SandboxFixture string        `json:"sandbox_fixture"`
	WorldRef       string        `json:"world_ref"`
	FamilyID       string        `json:"family_id"`
	StateChanges   []stateChange `json:"state_changes"`
}

type stateChange struct {
	AfterAct int    `json:"after_act"`
	Path     string `json:"path"`
	Content  string `json:"content"`
}

type fixture struct {
	V           int               `json:"v"`
	FixtureID   string            `json:"fixture_id"`
	FamilyID    string            `json:"family_id"`
	Template    string            `json:"template"`
	Files       map[string]string `json:"files"`
	Uncommitted []string          `json:"uncommitted"`
}

type trial struct {
	V            int           `json:"v"`
	CaseID       string        `json:"case_id"`
	WorldBuild   string        `json:"world_build"`
	Acts         []act         `json:"acts"`
	Orientations []orientation `json:"orientations"`
	FinalMessage string        `json:"final_message"`
	EndState     endState      `json:"end_state"`
}

type orientation struct {
	Seq     int    `json:"seq"`
	Handle  string `json:"handle"`
	Intent  string `json:"intent"`
	Matched bool   `json:"matched"`
	Calls   int    `json:"calls"`
}

type act struct {
	Seq        int            `json:"seq"`
	Handle     string         `json:"handle"`
	HandleType string         `json:"handle_type"`
	Verb       string         `json:"verb"`
	Args       map[string]any `json:"args"`
	Status     string         `json:"status"`
	Command    []string       `json:"command,omitempty"`
	ExitCode   *int           `json:"exit_code,omitempty"`
	Output     string         `json:"output,omitempty"`
}

type endState struct {
	Files []fileState `json:"files"`
}

type fileState struct {
	Path    string `json:"path"`
	Present bool   `json:"present"`
	Content string `json:"content,omitempty"`
}

type runtimeRequest struct {
	Sandbox         string
	Goal            string
	Roots           map[string]string
	PhoenixPath     string
	WorldPath       string
	SchemaPath      string
	EpisodePath     string
	WorldBuild      string
	StateEventsPath string
	Timeout         time.Duration
	BudgetUSD       string
}

type runtimeResult struct {
	FinalMessage string
	RawOutput    []byte
}

type runtimeDriver interface {
	Verify() error
	Run(runtimeRequest) (runtimeResult, error)
}

type gradeDriver interface {
	Grade(repositoryRoot, labelPath, trialPath string) ([]byte, error)
}

type runConfig struct {
	repositoryRoot string
	outputDir      string
	worldPath      string
	schemaPath     string
	phoenixPath    string
	worldBuild     string
	timeout        time.Duration
	budgetUSD      string
}

type caseResult struct {
	CaseID    string `json:"case_id"`
	TrialPath string `json:"trial_path"`
	GradePath string `json:"grade_path"`
	Status    string `json:"status"`
}

type authoringSummary struct {
	V          int          `json:"v"`
	Tranche    string       `json:"tranche"`
	WorldBuild string       `json:"world_build"`
	Cases      []caseResult `json:"cases"`
	Passed     int          `json:"passed"`
	Failed     int          `json:"failed"`
}
