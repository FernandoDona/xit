/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/fernandodona/xit/index"
	"github.com/fernandodona/xit/object"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a file to staging area",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called")

		for _, arg := range args {
			// verifica se arquivo existe
			if _, err := os.Stat(arg); err != nil {
				log.Fatal(err)
			}
			// abre o arquivo
			f, err := os.Open(arg)
			if err != nil {
				log.Fatal(err)
			}

			blob, err := object.CreateBlob(f)
			if err != nil {
				log.Fatal(err)
			}

			if err := index.Add(f, blob); err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
