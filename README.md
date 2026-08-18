# calvin
convert text to Calvin S ascii font (https://patorjk.com/software/taag/#p=display&amp;f=Calvin%20S&amp;t=)


example:

```
$ echo "Hello, World!" | go run cmd/calvin/calvin.go
╦ ╦ ┌─┐┬  ┬  ┌─┐   ╦ ╦ ┌─┐┬─┐┬  ┌┬┐┬    
╠═╣ ├┤ │  │  │ │   ║║║ │ │├┬┘│   │││    
╩ ╩ └─┘┴─┘┴─┘└─┘┘  ╚╩╝ └─┘┴└─┴─┘─┴┘o    

```

library usage example

```
package main

import (
	"github.com/0magnet/calvin"
)

func main() {
	println(calvin.AsciiFont("Hello, World!"))
}

```

## Characters

The font covers `a-z`, `A-Z`, `[ ] ! @ # $ % ^ & * - _ , . ?` and space.

Digits `0-9` are an **extension**. Calvin S itself defines no glyphs for them,
so [TAAG](https://patorjk.com/software/taag/#p=display&f=Calvin%20S) renders
`2026` as nothing at all. These are drawn in the font's lowercase style — light
box-drawing, three rows — and shaped to stay legible against the letters they
most resemble: `0` is slashed so it does not read as `o`, `8` closes its lower
bowl where `a` has feet, and `2` keeps a flat base against `z`'s closed one.

```
┌─┐┌┐ ┌─┐┌─┐┬ ┬┌──┌─ ──┐┌─┐┌─┐
│┼│ │  ┌┘ ─┤└─┤└─┐├─┐ ┌┘├─┤└─┤
└─┘─┴─└──└─┘  ┴└─┘└─┘ ┴ └─┘ ─┘
```

Anything else is skipped, as it is upstream.

## Multi-line input

Each line renders as its own block:

```
$ printf 'Hello\nWorld' | calvin
╦ ╦ ┌─┐┬  ┬  ┌─┐
╠═╣ ├┤ │  │  │ │
╩ ╩ └─┘┴─┘┴─┘└─┘
╦ ╦ ┌─┐┬─┐┬  ┌┬┐
║║║ │ │├┬┘│   ││
╚╩╝ └─┘┴└─┴─┘─┴┘
```

`\r\n` and `\r` are accepted, tabs expand to spaces, and a single trailing
newline is ignored so piped input does not render an empty trailing block.
