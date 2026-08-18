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
	Arm             string
	ArmBDocument    string
	FlatToolNames   []string
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
	FinalMessage    string
	RawOutput       []byte
	FirstModelToken bool
	Metrics         runtimeMetrics
	Failure         *runtimeFailure
}

type runtimeMetrics struct {
	InputTokens              int64   `json:"input_tokens"`
	CacheCreationInputTokens int64   `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int64   `json:"cache_read_input_tokens"`
	OutputTokens             int64   `json:"output_tokens"`
	TotalTokens              int64   `json:"total_tokens"`
	CostUSD                  float64 `json:"cost_usd"`
	Turns                    int     `json:"turns"`
	WallTimeMS               int64   `json:"wall_time_ms"`
	APITimeMS                int64   `json:"api_time_ms"`
}

type runtimeFailure struct {
	Code                  string `json:"code"`
	Message               string `json:"message"`
	BeforeFirstModelToken bool   `json:"before_first_model_token"`
	RetryEligible         bool   `json:"retry_eligible"`
	CapHit                bool   `json:"cap_hit"`
}

type runtimeDriver interface {
	Verify() error
	Run(runtimeRequest) (runtimeResult, error)
}

type gradeDriver interface {
	Grade(repositoryRoot, labelPath, trialPath string) ([]byte, error)
}

type runConfig struct {
	arm                string
	armBDocument       string
	armBDocumentDigest string
	flatToolNames      []string
	graderDigest       string
	repositoryRoot     string
	outputDir          string
	worldPath          string
	schemaPath         string
	phoenixPath        string
	worldBuild         string
	timeout            time.Duration
	budgetUSD          string
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
	Arm        string       `json:"arm"`
	WorldBuild string       `json:"world_build"`
	Cases      []caseResult `json:"cases"`
	Passed     int          `json:"passed"`
	Failed     int          `json:"failed"`
}

type attemptRecord struct {
	Attempt         int             `json:"attempt"`
	Status          string          `json:"status"`
	FirstModelToken bool            `json:"first_model_token"`
	RuntimePath     string          `json:"runtime_path"`
	Metrics         runtimeMetrics  `json:"metrics"`
	Failure         *runtimeFailure `json:"failure,omitempty"`
}

type assignedTrialResult struct {
	V              int             `json:"v"`
	LaunchIndex    int             `json:"launch_index"`
	FamilyBlock    int             `json:"family_block"`
	PairingIndex   int             `json:"pairing_index"`
	FamilyID       string          `json:"family_id"`
	CaseID         string          `json:"case_id"`
	Repetition     int             `json:"repetition"`
	Arm            string          `json:"arm"`
	Status         string          `json:"status"`
	ITTSuccess     bool            `json:"itt_success"`
	Termination    string          `json:"termination"`
	CapHit         bool            `json:"cap_hit"`
	Attempts       []attemptRecord `json:"attempts"`
	TotalCostUSD   float64         `json:"total_cost_usd"`
	TotalTokens    int64           `json:"total_tokens"`
	TrialPath      string          `json:"trial_path,omitempty"`
	GradePath      string          `json:"grade_path,omitempty"`
	AssignmentPath string          `json:"assignment_path,omitempty"`
}

type scheduledSummary struct {
	V              int                    `json:"v"`
	Tranche        string                 `json:"tranche"`
	ScheduleDigest string                 `json:"schedule_digest"`
	WorldBuild     string                 `json:"world_build"`
	Configuration  scheduledConfiguration `json:"configuration"`
	Status         string                 `json:"status"`
	StopReason     string                 `json:"stop_reason,omitempty"`
	RunBudgetUSD   float64                `json:"run_budget_usd"`
	SpentUSD       float64                `json:"spent_usd"`
	Assigned       int                    `json:"assigned"`
	Launched       int                    `json:"launched"`
	Passes         int                    `json:"passes"`
	Failures       int                    `json:"failures"`
	Unresolved     int                    `json:"unresolved"`
	BudgetStopped  int                    `json:"budget_stopped"`
	SafetyStopped  int                    `json:"safety_stopped"`
	CapHits        int                    `json:"cap_hits"`
	Results        []assignedTrialResult  `json:"results"`
}

type scheduledConfiguration struct {
	Runtime                    string              `json:"runtime"`
	RuntimeVersion             string              `json:"runtime_version"`
	Model                      string              `json:"model"`
	Effort                     string              `json:"effort"`
	ServiceTier                string              `json:"service_tier"`
	AccessMode                 string              `json:"access_mode"`
	Transport                  string              `json:"transport"`
	OutputFormat               string              `json:"output_format"`
	BuiltinTools               string              `json:"builtin_tools"`
	StrictMCPConfig            bool                `json:"strict_mcp_config"`
	SessionPersistence         bool                `json:"session_persistence"`
	ModelSeedSupport           bool                `json:"model_seed_support"`
	MaxTurns                   int                 `json:"max_turns"`
	TimeoutSeconds             int                 `json:"timeout_seconds"`
	MaxCostUSDPerTrial         float64             `json:"max_cost_usd_per_trial"`
	MaximumInfrastructureRetry int                 `json:"maximum_infrastructure_retries"`
	ScheduleSeed               uint64              `json:"schedule_seed"`
	Repetitions                int                 `json:"repetitions"`
	SystemPrompts              map[string]string   `json:"system_prompts"`
	AllowedTools               map[string][]string `json:"allowed_tools"`
	ArmBDocument               string              `json:"arm_b_document"`
	ArmBDocumentDigest         string              `json:"arm_b_document_digest"`
	GraderDigest               string              `json:"grader_digest"`
	TokenAccounting            string              `json:"token_accounting"`
}
