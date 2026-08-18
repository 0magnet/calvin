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

func TestAsciiFontDigitsAndBrackets(t *testing.T) {
	for _, ch := range `0123456789[]():;'"/\` {
		glyph, ok := boxFont[ch]
		if !ok {
			t.Errorf("no glyph for %q", ch)
			continue
		}
		if len(glyph) != glyphRows {
			t.Errorf("glyph %q has %d rows, want %d", ch, len(glyph), glyphRows)
		}
	}
	if AsciiFont("2026") == "" {
		t.Error("digits rendered as nothing")
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
