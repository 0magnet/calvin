package commands

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/0magnet/calvin"
)

// run executes the command with the given arguments and the given stdin,
// returning what it printed. stdin of nil means "a terminal": /dev/null is a
// character device, which is what the command uses to tell a pipe from a tty.
func run(t *testing.T, stdin []byte, args ...string) (string, error) {
	t.Helper()

	realIn, realOut := os.Stdin, os.Stdout
	t.Cleanup(func() { os.Stdin, os.Stdout = realIn, realOut })

	if stdin == nil {
		devNull, err := os.Open(os.DevNull)
		if err != nil {
			t.Fatalf("opening %s: %v", os.DevNull, err)
		}
		defer devNull.Close() //nolint:errcheck
		os.Stdin = devNull
	} else {
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("pipe: %v", err)
		}
		go func() {
			_, _ = w.Write(stdin) //nolint:errcheck
			_ = w.Close()         //nolint:errcheck
		}()
		defer r.Close() //nolint:errcheck
		os.Stdin = r
	}

	outR, outW, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = outW

	done := make(chan string, 1)
	go func() {
		b, _ := io.ReadAll(outR) //nolint:errcheck
		done <- string(b)
	}()

	RootCmd.SetArgs(args)
	runErr := RootCmd.RunE(RootCmd, args)

	_ = outW.Close() //nolint:errcheck
	return <-done, runErr
}

// Arguments have to win. Checking stdin first meant that whenever it was not a
// terminal — a script, a Makefile, CI — calvin ignored its arguments and
// blocked waiting for an EOF that never came.
func TestArgumentsWinOverStdin(t *testing.T) {
	out, err := run(t, []byte("from-stdin\n"), "hi")
	if err != nil {
		t.Fatalf("RunE: %v", err)
	}
	if want := calvin.AsciiFont("hi"); !strings.Contains(out, want) {
		t.Errorf("the argument was not rendered:\n%s", out)
	}
	if strings.Contains(out, calvin.AsciiFont("from-stdin")) {
		t.Error("stdin was rendered even though an argument was given")
	}
}

// The case that hung: a pipe with nothing in it, and an argument to render.
func TestAnArgumentIsRenderedWithAnEmptyPipeOnStdin(t *testing.T) {
	out, err := run(t, []byte{}, "hi")
	if err != nil {
		t.Fatalf("RunE: %v", err)
	}
	if want := calvin.AsciiFont("hi"); !strings.Contains(out, want) {
		t.Errorf("the argument was not rendered:\n%s", out)
	}
}

func TestSeveralArgumentsAreJoinedWithSpaces(t *testing.T) {
	out, err := run(t, nil, "a", "b")
	if err != nil {
		t.Fatalf("RunE: %v", err)
	}
	if want := calvin.AsciiFont("a b"); !strings.Contains(out, want) {
		t.Errorf("arguments were not joined with a space:\n%s", out)
	}
}

func TestStdinIsReadWhenThereAreNoArguments(t *testing.T) {
	out, err := run(t, []byte("hi\n"))
	if err != nil {
		t.Fatalf("RunE: %v", err)
	}
	if want := calvin.AsciiFont("hi"); !strings.Contains(out, want) {
		t.Errorf("piped input was not rendered:\n%s", out)
	}
}

func TestSeveralPipedLinesEachRender(t *testing.T) {
	out, err := run(t, []byte("ab\ncd\n"))
	if err != nil {
		t.Fatalf("RunE: %v", err)
	}
	for _, line := range []string{"ab", "cd"} {
		if !strings.Contains(out, calvin.AsciiFont(line)) {
			t.Errorf("%q was not rendered:\n%s", line, out)
		}
	}
}

// With no arguments and a terminal on stdin there is nothing to render, and
// blocking on a read the user never intended is the wrong answer.
func TestNoArgumentsAndATerminalIsAnError(t *testing.T) {
	_, err := run(t, nil)
	if err == nil {
		t.Fatal("no input produced no error")
	}
	if !strings.Contains(err.Error(), "no input") {
		t.Errorf("the error does not say what is missing: %v", err)
	}
}

func TestTheLongHelpShowsTheFontItself(t *testing.T) {
	if !strings.Contains(RootCmd.Long, calvin.AsciiFont("calvin")) {
		t.Error("the help does not render the font it describes")
	}
}
