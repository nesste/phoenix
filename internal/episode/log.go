// Package episode records append-only, reconstructable Phoenix sessions.
package episode

import "time"

type EpisodeStatus string

const (
	EpisodeActive      EpisodeStatus = "active"
	EpisodeCompleted   EpisodeStatus = "completed"
	EpisodeInterrupted EpisodeStatus = "interrupted"
)

type ActStatus string

const (
	ActStarted     ActStatus = "started"
	ActCompleted   ActStatus = "completed"
	ActInterrupted ActStatus = "interrupted"
)

type ActRequest struct {
	Handle string         `json:"handle"`
	Verb   string         `json:"verb"`
	Args   map[string]any `json:"args"`
	State  string         `json:"state,omitempty"`
	Intent string         `json:"intent,omitempty"`
}

type SuggestionLink struct {
	ActID string `json:"act_id"`
	Index int    `json:"index"`
}

type StartedAct struct {
	EpisodeID       string
	ActID           string
	WorldBuild      string
	Request         ActRequest
	SuggestionTaken *SuggestionLink
}

type CompletedAct struct {
	EpisodeID string
	ActID     string
	Result    any
}

type ActRecord struct {
	Sequence        int64           `json:"sequence"`
	ActID           string          `json:"act_id"`
	WorldBuild      string          `json:"world_build"`
	Status          ActStatus       `json:"status"`
	Request         ActRequest      `json:"request"`
	RequestDigest   string          `json:"request_digest"`
	Result          any             `json:"result"`
	ResultDigest    string          `json:"result_digest,omitempty"`
	SuggestionTaken *SuggestionLink `json:"suggestion_taken"`
	StartedAt       time.Time       `json:"started_at"`
	CompletedAt     *time.Time      `json:"completed_at"`
}

type Record struct {
	V          int           `json:"v"`
	EpisodeID  string        `json:"episode_id"`
	SessionID  string        `json:"session_id"`
	WorldBuild string        `json:"world_build"`
	Status     EpisodeStatus `json:"status"`
	Outcome    string        `json:"outcome,omitempty"`
	StartedAt  time.Time     `json:"started_at"`
	EndedAt    *time.Time    `json:"ended_at"`
	Acts       []ActRecord   `json:"acts"`
}

type Pointer struct {
	EpisodeID string `json:"episode_id"`
	ActID     string `json:"act_id"`
}
