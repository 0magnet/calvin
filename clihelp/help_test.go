package clihelp

import (
	"bytes"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// tree is a command tree shaped like the real one: prose, an example, flags,
// subcommands, and a global flag inherited from the root.
func tree(t *testing.T) *cobra.Command {
	t.Helper()
	root := &cobra.Command{
		Use:           "prog <command>",
		Short:         "does things",
		Long:          "prog — does things\n\nAt some length.",
		SilenceErrors: true,
		SilenceUsage:  true,
		Run:           func(c *cobra.Command, _ []string) { _ = c.Help() }, //nolint:errcheck,gosec
	}
	root.PersistentFlags().Bool("verbose", false, "say more")

	sub := &cobra.Command{
		Use:     "walk",
		Aliases: []string{"step"},
		Short:   "walks",
		Long:    "Walks somewhere.",
		Example: "  prog walk --far",
		Run:     func(*cobra.Command, []string) {},
	}
	sub.Flags().Bool("far", false, "go further")
	sub.Flags().IntP("count", "n", 3, "how many")
	root.AddCommand(sub, &cobra.Command{Use: "sit", Short: "sits", Run: func(*cobra.Command, []string) {}})

	// The tests compare plain text; a js host that exports TERM would otherwise
	// switch color on.
	color.NoColor = true

	// Init sets the banner, the coloring and the flags together.
	Init(root, "prog")
	return root
}

func help(t *testing.T, args ...string) string {
	t.Helper()
	root := tree(t)
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs(args)
	if err := root.Execute(); err != nil {
		t.Fatalf("%v: %v", args, err)
	}
	return out.String()
}

// TestHelpHasEverySection is the whole point: cobra's template cannot run under
// TinyGo, so this output is built by hand and something has to check that
// building it by hand did not quietly drop a section.
// TestHelpHasEverySection checks that the house style drops nothing cobra's
// own help would have shown.
func TestHelpHasEverySection(t *testing.T) {
	got := help(t, "walk", "--help")
	for _, want := range []string{
		"Walks somewhere.",
		"Usage:\n  prog walk [flags]\n",
		"Aliases:\n  walk, step\n",
		"Examples:\n  prog walk --far\n",
		"Flags:\n",
		"  -n, --count int   how many (default 3)\n",
		"      --far         go further\n",
		"Global Flags:\n",
		"      --verbose   say more\n",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
	if strings.Contains(got, "\n\n\n") {
		t.Errorf("a doubled blank line:\n%q", got)
	}
}

func TestRootHelpListsSubcommands(t *testing.T) {
	got := help(t, "--help")
	if !strings.Contains(got, "prog — does things\n\nAt some length.\n") {
		t.Errorf("prose missing:\n%s", got)
	}
	for _, sub := range [][2]string{{"sit", "sits"}, {"walk", "walks"}} {
		found := false
		for _, line := range strings.Split(got, "\n") {
			f := strings.Fields(line)
			if len(f) == 2 && f[0] == sub[0] && f[1] == sub[1] {
				found = true
			}
		}
		if !found {
			t.Errorf("no listing line for %s:\n%s", sub[0], got)
		}
	}
	// --help is hidden, so it does not sit in the listing describing itself.
	if strings.Contains(got, "-h, --help") {
		t.Errorf("the hidden --help flag was listed:\n%s", got)
	}
}
