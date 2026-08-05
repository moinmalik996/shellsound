// Command shellsound is invoked by a zsh precmd hook after every command,
// and plays a sound when the previous command's exit code was non-zero.
package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/moinabbas/shellsound/internal/audio"
	"github.com/moinabbas/shellsound/internal/config"
)

const (
	hookStart = "# >>> shellsound hook >>>"
	hookEnd   = "# <<< shellsound hook <<<"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	if os.Args[1] == "uninstall" {
		if err := runUninstall(); err != nil {
			fmt.Fprintf(os.Stderr, "shellsound uninstall: %v\n", err)
			os.Exit(1)
		}
		return
	}

	exitCode, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "shellsound: invalid argument %q\n", os.Args[1])
		printUsage()
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

func printUsage() {
	fmt.Fprintln(os.Stderr, "usage:")
	fmt.Fprintln(os.Stderr, "  shellsound <exit_code>")
	fmt.Fprintln(os.Stderr, "  shellsound uninstall")
}

func runUninstall() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	zshrc := filepath.Join(home, ".zshrc")
	installDir := filepath.Join(home, ".shellsound")

	if err := removeHookBlock(zshrc, hookStart, hookEnd); err != nil {
		return err
	}

	if err := os.RemoveAll(installDir); err != nil {
		return err
	}

	fmt.Printf("Removed install directory: %s\n", installDir)
	fmt.Println("Uninstall complete. Reload your shell with: exec zsh")
	return nil
}

func removeHookBlock(zshrc, startMarker, endMarker string) error {
	data, err := os.ReadFile(zshrc)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			fmt.Printf("No %s found. Skipping hook cleanup.\n", zshrc)
			return nil
		}
		return err
	}

	content := string(data)
	updated := removeDelimitedBlock(content, startMarker, endMarker)
	if updated == content {
		fmt.Printf("No shellsound hook block found in %s.\n", zshrc)
		return nil
	}

	mode := fs.FileMode(0o644)
	if stat, statErr := os.Stat(zshrc); statErr == nil {
		mode = stat.Mode().Perm()
	}

	backup := fmt.Sprintf("%s.backup.%s", zshrc, time.Now().Format("20060102-150405"))
	if err := os.WriteFile(backup, data, mode); err != nil {
		return err
	}

	if err := os.WriteFile(zshrc, []byte(updated), mode); err != nil {
		return err
	}

	fmt.Printf("Backed up %s to %s\n", zshrc, backup)
	fmt.Printf("Removed shellsound hook block from %s\n", zshrc)
	return nil
}

func removeDelimitedBlock(content, startMarker, endMarker string) string {
	pattern := "(?s)\\n?" + regexp.QuoteMeta(startMarker) + ".*?" + regexp.QuoteMeta(endMarker) + "\\n?"
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(content, "\n")
}
