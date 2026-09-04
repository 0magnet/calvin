//go:build js && wasm

// Command web is the calvin demo: type into the box, watch the letters turn.
//
// The whole package is two functions over a glyph table with no
// dependencies, so the demo is the same two calls a program would make —
// there is nothing else to show, and a page that showed more would be
// showing the page rather than the library.
package main

import (
	"strings"
	"syscall/js"

	"github.com/0magnet/calvin"
)

func main() {
	doc := js.Global().Get("document")
	in := doc.Call("getElementById", "in")
	box := doc.Call("getElementById", "box")
	bbb := doc.Call("getElementById", "bbb")

	render := func() {
		s := in.Get("value").String()
		if strings.TrimSpace(s) == "" {
			s = "calvin"
		}
		box.Set("textContent", calvin.AsciiFont(s))
		bbb.Set("textContent", calvin.BlackboardBold(s))
	}

	cb := js.FuncOf(func(js.Value, []js.Value) any { render(); return nil })
	in.Call("addEventListener", "input", cb)
	render()
	in.Call("focus")

	select {}
}
