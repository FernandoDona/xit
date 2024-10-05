// init.go
/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fernandodona/xit/commit"
	"github.com/fernandodona/xit/index" // Import the index package
	"github.com/fernandodona/xit/object"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the repository",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")

		workingDirectory, _ := os.Getwd()

		fmt.Println(workingDirectory)
		if err := setupFiles(); err != nil {
			log.Fatal(err)
		}
	},
}

func setupFiles() error {
	// verifica se existem os arquivos para não sobreescrever
	if _, err := os.Stat(object.DirPath); err == nil {
		return err
	}
	if _, err := os.Stat(index.FilePath); err == nil {
		return err
	}
	if _, err := os.Stat(commit.HeadPath); err == nil {
		return err
	}

	//create objects folder
	err := os.MkdirAll(object.DirPath, 0644)
	if err != nil {
		return err
	}
	//create index file
	indexFile, err := os.OpenFile(index.FilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
		return nil
	}
	defer indexFile.Close()

	idx := index.Index{Objects: make(map[string]index.IndexEntry)}
	content, err := json.Marshal(idx)
	if err != nil {
		return err
	}
	io.WriteString(indexFile, string(content))

	_, err = os.Create(commit.HeadPath)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
