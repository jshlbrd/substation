package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/brexhq/substation/v2"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print version",
	Long: `'substation version' prints the version to stdout.
`,
	// Examples:
	//  substation version
	Example: `  substation version
`,
	Args: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, line := range strings.Split(substation.Version, "\n") {
			if strings.HasPrefix(line, "v") {
				fmt.Printf("%s\n", line)

				return nil
			}
		}

		return nil
	},
}
