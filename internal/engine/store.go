package engine

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/gitstq/flakehunter/internal/model"
)

// Store is an append-only JSONL history of test runs kept entirely on disk.
type Store struct {
	Dir     string
	MaxRuns int
}

// OpenStore opens (creating when needed) the history directory.
func OpenStore(dir string, maxRuns int) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	if maxRuns <= 0 {
		maxRuns = 100
	}
	return &Store{Dir: dir, MaxRuns: maxRuns}, nil
}

// historyPath returns the JSONL history file location.
func (s *Store) historyPath() string {
	return filepath.Join(s.Dir, "history.jsonl")
}

// Append records a run and rotates history to keep at most MaxRuns entries.
func (s *Store) Append(run model.Run) error {
	line, err := json.Marshal(run)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(s.historyPath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	if _, err := f.Write(append(line, '\n')); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return s.rotate()
}

// rotate rewrites the history file keeping only the most recent MaxRuns runs.
func (s *Store) rotate() error {
	runs, err := s.Load()
	if err != nil {
		return err
	}
	if len(runs) <= s.MaxRuns {
		return nil
	}
	runs = runs[len(runs)-s.MaxRuns:]
	tmp := s.historyPath() + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	for _, run := range runs {
		line, err := json.Marshal(run)
		if err != nil {
			f.Close()
			return err
		}
		if _, err := w.Write(append(line, '\n')); err != nil {
			f.Close()
			return err
		}
	}
	if err := w.Flush(); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, s.historyPath())
}

// Load reads the full run history, oldest first. Corrupt lines are skipped.
func (s *Store) Load() ([]model.Run, error) {
	f, err := os.Open(s.historyPath())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	var runs []model.Run
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 64*1024*1024)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		var run model.Run
		if err := json.Unmarshal([]byte(text), &run); err != nil {
			continue
		}
		runs = append(runs, run)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return runs, nil
}
