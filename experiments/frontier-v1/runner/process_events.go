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
	"sync"
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
	// mu serializes appends: the run loop and the signal goroutine both write
	// this file, and an interleaved line would be unparseable at resume.
	mu sync.Mutex
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
	log.mu.Lock()
	defer log.mu.Unlock()
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
// process is signalled and then stops the process. The returned stop function
// releases the handler.
//
// Recording must not suppress the termination. A scheduled run that swallowed
// its supervisor's stop signal would contradict the section 9 host
// requirement for a supervised custodian process, and would leave the log
// asserting a termination that never happened while the run kept launching
// trials. The handler therefore records the entry durably and exits.
func watchProcessTermination(log *processEventLog, progress func() processEventProgress) func() {
	return watchProcessTerminationWith(log, progress, os.Exit)
}

// watchProcessTerminationWith is watchProcessTermination with an injectable
// exit, so the stop behaviour itself is testable.
func watchProcessTerminationWith(log *processEventLog, progress func() processEventProgress, exit func(int)) func() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	return watchProcessTerminationOn(log, progress, exit, signals, func() { signal.Stop(signals) })
}

// watchProcessTerminationOn is the seam the tests use: it takes the signal
// channel directly, so the stop behaviour can be exercised on a host whose
// os.Process.Signal does not support delivering an interrupt to itself.
func watchProcessTerminationOn(
	log *processEventLog,
	progress func() processEventProgress,
	exit func(int),
	signals chan os.Signal,
	release func(),
) func() {
	done := make(chan struct{})
	go func() {
		select {
		case received := <-signals:
			_ = log.recordTermination(progress(), received.String())
			release()
			exit(terminationExitCode(received))
		case <-done:
		}
	}()
	return func() {
		release()
		close(done)
	}
}

// terminationExitCode follows the shell convention of 128 plus the signal
// number, so a supervisor can tell an interrupted run from a failed one.
func terminationExitCode(received os.Signal) int {
	if received == syscall.SIGTERM {
		return 143
	}
	return 130
}

// liveProgress carries the current outcome-free progress to the termination
// handler, which runs on the signal goroutine while the run loop advances.
type liveProgress struct {
	mu    sync.Mutex
	value processEventProgress
}

func (live *liveProgress) set(value processEventProgress) {
	live.mu.Lock()
	live.value = value
	live.mu.Unlock()
}

func (live *liveProgress) get() processEventProgress {
	live.mu.Lock()
	defer live.mu.Unlock()
	return live.value
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
