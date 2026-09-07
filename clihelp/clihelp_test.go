package clihelp

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// A test binary is stamped by the toolchain, so BuildInfo is present here and
// the accessors answer. What they answer is not fixed — the version is
// "(devel)" under `go test` — so these assert shape, not content.
func TestBuildInfoIsAvailable(t *testing.T) {
	if BuildInfo() == nil {
		t.Fatal("BuildInfo() is nil; a test binary should carry one")
	}
	if GoVersion() == "" {
		t.Error("GoVersion() is empty")
	}
	if !strings.HasPrefix(GoVersion(), "go") {
		t.Errorf("GoVersion() = %q, want a go… string", GoVersion())
	}
}

// The banner is the ASCII name with the build underneath. The name must not
// appear as plain text — that would mean the font silently did nothing.
func TestBannerIsAsciiPlusBuild(t *testing.T) {
	b := Banner("dict")
	if strings.Contains(b, "\ndict\n") {
		t.Error("banner contains the plain name; the ASCII font did not run")
	}
	if !strings.Contains(b, "built with go") {
		t.Errorf("banner has no toolchain line:\n%s", b)
	}
	// Three rows of font plus at least one build line. A test binary carries no
	// VCS settings and reports "(devel)", so only the toolchain line is certain.
	if lines := strings.Count(b, "\n"); lines < 3 {
		t.Errorf("banner is %d newlines, want the font plus a build line:\n%s", lines, b)
	}
}

func TestBannerWithAppendsProse(t *testing.T) {
	b := BannerWith("dict", "look up a word\n")
	if !strings.HasSuffix(b, "look up a word") {
		t.Errorf("prose missing or not trimmed:\n%q", b)
	}
	if !strings.Contains(b, "built with go") {
		t.Error("BannerWith dropped the build block")
	}
}

// Init must leave the command runnable and answering the two flags, and must
// not disturb a Run the command already had.
func TestInitKeepsExistingRun(t *testing.T) {
	ran := false
	cmd := &cobra.Command{Use: "thing", Run: func(*cobra.Command, []string) { ran = true }}
	Init(cmd, "thing", false)

	cmd.SetArgs(nil)
	cmd.SetOut(&bytes.Buffer{})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !ran {
		t.Error("Init discarded the command's Run")
	}
}

func TestVersionFlagPrintsVersion(t *testing.T) {
	cmd := &cobra.Command{Use: "thing", Run: func(*cobra.Command, []string) {}}
	Init(cmd, "thing", false)

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"-b"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(out.String()); got != Version() {
		t.Errorf("-b printed %q, want %q", got, Version())
	}
}

func TestInfoFlagPrintsBuildInfo(t *testing.T) {
	cmd := &cobra.Command{Use: "thing", Run: func(*cobra.Command, []string) {}}
	Init(cmd, "thing", false)

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"-d"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "go\t") && !strings.Contains(out.String(), "mod\t") {
		t.Errorf("-d did not print a debug.BuildInfo:\n%s", out.String())
	}
}

// A command with no Run of its own should print help rather than fall through
// silently, which is what cobra does for a bare root with subcommands.
func TestRunlessCommandShowsHelp(t *testing.T) {
	cmd := &cobra.Command{Use: "thing"}
	// The subcommand needs a Run: cobra.IsAvailableCommand excludes one that is
	// neither runnable nor a parent, and the template filters on it.
	cmd.AddCommand(&cobra.Command{Use: "sub", Short: "a subcommand", Run: func(*cobra.Command, []string) {}})
	Init(cmd, "thing", false)

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs(nil)
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "sub") {
		t.Errorf("help did not list the subcommand:\n%s", out.String())
	}
}

// --help must not appear in the listing: it is there so the tree can be asked
// for help, not so it can describe itself.
func TestHelpFlagIsHidden(t *testing.T) {
	cmd := &cobra.Command{Use: "thing", Run: func(*cobra.Command, []string) {}}
	Init(cmd, "thing", false)
	f := cmd.PersistentFlags().Lookup("help")
	if f == nil {
		t.Fatal("no --help flag installed")
	}
	if !f.Hidden {
		t.Error("--help is not hidden")
	}
}
