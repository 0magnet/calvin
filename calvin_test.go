package calvin

import (
	"strings"
	"testing"
)

func TestAsciiFontSingleLine(t *testing.T) {
	got := AsciiFont("hi")
	want := strings.Join([]string{
		`┬ ┬┬`,
		`├─┤│`,
		`┴ ┴┴`,
	}, "\n")
	if got != want {
		t.Errorf("AsciiFont(\"hi\") =\n%s\nwant\n%s", got, want)
	}
}

// Each input line gets its own three-row block. Before this, newlines were not
// in the font map and so were dropped, collapsing every line into one.
func TestAsciiFontNewlines(t *testing.T) {
	got := AsciiFont("hi\nhi")
	if lines := strings.Split(got, "\n"); len(lines) != 6 {
		t.Fatalf("got %d lines, want 6:\n%s", len(lines), got)
	}
	half := AsciiFont("hi")
	if got != half+"\n"+half {
		t.Errorf("two lines should render as the same block twice, got\n%s", got)
	}
}

func TestAsciiFontLineEndings(t *testing.T) {
	want := AsciiFont("a\nb")
	for _, input := range []string{"a\r\nb", "a\rb"} {
		if got := AsciiFont(input); got != want {
			t.Errorf("AsciiFont(%q) did not match the \\n rendering:\n%s", input, got)
		}
	}
}

// Piped input almost always ends in a newline; it should not add an empty block.
func TestAsciiFontTrailingNewline(t *testing.T) {
	if got, want := AsciiFont("hi\n"), AsciiFont("hi"); got != want {
		t.Errorf("a trailing newline changed the output:\n%s", got)
	}
	// An explicit blank line, however, is real input and should render one.
	if lines := strings.Split(AsciiFont("hi\n\n"), "\n"); len(lines) != 6 {
		t.Errorf("a blank line should render an empty block, got %d lines", len(lines))
	}
}

func TestAsciiFontTabs(t *testing.T) {
	if got, want := AsciiFont("a\tb"), AsciiFont("a    b"); got != want {
		t.Errorf("tab did not expand to %d spaces:\n%s", tabWidth, got)
	}
}

// The font now covers every printable ASCII character, so anything typeable
// renders rather than silently vanishing. Upstream Calvin S stops at the
// letters and a handful of symbols; the rest are this port's extension.
func TestPrintableASCIIIsComplete(t *testing.T) {
	for ch := rune(' '); ch <= '~'; ch++ {
		if _, ok := boxFont[ch]; !ok {
			t.Errorf("no glyph for %q (0x%02x)", ch, ch)
		}
	}
	if AsciiFont("2026") == "" {
		t.Error("digits rendered as nothing")
	}
	// A whole line of symbols must survive intact, not collapse to blanks.
	if strings.TrimSpace(AsciiFont(`~/.cache | grep -c "x" > n & echo $?`)) == "" {
		t.Error("a line of symbols rendered as nothing")
	}
}

// Characters that are distinct when typed must stay distinct when rendered,
// or the output is ambiguous to read back.
func TestGlyphsAreDistinct(t *testing.T) {
	seen := make(map[string]rune, len(boxFont))
	for ch := rune(' '); ch <= '~'; ch++ {
		glyph, ok := boxFont[ch]
		if !ok {
			continue
		}
		key := strings.Join(glyph, "\n")
		if prev, dup := seen[key]; dup {
			t.Errorf("%q and %q render identically:\n%s", prev, ch, key)
			continue
		}
		seen[key] = ch
	}
}

// Every glyph must be rectangular, or rows drift out of alignment as they are
// concatenated.
func TestGlyphsAreRectangular(t *testing.T) {
	for ch, glyph := range boxFont {
		if len(glyph) != glyphRows {
			t.Errorf("glyph %q has %d rows, want %d", ch, len(glyph), glyphRows)
			continue
		}
		width := len([]rune(glyph[0]))
		for i, row := range glyph {
			if w := len([]rune(row)); w != width {
				t.Errorf("glyph %q row %d is %d wide, want %d", ch, i, w, width)
			}
		}
	}
}

// Characters the font does not define are skipped, as they are upstream.
func TestAsciiFontUnknownRunes(t *testing.T) {
	if got, want := AsciiFont("héi"), AsciiFont("hi"); got != want {
		t.Errorf("undefined rune was not skipped:\n%s", got)
	}
}

func TestBlackboardBold(t *testing.T) {
	if got, want := BlackboardBold("Go 42"), "𝔾𝕠 𝟜𝟚"; got != want {
		t.Errorf("BlackboardBold = %q, want %q", got, want)
	}
	// Unmapped runes, newlines included, pass through untouched.
	if got, want := BlackboardBold("a\n!"), "𝕒\n!"; got != want {
		t.Errorf("BlackboardBold = %q, want %q", got, want)
	}
}

// A space is two columns wide, matching the reference font's two hardblanks.
// Getting this wrong silently changes the spacing of every rendered string.
func TestSpaceGlyph(t *testing.T) {
	glyph, ok := boxFont[' ']
	if !ok {
		t.Fatal("no glyph for space")
	}
	for i, row := range glyph {
		if row != "  " {
			t.Errorf("space row %d = %q, want two spaces", i, row)
		}
	}
	// A space must actually separate its neighbours.
	if AsciiFont("a b") == AsciiFont("ab") {
		t.Error("space did not separate the surrounding glyphs")
	}
	if got, want := AsciiFont("a  b"), AsciiFont("a b"); got == want {
		t.Error("two spaces rendered the same as one")
	}
}
