// Command sklp-create scaffolds a new skalpai-style project from the
// embedded bootstrap template.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	root := &cobra.Command{
		Use:   "sklp-create",
		Short: "Scaffold a new skalpai-style project",
	}
	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the CLI version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	})
	root.AddCommand(newStartCmd())
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
