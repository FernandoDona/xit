package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "xit",
	Short: "Xit is my git clone study project.",
	Long: `Built with go that i didnt know until now. This is xit.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},	
}

rootCmd.AddCommand(addCmd)


func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
