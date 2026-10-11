// package commands cmd/calvin/commands/serve.go
package commands

import (
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// serveCmd is a static file server for the demo page, so the wasm can be
// loaded over http rather than file:// (which blocks instantiateStreaming).
func serveCmd() *cobra.Command {
	var addr, dir string
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "serve the demo page over http",
		Long:  "serve the demo page over http",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			log.Printf("serving %s on http://localhost%s", dir, addr)
			return http.ListenAndServe(addr, http.FileServer(http.Dir(dir))) //nolint:gosec
		},
	}
	cmd.Flags().StringVar(&addr, "addr", ":8797", "listen address")
	cmd.Flags().StringVar(&dir, "dir", "docs", "directory to serve")
	return cmd
}
