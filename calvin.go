// Package calvin calvin.go
/*
convert text to ascii art
*/
package calvin

import (
	"strings"
)

// $ awk '{if (NF == 0) next; if (/^.$/) {if (NR > 1) printf "},\n"; printf "\t\047%s\047: {", $0; getline; gsub(/\|$/, ""); printf "\042%s\042", $0; getline; gsub(/\|$/, ""); printf ", \042%s\042", $0; getline; gsub(/\|$/, ""); printf ", \042%s\042", $0}} END {print "}"}' calvin.txt

var boxFont = map[rune][]string{
	'a': {`┌─┐`, `├─┤`, `┴ ┴`},
	'b': {`┌┐ `, `├┴┐`, `└─┘`},
	'c': {`┌─┐`, `│  `, `└─┘`},
	'd': {`┌┬┐`, ` ││`, `─┴┘`},
	'e': {`┌─┐`, `├┤ `, `└─┘`},
	'f': {`┌─┐`, `├┤ `, `└  `},
	'g': {`┌─┐`, `│ ┬`, `└─┘`},
	'h': {`┬ ┬`, `├─┤`, `┴ ┴`},
	'i': {`┬`, `│`, `┴`},
	'j': {` ┬`, ` │`, `└┘`},
	'k': {`┬┌─`, `├┴┐`, `┴ ┴`},
	'l': {`┬  `, `│  `, `┴─┘`},
	'm': {`┌┬┐`, `│││`, `┴ ┴`},
	'n': {`┌┐┌`, `│││`, `┘└┘`},
	'o': {`┌─┐`, `│ │`, `└─┘`},
	'p': {`┌─┐`, `├─┘`, `┴  `},
	'q': {`┌─┐ `, `│─┼┐`, `└─┘└`},
	'r': {`┬─┐`, `├┬┘`, `┴└─`},
	's': {`┌─┐`, `└─┐`, `└─┘`},
	't': {`┌┬┐`, ` │ `, ` ┴ `},
	'u': {`┬ ┬`, `│ │`, `└─┘`},
	'v': {`┬  ┬`, `└┐┌┘`, ` └┘ `},
	'w': {`┬ ┬`, `│││`, `└┴┘`},
	'x': {`─┐ ┬`, `┌┴┬┘`, `┴ └─`},
	'y': {`┬ ┬`, `└┬┘`, ` ┴ `},
	'z': {`┌─┐`, `┌─┘`, `└─┘`},
	'A': {`╔═╗ `, `╠═╣ `, `╩ ╩ `},
	'B': {`╔╗  `, `╠╩╗ `, `╚═╝ `},
	'C': {`╔═╗ `, `║   `, `╚═╝ `},
	'D': {`╔╦╗ `, ` ║║ `, `═╩╝ `},
	'E': {`╔═╗ `, `║╣  `, `╚═╝ `},
	'F': {`╔═╗ `, `╠╣  `, `╚   `},
	'G': {`╔═╗ `, `║ ╦ `, `╚═╝ `},
	'H': {`╦ ╦ `, `╠═╣ `, `╩ ╩ `},
	'I': {`╦   `, `║   `, `╩   `},
	'J': {` ╦  `, ` ║  `, `╚╝  `},
	'K': {`╦╔═ `, `╠╩╗ `, `╩ ╩ `},
	'L': {`╦   `, `║   `, `╩═╝ `},
	'M': {`╔╦╗ `, `║║║ `, `╩ ╩ `},
	'N': {`╔╗╔ `, `║║║ `, `╝╚╝ `},
	'O': {`╔═╗ `, `║ ║ `, `╚═╝ `},
	'P': {`╔═╗ `, `╠═╝ `, `╩   `},
	'Q': {`╔═╗ `, `║═╬╗`, `╚═╝╚`},
	'R': {`╦═╗ `, `╠╦╝ `, `╩╚═ `},
	'S': {`╔═╗ `, `╚═╗ `, `╚═╝ `},
	'T': {`╔╦╗ `, ` ║  `, ` ╩  `},
	'U': {`╦ ╦ `, `║ ║ `, `╚═╝ `},
	'V': {`╦  ╦`, `╚╗╔╝`, ` ╚╝ `},
	'W': {`╦ ╦ `, `║║║ `, `╚╩╝ `},
	'X': {`═╗ ╦`, `╔╩╦╝`, `╩ ╚═`},
	'Y': {`╦ ╦ `, `╚╦╝ `, ` ╩  `},
	'Z': {`╔═╗ `, `╔═╝ `, `╚═╝ `},
	'!': {`┬    `, `│    `, `o    `},
	'@': {`┌─┐  `, `│└┘  `, `└──  `},
	'#': {`─┼─┼─`, `─┼─┼─`, `     `},
	'$': {`┌┼┐  `, `└┼┐  `, `└┼┘  `},
	'%': {`O┬   `, `┌┘   `, `┴O   `},
	'^': {`/\   `, `     `, `     `},
	'&': {` ┬   `, `┌┼─  `, `└┘   `},
	'*': {`\│/  `, `─ ─  `, `/│\  `},
	'-': {`   `, `───`, `   `},
	'_': {`    `, `    `, `────`},
	',': {` `, ` `, `┘`},
	'.': {` `, ` `, `o`},
	'?': {`┌─┐`, ` ┌┘`, ` o `},
	' ': {`  `, `  `, `  `},
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

var charMap = map[rune]string{
	'A': "𝔸", 'B': "𝔹", 'C': "ℂ", 'D': "𝔻", 'E': "𝔼", 'F': "𝔽",
	'G': "𝔾", 'H': "ℍ", 'I': "𝕀", 'J': "𝕁", 'K': "𝕂", 'L': "𝕃",
	'M': "𝕄", 'N': "ℕ", 'O': "𝕆", 'P': "ℙ", 'Q': "ℚ", 'R': "ℝ",
	'S': "𝕊", 'T': "𝕋", 'U': "𝕌", 'V': "𝕍", 'W': "𝕎", 'X': "𝕏",
	'Y': "𝕐", 'Z': "ℤ",
	'a': "𝕒", 'b': "𝕓", 'c': "𝕔", 'd': "𝕕", 'e': "𝕖", 'f': "𝕗",
	'g': "𝕘", 'h': "𝕙", 'i': "𝕚", 'j': "𝕛", 'k': "𝕜", 'l': "𝕝",
	'm': "𝕞", 'n': "𝕟", 'o': "𝕠", 'p': "𝕡", 'q': "𝕢", 'r': "𝕣",
	's': "𝕤", 't': "𝕥", 'u': "𝕦", 'v': "𝕧", 'w': "𝕨", 'x': "𝕩",
	'y': "𝕪", 'z': "𝕫",
	'0': "𝟘", '1': "𝟙", '2': "𝟚", '3': "𝟛", '4': "𝟜",
	'5': "𝟝", '6': "𝟞", '7': "𝟟", '8': "𝟠", '9': "𝟡",
}

// BlackboardBold converts a string
func BlackboardBold(input string) string {
	var result strings.Builder
	for _, ch := range input {
		if specialChar, exists := charMap[ch]; exists {
			result.WriteString(specialChar)
		} else {
			result.WriteRune(ch)
		}
	}
	return result.String()
}
