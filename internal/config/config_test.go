package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_CreatesDefaultConfigOnFirstRun(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	cfg := Load()

	if cfg != defaultConfig() {
		t.Errorf("Load() = %+v, want default %+v", cfg, defaultConfig())
	}

	if _, err := os.Stat(path()); err != nil {
		t.Errorf("expected config file to be written at %s: %v", path(), err)
	}
}

func TestLoad_ReadsExistingConfig(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	if err := os.MkdirAll(dir(), 0o755); err != nil {
		t.Fatal(err)
	}

	custom := `{
		"failureSound": "/tmp/custom-fail.mp3",
		"playOnFailure": true,
		"playOnSuccess": true
	}`
	if err := os.WriteFile(path(), []byte(custom), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := Load()

	if cfg.FailureSound != "/tmp/custom-fail.mp3" {
		t.Errorf("FailureSound = %q, want %q", cfg.FailureSound, "/tmp/custom-fail.mp3")
	}
	if !cfg.PlayOnSuccess {
		t.Error("PlayOnSuccess = false, want true (from custom config)")
	}
	// SuccessSound wasn't in the custom JSON, so it should keep its default.
	if cfg.SuccessSound != defaultConfig().SuccessSound {
		t.Errorf("SuccessSound = %q, want default %q", cfg.SuccessSound, defaultConfig().SuccessSound)
	}
}

func TestLoad_InvalidJSONFallsBackToDefault(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	if err := os.MkdirAll(dir(), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path(), []byte("{not valid json"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := Load()

	if cfg != defaultConfig() {
		t.Errorf("Load() with invalid JSON = %+v, want default %+v", cfg, defaultConfig())
	}
}

func TestDir_UsesHomeDotShellsound(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	want := filepath.Join(home, ".shellsound")
	if got := dir(); got != want {
		t.Errorf("dir() = %q, want %q", got, want)
	}
}
