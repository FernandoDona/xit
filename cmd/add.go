package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var add = &cobra.Command{
	Use:   "add",
	Short: "Add file to stage",
	Long:  `This command add the file to stage to prepare for a commit.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Chamou o Add, ihaaa")
	},
}
