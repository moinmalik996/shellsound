package main

import (
	"testing"

	"github.com/moinabbas/shellsound/internal/config"
)

func TestSoundToPlay(t *testing.T) {
	tests := []struct {
		name     string
		exitCode int
		cfg      config.Config
		want     string
	}{
		{
			name:     "failure plays failure sound when enabled",
			exitCode: 1,
			cfg:      config.Config{FailureSound: "fail.mp3", PlayOnFailure: true},
			want:     "fail.mp3",
		},
		{
			name:     "failure stays silent when disabled",
			exitCode: 1,
			cfg:      config.Config{FailureSound: "fail.mp3", PlayOnFailure: false},
			want:     "",
		},
		{
			name:     "success stays silent by default",
			exitCode: 0,
			cfg:      config.Config{SuccessSound: "ok.mp3", PlayOnSuccess: false},
			want:     "",
		},
		{
			name:     "success plays success sound when enabled",
			exitCode: 0,
			cfg:      config.Config{SuccessSound: "ok.mp3", PlayOnSuccess: true},
			want:     "ok.mp3",
		},
		{
			name:     "nonzero exit code other than 1 still counts as failure",
			exitCode: 127,
			cfg:      config.Config{FailureSound: "fail.mp3", PlayOnFailure: true},
			want:     "fail.mp3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := soundToPlay(tt.exitCode, tt.cfg)
			if got != tt.want {
				t.Errorf("soundToPlay(%d, %+v) = %q, want %q", tt.exitCode, tt.cfg, got, tt.want)
			}
		})
	}
}

func TestRemoveDelimitedBlock(t *testing.T) {
	t.Run("removes block when markers exist", func(t *testing.T) {
		input := "line1\n# >>> shellsound hook >>>\na\nb\n# <<< shellsound hook <<<\nline2\n"
		got := removeDelimitedBlock(input, hookStart, hookEnd)
		want := "line1\nline2\n"
		if got != want {
			t.Errorf("removeDelimitedBlock() = %q, want %q", got, want)
		}
	})

	t.Run("no-op when markers missing", func(t *testing.T) {
		input := "line1\nline2\n"
		got := removeDelimitedBlock(input, hookStart, hookEnd)
		if got != input {
			t.Errorf("removeDelimitedBlock() should not modify content without markers")
		}
	})
}
