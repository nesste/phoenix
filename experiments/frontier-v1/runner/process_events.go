package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"syscall"
	"time"
)

const scheduledProcessEventLogName = "process-events.jsonl"

// Protocol v5 section 9 process-event kinds. The log is the custodian record
// that a resume attestation cites when it classifies an interruption cause:
// without it a process kill cannot be established as participant-free and the
// tranche closes indeterminate.
const (
	processEventStart       = "start"
	processEventCheckpoint  = "checkpoint"
	processEventStop        = "stop"
	processEventTermination = "termination"
)

// processPrincipal records the originating principal where the host exposes
// it. A supervised custodian process exposes its own owner; an asynchronous
// signal does not carry the sending principal on either supported host, so
// the record says so explicitly rather than guessing.
type processPrincipal struct {
	Exposed bool   `json:"exposed"`
	Source  string `json:"source"`
	User    string `json:"user,omitempty"`
	Host    string `json:"host,omitempty"`
	PID     int    `json:"pid"`
	Reason  string `json:"reason,omitempty"`
}

type processEvent struct {
	V                    int              `json:"v"`
	Event                string           `json:"event"`
	RecordedAt           string           `json:"recorded_at"`
	Tranche              string           `json:"tranche"`
	ScheduleDigest       string           `json:"schedule_digest"`
	WorldBuild           string           `json:"world_build,omitempty"`
	NextLaunchIndex      int              `json:"next_launch_index"`
	CompletedPairingKeys int              `json:"completed_pairing_keys"`
	SpentUSD             float64          `json:"spent_usd"`
	Resumes              int              `json:"resumes"`
	Principal            processPrincipal `json:"principal"`
	Detail               string           `json:"detail,omitempty"`
}

// processEventLog appends outcome-free lifecycle events for one scheduled
// output directory. Every write is an append followed by a sync so an
// abrupt termination cannot lose the events already recorded.
type processEventLog struct {
	path           string
	tranche        string
	scheduleDigest string
	worldBuild     string
	clock          func() time.Time
}

func newProcessEventLog(outputDir, tranche, scheduleDigest, worldBuild string) *processEventLog {
	return &processEventLog{
		path:           filepath.Join(outputDir, scheduledProcessEventLogName),
		tranche:        tranche,
		scheduleDigest: scheduleDigest,
		worldBuild:     worldBuild,
		clock:          func() time.Time { return time.Now().UTC() },
	}
}

type processEventProgress struct {
	nextLaunchIndex      int
	completedPairingKeys int
	spentUSD             float64
	resumes              int
}

func (log *processEventLog) append(kind string, progress processEventProgress, principal processPrincipal, detail string) error {
	if log == nil {
		return nil
	}
	event := processEvent{
		V: 1, Event: kind, RecordedAt: log.clock().Format(time.RFC3339Nano),
		Tranche: log.tranche, ScheduleDigest: log.scheduleDigest, WorldBuild: log.worldBuild,
		NextLaunchIndex: progress.nextLaunchIndex, CompletedPairingKeys: progress.completedPairingKeys,
		SpentUSD: progress.spentUSD, Resumes: progress.resumes, Principal: principal, Detail: detail,
	}
	line, err := json.Marshal(event)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(log.path), 0o755); err != nil {
		return err
	}
	file, err := os.OpenFile(log.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	if _, err := file.Write(append(line, '\n')); err != nil {
		file.Close()
		return err
	}
	if err := file.Sync(); err != nil {
		file.Close()
		return err
	}
	return file.Close()
}

// record appends an event attributed to this custodian process.
func (log *processEventLog) record(kind string, progress processEventProgress, detail string) error {
	return log.append(kind, progress, hostProcessPrincipal(), detail)
}

// recordTermination appends the termination event for an asynchronous signal.
// Neither Windows nor Linux exposes the sending principal to a Go signal
// handler, so the entry states that the principal is unavailable; a resume
// after a kill therefore needs the custodian's written establishment.
func (log *processEventLog) recordTermination(progress processEventProgress, signalName string) error {
	principal := hostProcessPrincipal()
	principal.Exposed = false
	principal.Source = "signal_without_sender_identity"
	principal.Reason = "the host signal interface does not expose the sending principal"
	return log.append(processEventTermination, progress, principal, "received "+signalName)
}

func hostProcessPrincipal() processPrincipal {
	principal := processPrincipal{Exposed: true, Source: "host_process_owner", PID: os.Getpid()}
	if owner, err := user.Current(); err == nil {
		principal.User = owner.Username
	} else {
		principal.Exposed = false
		principal.Reason = "the host does not expose the process owner"
	}
	if host, err := os.Hostname(); err == nil {
		principal.Host = host
	}
	return principal
}

// watchProcessTermination records a termination event when the custodian
// process is signalled. The returned stop function releases the handler.
func watchProcessTermination(log *processEventLog, progress func() processEventProgress) func() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	done := make(chan struct{})
	go func() {
		select {
		case received := <-signals:
			_ = log.recordTermination(progress(), received.String())
		case <-done:
		}
	}()
	return func() {
		signal.Stop(signals)
		close(done)
	}
}

// digestProcessEventLog returns the raw SHA-256 of the process-event log,
// the identity the custody record commits and a resume attestation cites.
func digestProcessEventLog(outputDir string) (string, error) {
	contents, err := os.ReadFile(filepath.Join(outputDir, scheduledProcessEventLogName))
	if err != nil {
		return "", fmt.Errorf("read process-event log: %w", err)
	}
	sum := sha256.Sum256(contents)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}
