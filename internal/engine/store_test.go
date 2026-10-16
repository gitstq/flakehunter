package engine

import (
	"testing"
	"time"

	"github.com/gitstq/flakehunter/internal/model"
)

func TestStoreAppendAndLoad(t *testing.T) {
	dir := t.TempDir()
	store, err := OpenStore(dir, 10)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		run := model.Run{
			RunID:     string(rune('a' + i)),
			Timestamp: time.Date(2026, 10, 16, 10, i, 0, 0, time.UTC),
			Results:   []model.TestResult{{Fingerprint: "s#t", Outcome: model.OutcomePassed}},
		}
		if err := store.Append(run); err != nil {
			t.Fatalf("append %d: %v", i, err)
		}
	}
	runs, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(runs) != 3 {
		t.Fatalf("want 3 runs, got %d", len(runs))
	}
	if runs[0].RunID != "a" || runs[2].RunID != "c" {
		t.Errorf("runs out of order: %s ... %s", runs[0].RunID, runs[2].RunID)
	}
}

func TestStoreRotation(t *testing.T) {
	dir := t.TempDir()
	store, err := OpenStore(dir, 2)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 5; i++ {
		run := model.Run{
			RunID:     string(rune('a' + i)),
			Timestamp: time.Date(2026, 10, 16, 10, i, 0, 0, time.UTC),
			Results:   []model.TestResult{{Fingerprint: "s#t", Outcome: model.OutcomePassed}},
		}
		if err := store.Append(run); err != nil {
			t.Fatalf("append %d: %v", i, err)
		}
	}
	runs, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(runs) != 2 {
		t.Fatalf("want 2 runs after rotation, got %d", len(runs))
	}
	if runs[0].RunID != "d" || runs[1].RunID != "e" {
		t.Errorf("rotation kept wrong runs: %s, %s", runs[0].RunID, runs[1].RunID)
	}
}

func TestStoreLoadMissing(t *testing.T) {
	store, err := OpenStore(t.TempDir(), 10)
	if err != nil {
		t.Fatal(err)
	}
	runs, err := store.Load()
	if err != nil {
		t.Fatalf("missing history should not error: %v", err)
	}
	if len(runs) != 0 {
		t.Errorf("want 0 runs, got %d", len(runs))
	}
}

func TestConfigDefaultsAndOverride(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.MinRuns != 3 || cfg.MaxRuns != 100 || cfg.MaxFlakeRate != 5 {
		t.Errorf("unexpected defaults: %+v", cfg)
	}
	dir := t.TempDir()
	cfg.MinRuns = 5
	if err := SaveConfig(dir, cfg); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadConfig(dir)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.MinRuns != 5 {
		t.Errorf("want MinRuns=5 after reload, got %d", loaded.MinRuns)
	}
	// Missing config falls back to defaults.
	fallback, err := LoadConfig(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if fallback.MinRuns != 3 {
		t.Errorf("want default MinRuns=3, got %d", fallback.MinRuns)
	}
}
