/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/fernandodona/xit/commit"
	"github.com/fernandodona/xit/hash"
	"github.com/fernandodona/xit/index"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the added, modified and untracked files",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("status called")

		dirsToIgnore := []string{".git", ".xit"}
		var untrackedFiles []string
		var modifiedFiles []string
		var stagedFiles []string
		filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			for _, dir := range dirsToIgnore {
				if path == dir {
					return filepath.SkipDir
				}
			}

			if info.IsDir() {
				return nil
			}

			idx, err := index.Build()
			if err != nil {
				return err
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			content, err := io.ReadAll(file)
			if err != nil {
				return err
			}

			currentVersionHash, err := hash.GetHashCode(content)
			if err != nil {
				return err
			}

			relativePath, _ := filepath.Rel(".", path)
			blobInfo, exists := idx.Objects[relativePath]
			// nunca adicionados
			if !exists {
				untrackedFiles = append(untrackedFiles, relativePath)
			}
			// adicionados porém modificados
			if exists && blobInfo.Hash != currentVersionHash {
				modifiedFiles = append(modifiedFiles, relativePath)
			}
			// adicionados na staging area
			// TODO: comparar com versão do commit atual
			head, err := commit.GetHeadCommit()
			if err != nil {
				head = &commit.Commit{Index: index.Index{Objects: make(map[string]index.IndexEntry)}}
			}
			if exists && blobInfo.Hash == currentVersionHash && blobInfo.Hash != head.Index.Objects[relativePath].Hash {
				stagedFiles = append(stagedFiles, relativePath)
			}

			return nil
		})

		fmt.Println("Staged Files:")
		for _, file := range stagedFiles {
			fmt.Println(color.GreenString("\t" + file))
		}

		fmt.Println("Modified Files:")
		for _, file := range modifiedFiles {
			fmt.Println(color.MagentaString("\t" + file))
		}

		fmt.Println("Untracked Files:")
		for _, file := range untrackedFiles {
			fmt.Println(color.RedString("\t" + file))
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
