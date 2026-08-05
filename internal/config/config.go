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
		SuccessSound:  bundledSuccessSoundPath(),
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
	return filepath.Join(dir(), "sounds", "error.mp3")
}

func bundledSuccessSoundPath() string {
	return filepath.Join(dir(), "sounds", "normal.wav")
}

// extractBundledSound writes the embedded sounds to disk if they aren't there
// yet, so the binary works standalone on a machine that never built it.
func extractBundledSound() error {
	if err := os.MkdirAll(filepath.Join(dir(), "sounds"), 0o755); err != nil {
		return err
	}

	files := map[string][]byte{
		bundledFailureSoundPath(): assets.FailureSound,
		bundledSuccessSoundPath(): assets.SuccessSound,
	}

	for path, data := range files {
		if _, err := os.Stat(path); err == nil {
			continue
		}
		if err := os.WriteFile(path, data, 0o644); err != nil {
			return err
		}
	}

	return nil
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
