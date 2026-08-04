package audio

import "testing"

func TestPlay_EmptyPathDoesNothing(t *testing.T) {
	orig := runCommand
	defer func() { runCommand = orig }()

	called := false
	runCommand = func(name string, arg ...string) error {
		called = true
		return nil
	}

	Play("")

	if called {
		t.Error("Play(\"\") should not invoke the player command")
	}
}

func TestPlay_InvokesAfplayWithPath(t *testing.T) {
	orig := runCommand
	defer func() { runCommand = orig }()

	var gotName string
	var gotArgs []string
	runCommand = func(name string, arg ...string) error {
		gotName = name
		gotArgs = arg
		return nil
	}

	Play("/tmp/shellsound.mp3")

	if gotName != "afplay" {
		t.Errorf("player command = %q, want %q", gotName, "afplay")
	}
	if len(gotArgs) != 1 || gotArgs[0] != "/tmp/shellsound.mp3" {
		t.Errorf("player args = %v, want [%q]", gotArgs, "/tmp/shellsound.mp3")
	}
}
