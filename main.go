package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "xit",
		Short: "Xit is my git clone study project.",
		Long:  `Built with go that i didnt know until now. This is xit.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Starts a new repository",
		Long:  `This command starts a new repository in the current directory.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Chamou o init, uhulll")
		},
	}

	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add file to stage",
		Long:  `This command add the file to stage to prepare for a commit.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Chamou o Add, ihaaa")
		},
	}

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(initCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
