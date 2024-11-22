# calvin
convert text to Calvin S ascii font (https://patorjk.com/software/taag/#p=display&amp;f=Calvin%20S&amp;t=)


example:

```
$ echo "test" | go run cmd/calvin/calvin.go
┌┬┐┌─┐┌─┐┌┬┐
 │ ├┤ └─┐ │
 ┴ └─┘└─┘ ┴

```

library usage example

```
package main

import (
	"fmt"
	"github.com/0magnet/calvin"
)

func main() {
		fmt.Println(calvin.AsciiFont("test"))
}

```
