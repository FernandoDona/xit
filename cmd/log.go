/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/fernandodona/xit/commit"
	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Shows the last commits",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("log called")

		currentCommit, err := commit.GetHeadCommit()
		if err != nil {
			fmt.Println("there are no commits")
			return
		}

		for {
			output := currentCommit.String()
			fmt.Println(output)

			if currentCommit.ParentHash == "" {
				break
			}
			fmt.Println("")
			currentCommit, _ = commit.Build(currentCommit.ParentHash)
		}
	},
}

func init() {
	rootCmd.AddCommand(logCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
