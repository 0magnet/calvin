// Package calvin calvin.go
/*
convert text to ascii art
*/
package calvin

import (
	"strings"
)

// $ awk '{ if (NF == 0) next; if (/^[a-z]$/) { if (NR > 1) printf "},\n"; printf "\t\047%s\047: {", $0; getline; gsub(/\|$/, ""); printf "\042%s\042", $0; getline; gsub(/\|$/, ""); printf ", \042%s\042", $0; getline; gsub(/\|$/, ""); printf ", \042%s\042", $0 } } END { print "}" }' calvin.txt


var boxFont = map[rune][]string{
  'a': {"┌─┐", "├─┤", "┴ ┴"},
	'b': {"┌┐ ", "├┴┐", "└─┘"},
	'c': {"┌─┐", "│  ", "└─┘"},
	'd': {"┌┬┐", " ││", "─┴┘"},
	'e': {"┌─┐", "├┤ ", "└─┘"},
	'f': {"┌─┐", "├┤ ", "└  "},
	'g': {"┌─┐", "│ ┬", "└─┘"},
	'h': {"┬ ┬", "├─┤", "┴ ┴"},
	'i': {"┬", "│", "┴"},
	'j': {" ┬", " │", "└┘"},
	'k': {"┬┌─", "├┴┐", "┴ ┴"},
	'l': {"┬  ", "│  ", "┴─┘"},
	'm': {"┌┬┐", "│││", "┴ ┴"},
	'n': {"┌┐┌", "│││", "┘└┘"},
	'o': {"┌─┐", "│ │", "└─┘"},
	'p': {"┌─┐", "├─┘", "┴  "},
	'q': {"┌─┐ ", "│─┼┐", "└─┘└"},
	'r': {"┬─┐", "├┬┘", "┴└─"},
	's': {"┌─┐", "└─┐", "└─┘"},
	't': {"┌┬┐", " │ ", " ┴ "},
	'u': {"┬ ┬", "│ │", "└─┘"},
	'v': {"┬  ┬", "└┐┌┘", " └┘ "},
	'w': {"┬ ┬", "│││", "└┴┘"},
	'x': {"─┐ ┬", "┌┴┬┘", "┴ └─"},
	'y': {"┬ ┬", "└┬┘", " ┴ "},
	'z': {"┌─┐", "┌─┘", "└─┘"},
}

// ConvertToBoxFont converts a lowercase string to box drawing characters.
func AsciiFont(input string) string {
	var output [3]string

	for _, char := range input {
		if row, ok := boxFont[char]; ok {
			for i := 0; i < len(row); i++ {
				output[i] += row[i]
			}
		}
	}

	return strings.Join(output[:], "\n")
}
