package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/moinabbas/shellsound/internal/assets"
)

// Config controls which sounds play and when.
type Config struct {
	FailureSound  string `json:"failureSound"`
	SuccessSound  string `json:"successSound"`
	PlayOnFailure bool   `json:"playOnFailure"`
	PlayOnSuccess bool   `json:"playOnSuccess"`
}

func defaultConfig() Config {
	return Config{
		FailureSound:  bundledFailureSoundPath(),
		SuccessSound:  "/System/Library/Sounds/Glass.aiff",
		PlayOnFailure: true,
		PlayOnSuccess: true,
	}
}

func dir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".shellsound"
	}
	return filepath.Join(home, ".shellsound")
}

func path() string {
	return filepath.Join(dir(), "config.json")
}

// bundledFailureSoundPath returns where the sound embedded in the binary
// gets extracted to on disk, since afplay needs a real file path to play.
func bundledFailureSoundPath() string {
	return filepath.Join(dir(), "sounds", "shellsound.mp3")
}

// extractBundledSound writes the embedded sound to disk if it isn't there
// yet, so the binary works standalone on a machine that never built it.
func extractBundledSound() error {
	path := bundledFailureSoundPath()
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, assets.FailureSound, 0o644)
}

// Load reads the config file, writing out the defaults on first run.
func Load() Config {
	_ = extractBundledSound()

	data, err := os.ReadFile(path())
	if err != nil {
		cfg := defaultConfig()
		_ = save(cfg)
		return cfg
	}

	cfg := defaultConfig()
	if err := json.Unmarshal(data, &cfg); err != nil {
		return defaultConfig()
	}

	return cfg
}

func save(cfg Config) error {
	if err := os.MkdirAll(dir(), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path(), data, 0o644)
}
