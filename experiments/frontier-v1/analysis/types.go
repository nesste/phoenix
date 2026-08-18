package main

type scheduledSummary struct {
	V              int                    `json:"v"`
	Tranche        string                 `json:"tranche"`
	ScheduleDigest string                 `json:"schedule_digest"`
	WorldBuild     string                 `json:"world_build"`
	Status         string                 `json:"status"`
	StopReason     string                 `json:"stop_reason"`
	Assigned       int                    `json:"assigned"`
	Results        []assignedTrialResult  `json:"results"`
	Configuration  scheduledConfiguration `json:"configuration"`
}

type scheduledConfiguration struct {
	RuntimeVersion     string `json:"runtime_version"`
	Model              string `json:"model"`
	GraderDigest       string `json:"grader_digest"`
	ArmBDocumentDigest string `json:"arm_b_document_digest"`
}

type assignedTrialResult struct {
	LaunchIndex    int             `json:"launch_index"`
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
	TrialPath      string          `json:"trial_path"`
	AssignmentPath string          `json:"assignment_path"`
}

type attemptRecord struct {
	RuntimePath string         `json:"runtime_path"`
	Metrics     runtimeMetrics `json:"metrics"`
}

type runtimeMetrics struct {
	Turns      int   `json:"turns"`
	WallTimeMS int64 `json:"wall_time_ms"`
	APITimeMS  int64 `json:"api_time_ms"`
}

type trancheManifest struct {
	V       int            `json:"v"`
	Tranche string         `json:"tranche"`
	Cases   []manifestCase `json:"cases"`
}

type manifestCase struct {
	CaseID   string `json:"case_id"`
	Class    string `json:"class"`
	FamilyID string `json:"family_id"`
}

type trialEvidence struct {
	Acts         []evidenceAct `json:"acts"`
	Orientations []orientation `json:"orientations"`
	FinalMessage string        `json:"final_message"`
}

type evidenceAct struct {
	Verb   string         `json:"verb"`
	Args   map[string]any `json:"args"`
	Status string         `json:"status"`
}

type orientation struct {
	Calls int `json:"calls"`
}

type observation struct {
	FamilyID       string
	CaseID         string
	Class          string
	Repetition     int
	Arm            string
	ITTSuccess     bool
	CompleteCase   bool
	Unresolved     bool
	BudgetStopped  bool
	SafetyStopped  bool
	CapHit         bool
	TimeoutOrCap   bool
	DeadEnd        bool
	HelpRequest    bool
	Tokens         float64
	CostUSD        float64
	Turns          float64
	WallTimeMS     float64
	APITimeMS      float64
	Acts           int
	WrongActs      int
	FrontierShown  int
	FrontierTaken  int
	FrontierPassed int
	Recovery       bool
}

type analysisReport struct {
	V                  int                  `json:"v"`
	Tranche            string               `json:"tranche"`
	AnalysisSeed       uint64               `json:"analysis_seed"`
	ScheduleDigest     string               `json:"schedule_digest"`
	WorldBuild         string               `json:"world_build"`
	RuntimeVersion     string               `json:"runtime_version"`
	Model              string               `json:"model"`
	GraderDigest       string               `json:"grader_digest"`
	ArmBDocument       string               `json:"arm_b_document_digest"`
	Status             string               `json:"status"`
	StatusReasons      []string             `json:"status_reasons"`
	IncompletePairings int                  `json:"incomplete_pairings"`
	Arms               map[string]armReport `json:"arms"`
	Claims             []claimReport        `json:"claims"`
	Measures           descriptiveMeasures  `json:"measures"`
}

type armReport struct {
	Assigned              int     `json:"assigned"`
	Successes             int     `json:"successes"`
	SuccessRate           float64 `json:"success_rate"`
	Unresolved            int     `json:"unresolved"`
	UnresolvedRate        float64 `json:"unresolved_rate"`
	CapHits               int     `json:"cap_hits"`
	CapHitRate            float64 `json:"cap_hit_rate"`
	TotalTokens           float64 `json:"total_tokens"`
	TotalCostUSD          float64 `json:"total_cost_usd"`
	TotalTurns            float64 `json:"total_turns"`
	WallTimeMS            float64 `json:"wall_time_ms"`
	MeanWallTimeMS        float64 `json:"mean_wall_time_ms"`
	MeanTokensToSuccess   float64 `json:"mean_tokens_to_success"`
	MeanFailedTrialTokens float64 `json:"mean_failed_trial_tokens"`
	MeanActsToSuccess     float64 `json:"mean_acts_to_success"`
	MeanFailedTrialActs   float64 `json:"mean_failed_trial_acts"`
}

type claimReport struct {
	ID                string            `json:"id"`
	Intervention      string            `json:"intervention"`
	Comparator        string            `json:"comparator"`
	Population        string            `json:"population"`
	Endpoint          string            `json:"endpoint"`
	Families          int               `json:"families"`
	Pairs             int               `json:"pairs"`
	Method            string            `json:"method"`
	PointEstimate     float64           `json:"point_estimate"`
	LowerBound        *float64          `json:"lower_bound,omitempty"`
	UpperBound        *float64          `json:"upper_bound,omitempty"`
	PWorse            *float64          `json:"p_worse,omitempty"`
	Decision          string            `json:"decision"`
	Reason            string            `json:"reason"`
	CostRatio         *ratioReport      `json:"paired_success_token_ratio,omitempty"`
	DeadEndDifference *inferenceReport  `json:"dead_end_difference,omitempty"`
	CapHitDifference  *inferenceReport  `json:"cap_hit_difference,omitempty"`
	DownstreamSuccess *inferenceReport  `json:"downstream_success,omitempty"`
	Sensitivity       sensitivityReport `json:"sensitivity"`
}

type inferenceReport struct {
	Point      float64  `json:"point_estimate"`
	LowerBound *float64 `json:"lower_bound,omitempty"`
	UpperBound *float64 `json:"upper_bound,omitempty"`
	PWorse     *float64 `json:"p_worse,omitempty"`
}

type ratioReport struct {
	Determinate bool     `json:"determinate"`
	Pairs       int      `json:"paired_successes"`
	Method      string   `json:"method,omitempty"`
	Ratio       *float64 `json:"ratio,omitempty"`
	UpperBound  *float64 `json:"upper_bound,omitempty"`
	Reason      string   `json:"reason,omitempty"`
}

type sensitivityReport struct {
	CompleteCase     *float64 `json:"complete_case,omitempty"`
	OneVotePerCase   float64  `json:"one_vote_per_case"`
	OneVotePerFamily float64  `json:"one_vote_per_family"`
}

type descriptiveMeasures struct {
	WrongVerbRate              float64  `json:"wrong_verb_rate"`
	DeadEndRate                float64  `json:"dead_end_rate"`
	TimeoutOrCapHitRate        float64  `json:"timeout_or_cap_hit_rate"`
	FrontierTakeRate           *float64 `json:"frontier_take_rate,omitempty"`
	FrontierTakeAndSucceedRate *float64 `json:"frontier_take_and_succeed_rate,omitempty"`
}
