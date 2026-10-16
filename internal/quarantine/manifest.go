// Package quarantine tracks tests that are temporarily excluded from gating
// while their owners fix them.
package quarantine

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/gitstq/flakehunter/internal/engine"
)

// Entry is one quarantined test.
type Entry struct {
	Fingerprint string    `json:"fingerprint"`
	Suite       string    `json:"suite,omitempty"`
	Name        string    `json:"name,omitempty"`
	Reason      string    `json:"reason,omitempty"`
	Framework   string    `json:"framework,omitempty"`
	AddedAt     time.Time `json:"added_at"`
	StableRuns  int       `json:"stable_runs"`
}

// Manifest is the quarantine list stored at .flakehunter/quarantine.json.
type Manifest struct {
	Entries []Entry `json:"entries"`
}

// LoadManifest reads the quarantine manifest; a missing file yields an empty manifest.
func LoadManifest(dir string) (*Manifest, error) {
	data, err := os.ReadFile(filepath.Join(dir, "quarantine.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return &Manifest{}, nil
		}
		return nil, err
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// Save writes the manifest back to disk.
func (m *Manifest) Save(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "quarantine.json"), append(data, '\n'), 0o644)
}

// IsQuarantined reports whether fp is on the list.
func (m *Manifest) IsQuarantined(fp string) bool {
	for _, e := range m.Entries {
		if e.Fingerprint == fp {
			return true
		}
	}
	return false
}

// Add adds a new entry or refreshes an existing one.
func (m *Manifest) Add(fp, suite, name, reason, framework string) {
	for i := range m.Entries {
		if m.Entries[i].Fingerprint == fp {
			m.Entries[i].Reason = reason
			m.Entries[i].Framework = framework
			m.Entries[i].StableRuns = 0
			return
		}
	}
	m.Entries = append(m.Entries, Entry{
		Fingerprint: fp,
		Suite:       suite,
		Name:        name,
		Reason:      reason,
		Framework:   framework,
		AddedAt:     time.Now().UTC(),
	})
}

// Sync increments the stable-run counter of entries whose test is currently
// stable and drops the ones that reached the threshold. It returns the
// entries that were removed.
func (m *Manifest) Sync(status map[string]string, threshold int) []Entry {
	var kept []Entry
	var removed []Entry
	for _, e := range m.Entries {
		if status[e.Fingerprint] == engine.StatusStablePass {
			e.StableRuns++
		} else {
			e.StableRuns = 0
		}
		if e.StableRuns >= threshold {
			removed = append(removed, e)
			continue
		}
		kept = append(kept, e)
	}
	m.Entries = kept
	return removed
}
