// Package world loads immutable Phoenix world definitions and creates isolated
// session graphs from them.
package world

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gowebpki/jcs"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

type Definition struct {
	V           int                   `json:"v"`
	ID          string                `json:"id"`
	Roots       []Root                `json:"roots"`
	HandleTypes map[string]HandleType `json:"handle_types"`
	Transitions []Transition          `json:"transitions"`
	digest      string
}

type Root struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Label    string   `json:"label"`
	Resource Resource `json:"resource"`
}

type Resource struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

type HandleType struct {
	Verbs map[string]Verb `json:"verbs"`
}

type Verb struct {
	ArgsSchema   json.RawMessage `json:"args_schema"`
	ResultSchema json.RawMessage `json:"result_schema"`
	Grants       []string        `json:"grants"`
	MayRevoke    bool            `json:"may_revoke"`
	Refusals     []RefusalRule   `json:"refusals"`
}

type RefusalRule struct {
	ID            string          `json:"id"`
	When          json.RawMessage `json:"when"`
	What          string          `json:"what"`
	Why           string          `json:"why"`
	Instead       *CallTemplate   `json:"instead"`
	NoAlternative string          `json:"no_alternative,omitempty"`
}

type Transition struct {
	ID          string       `json:"id"`
	Match       Match        `json:"match"`
	Suggestions []Suggestion `json:"suggestions"`
}

type Match struct {
	HandleType string          `json:"handle_type"`
	Verb       string          `json:"verb"`
	Status     string          `json:"status"`
	ResultWhen json.RawMessage `json:"result_when"`
}

type Suggestion struct {
	Call  CallTemplate `json:"call"`
	Why   string       `json:"why"`
	Score float64      `json:"score"`
}

type CallTemplate struct {
	Handle HandleSelector     `json:"handle"`
	Verb   string             `json:"verb"`
	Args   map[string]Binding `json:"args"`
	State  *Binding           `json:"state,omitempty"`
}

type HandleSelector struct {
	Source        string `json:"source"`
	Name          string `json:"name,omitempty"`
	ResultPointer string `json:"result_pointer,omitempty"`
}

type Binding struct {
	Literal       json.RawMessage `json:"literal,omitempty"`
	ResultPointer *string         `json:"result_pointer,omitempty"`
	StatePointer  *string         `json:"state_pointer,omitempty"`
	StateDigest   bool            `json:"state_digest,omitempty"`
}

func Load(schemaPath, definitionPath string) (*Definition, error) {
	contents, err := os.ReadFile(definitionPath)
	if err != nil {
		return nil, fmt.Errorf("read world definition: %w", err)
	}
	canonical, err := jcs.Transform(contents)
	if err != nil {
		return nil, fmt.Errorf("canonicalize world definition: %w", err)
	}

	compiler := jsonschema.NewCompiler()
	schema, err := compiler.Compile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("compile world schema: %w", err)
	}
	instance, err := jsonschema.UnmarshalJSON(bytes.NewReader(canonical))
	if err != nil {
		return nil, fmt.Errorf("decode world definition: %w", err)
	}
	if err := schema.Validate(instance); err != nil {
		return nil, fmt.Errorf("world schema validation failed: %w", err)
	}

	var definition Definition
	if err := json.Unmarshal(canonical, &definition); err != nil {
		return nil, fmt.Errorf("decode typed world definition: %w", err)
	}
	definition.digest = digestBytes(canonical)
	if err := ValidateDefinition(&definition); err != nil {
		return nil, fmt.Errorf("world semantic validation failed: %w", err)
	}
	return &definition, nil
}

func (definition *Definition) Digest() string {
	if definition == nil {
		return ""
	}
	return definition.digest
}

func digestValue(value any) (string, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	canonical, err := jcs.Transform(encoded)
	if err != nil {
		return "", err
	}
	return digestBytes(canonical), nil
}

func digestBytes(contents []byte) string {
	sum := sha256.Sum256(contents)
	return "sha256:" + hex.EncodeToString(sum[:])
}
