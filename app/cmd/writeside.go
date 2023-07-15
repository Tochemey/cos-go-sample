package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// writesideCmd represents the runWriteside command
var writesideCmd = &cobra.Command{
	Use:   "writeside",
	Short: "Run the commands and events handler service",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("runWriteside called")
	},
}

func init() { rootCmd.AddCommand(writesideCmd) }
