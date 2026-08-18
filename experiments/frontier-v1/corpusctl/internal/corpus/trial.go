package corpus

import "fmt"

// Trial is captured evidence from one runner execution against one case. It
// carries no expected outcome; grading happens externally by comparing a
// Trial against a Label (see grade.go).
type Trial struct {
	V            int           `json:"v"`
	CaseID       string        `json:"case_id"`
	WorldBuild   string        `json:"world_build"`
	Acts         []Act         `json:"acts"`
	Orientations []Orientation `json:"orientations"`
	FinalMessage string        `json:"final_message"`
	EndState     EndState      `json:"end_state"`
}

type Orientation struct {
	Seq     int    `json:"seq"`
	Handle  string `json:"handle"`
	Intent  string `json:"intent"`
	Matched bool   `json:"matched"`
	Calls   int    `json:"calls"`
}

// Act is one structured step the trial took. Its position in Trial.Acts is
// the ordering used to match a label's acceptable_paths.
type Act struct {
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

// EndState captures observable sandbox facts without re-touching the
// filesystem at grade time.
type EndState struct {
	Files []FileState `json:"files"`
}

// FileState is one file's presence and content at the end of a trial.
type FileState struct {
	Path    string `json:"path"`
	Present bool   `json:"present"`
	Content string `json:"content,omitempty"`
}

// LoadTrial reads, canonicalizes, schema-validates, and decodes a trial
// record at root-relative path.
func LoadTrial(root, path string) (Document, Trial, error) {
	schemas, err := loadSchemas(root)
	if err != nil {
		return Document{}, Trial{}, err
	}
	document, err := loadDocument(root, path)
	if err != nil {
		return Document{}, Trial{}, err
	}
	if err := schemas.validate("trial", document); err != nil {
		return Document{}, Trial{}, err
	}
	var trial Trial
	if err := document.Into(&trial); err != nil {
		return Document{}, Trial{}, err
	}
	if err := validateTrial(trial); err != nil {
		return Document{}, Trial{}, fmt.Errorf("%s: %w", document.Path, err)
	}
	return document, trial, nil
}

func validateTrial(trial Trial) error {
	for index, orientation := range trial.Orientations {
		if orientation.Seq != index {
			return fmt.Errorf("orientation seq %d at index %d; sequence must be contiguous from zero", orientation.Seq, index)
		}
		if orientation.Matched != (orientation.Calls > 0) {
			return fmt.Errorf("orientation %d matched flag disagrees with its call count", index)
		}
	}
	for index, act := range trial.Acts {
		if act.Seq != index {
			return fmt.Errorf("act seq %d at index %d; sequence must be contiguous from zero", act.Seq, index)
		}
	}
	seenFiles := make(map[string]struct{}, len(trial.EndState.Files))
	for _, file := range trial.EndState.Files {
		if _, exists := seenFiles[file.Path]; exists {
			return fmt.Errorf("duplicate end_state file %s", file.Path)
		}
		seenFiles[file.Path] = struct{}{}
		if !file.Present && file.Content != "" {
			return fmt.Errorf("absent end_state file %s carries content", file.Path)
		}
	}
	return nil
}
