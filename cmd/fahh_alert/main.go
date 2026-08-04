// Command fahh_alert is invoked by a zsh precmd hook after every command,
// and plays a sound when the previous command's exit code was non-zero.
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/moinabbas/fahh_alert/internal/audio"
	"github.com/moinabbas/fahh_alert/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: fahh_alert <exit_code>")
		os.Exit(1)
	}

	exitCode, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "fahh_alert: invalid exit code %q\n", os.Args[1])
		os.Exit(1)
	}

	cfg := config.Load()
	audio.Play(soundToPlay(exitCode, cfg))
}

// soundToPlay returns which sound (if any) should play for the given
// exit code under cfg. An empty string means stay silent.
func soundToPlay(exitCode int, cfg config.Config) string {
	switch {
	case exitCode != 0 && cfg.PlayOnFailure:
		return cfg.FailureSound
	case exitCode == 0 && cfg.PlayOnSuccess:
		return cfg.SuccessSound
	default:
		return ""
	}
}
