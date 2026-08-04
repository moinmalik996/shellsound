package audio

import "os/exec"

// runCommand starts name with arg and returns immediately; swapped out in tests.
var runCommand = func(name string, arg ...string) error {
	return exec.Command(name, arg...).Start()
}

// Play starts the given sound file playing via afplay without blocking
// the caller. It is a no-op if path is empty or afplay is unavailable.
func Play(path string) {
	if path == "" {
		return
	}

	_ = runCommand("afplay", path)
}
