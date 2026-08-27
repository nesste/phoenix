package surface

import "fmt"

// WireEnvelope is the compressed transport form of Envelope. Null and empty
// fields are omitted, and the full world_build digest appears only in the
// first response of a session; later responses carry the 8-character prefix
// in wb. The archival Envelope is recovered losslessly by Expand, which is
// committed as a frozen round-trip invariant: Expand(Compress(x)) == x.
type WireEnvelope struct {
	V                int             `json:"v"`
	WorldBuild       string          `json:"world_build,omitempty"`
	WorldBuildPrefix string          `json:"wb,omitempty"`
	SessionID        string          `json:"session_id"`
	ActID            string          `json:"act_id"`
	Handle           string          `json:"handle"`
	Verb             string          `json:"verb"`
	Status           Status          `json:"status"`
	Result           any             `json:"result,omitempty"`
	Error            *Error          `json:"error,omitempty"`
	Text             string          `json:"text"`
	Handles          *HandleDelta    `json:"handles,omitempty"`
	Frontier         []FrontierEntry `json:"frontier,omitempty"`
	Refusal          *Refusal        `json:"refusal,omitempty"`
}

// worldBuildPrefixLength is the frozen wb length: 8 hex characters of the
// digest, after the sha256: scheme.
const worldBuildPrefixLength = 8

func worldBuildPrefix(digest string) string {
	hexPart := digest
	if len(digest) > len("sha256:") && digest[:len("sha256:")] == "sha256:" {
		hexPart = digest[len("sha256:"):]
	}
	if len(hexPart) < worldBuildPrefixLength {
		return hexPart
	}
	return hexPart[:worldBuildPrefixLength]
}

// Compress produces the wire form of one envelope. The first response of a
// session announces the full world_build digest; every later response carries
// only the prefix.
func Compress(envelope Envelope, firstResponse bool) WireEnvelope {
	wire := WireEnvelope{
		V: envelope.V, SessionID: envelope.SessionID, ActID: envelope.ActID,
		Handle: envelope.Handle, Verb: envelope.Verb, Status: envelope.Status,
		Result: envelope.Result, Error: envelope.Error, Text: envelope.Text,
		Refusal: envelope.Refusal,
	}
	if firstResponse {
		wire.WorldBuild = envelope.WorldBuild
	} else {
		wire.WorldBuildPrefix = worldBuildPrefix(envelope.WorldBuild)
	}
	if len(envelope.Handles.Grant) > 0 || len(envelope.Handles.Revoke) > 0 {
		delta := envelope.Handles
		wire.Handles = &delta
	}
	if len(envelope.Frontier) > 0 {
		wire.Frontier = envelope.Frontier
	}
	return wire
}

// Expand reconstructs the archival envelope from its wire form. worldBuild is
// the session's full digest; it must match the wire prefix when the wire form
// carries one. Expansion is deterministic and lossless.
func Expand(wire WireEnvelope, worldBuild string) (Envelope, error) {
	envelope := Envelope{
		V: wire.V, SessionID: wire.SessionID, ActID: wire.ActID,
		Handle: wire.Handle, Verb: wire.Verb, Status: wire.Status,
		Result: wire.Result, Error: wire.Error, Text: wire.Text,
		Handles:  HandleDelta{Grant: []HandleGrant{}, Revoke: []string{}},
		Frontier: []FrontierEntry{}, Refusal: wire.Refusal,
	}
	switch {
	case wire.WorldBuild != "":
		envelope.WorldBuild = wire.WorldBuild
	case wire.WorldBuildPrefix != "":
		if worldBuildPrefix(worldBuild) != wire.WorldBuildPrefix {
			return Envelope{}, fmt.Errorf("wire world-build prefix %q does not match session digest", wire.WorldBuildPrefix)
		}
		envelope.WorldBuild = worldBuild
	default:
		return Envelope{}, fmt.Errorf("wire envelope carries neither world_build nor wb")
	}
	if wire.Handles != nil {
		envelope.Handles = *wire.Handles
	}
	if wire.Frontier != nil {
		envelope.Frontier = wire.Frontier
	}
	return envelope, nil
}

// firstWireResponse reports whether this session has announced its full
// world_build digest yet and marks it announced.
func (session *Session) firstWireResponse() bool {
	session.mu.Lock()
	defer session.mu.Unlock()
	if session.announced {
		return false
	}
	session.announced = true
	return true
}
