package commit

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/fernandodona/xit/hash"
	"github.com/fernandodona/xit/index"
	"github.com/fernandodona/xit/object"
	"github.com/fernandodona/xit/utils"
)

var HeadPath string = filepath.Join(utils.GetRepoBasePath(), "HEAD")

type Commit struct {
	Hash       string      `json:"hash"`
	ParentHash string      `json:"parent_hash"`
	Index      index.Index `json:"index"`
	Message    string      `json:"message"`
	Date       time.Time   `json:"date"`
	Parent     *Commit
}

func (commit *Commit) String() string {
	yellow := color.New(color.FgYellow).SprintFunc()
	// Usar um string builder pra concatenar
	output := fmt.Sprintf("Commit: %s\nMessage: %s\nDate:%s", yellow(commit.Hash), commit.Message, commit.Date)

	changedFiles, _ := GetCommitedFiles(commit)
	if len(changedFiles) > 0 {
		output += "\nChanges:"
	}
	for _, file := range changedFiles {
		output += fmt.Sprintf("\n\t%s", file)
	}

	return output
}

func Create(message string) error {
	var commit = &Commit{Message: message}
	idx, err := index.Build()
	if err != nil {
		return err
	}
	headHash, err := getHeadHash()
	if err != nil {
		return err
	}
	hashes := make([]string, len(idx.Objects))
	for _, entry := range idx.Objects {
		hashes = append(hashes, entry.Hash)
	}
	sort.Strings(hashes)
	stringToHash := strings.Join(hashes, ".")
	hash, err := hash.GetHashCode([]byte(stringToHash))
	if err != nil {
		return err
	}

	if hash == headHash {
		return errors.New("commit already created")
	}

	commit.Index = idx
	commit.Hash = hash
	commit.Date = time.Now()
	commit.ParentHash = headHash

	// Refatorar trecho
	// faz a mesma coisa que no package Object
	content, err := json.Marshal(commit)
	if err != nil {
		return err
	}

	folder := hash[:2]
	filename := hash[2:]
	path := filepath.Join(object.DirPath, folder)

	if err := os.Mkdir(path, 0644); err != nil && !os.IsExist(err) {
		return err
	}

	blob, err := os.OpenFile(filepath.Join(path, filename), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer blob.Close()

	var b bytes.Buffer
	compressor := gzip.NewWriter(&b)
	if _, err := compressor.Write(content); err != nil {
		return err
	}
	compressor.Close()

	if _, err := io.WriteString(blob, b.String()); err != nil {
		return err
	}
	// fim do trecho

	headFile, err := os.Create(HeadPath)
	if err != nil {
		return err
	}
	if _, err = io.WriteString(headFile, hash); err != nil {
		return err
	}

	defer headFile.Close()

	return nil
}

func Build(hash string) (*Commit, error) {
	commitFile, err := os.Open(filepath.Join(object.DirPath, hash[:2], hash[2:]))
	if err != nil {
		return nil, err
	}
	defer commitFile.Close()

	decompressor, err := gzip.NewReader(commitFile)
	if err != nil {
		return nil, err
	}
	defer decompressor.Close()

	content, err := io.ReadAll(decompressor)
	if err != nil {
		return nil, err
	}

	var commit Commit
	if err := json.Unmarshal(content, &commit); err != nil {
		return nil, err
	}

	return &commit, nil
}

func GetHeadCommit() (*Commit, error) {
	hash, err := getHeadHash()
	if err != nil {
		return nil, err
	}

	if hash == "" {
		return nil, errors.New("there is no HEAD")
	}

	commit, err := Build(hash)
	if err != nil {
		return nil, err
	}

	return commit, nil
}

func GetCommitedFiles(commit *Commit) ([]string, error) {
	commitedFiles := make([]string, 0)

	if commit.ParentHash == "" {
		for name := range commit.Index.Objects {
			commitedFiles = append(commitedFiles, name)
		}
		return commitedFiles, nil
	}

	parent, err := Build(commit.ParentHash)
	if err != nil {
		return nil, err
	}

	for name, entry := range commit.Index.Objects {
		if entry.Hash != parent.Index.Objects[name].Hash {
			commitedFiles = append(commitedFiles, name)
		}
	}

	return commitedFiles, nil
}

func getHeadHash() (string, error) {
	headFile, err := os.Open(HeadPath)
	if err != nil {
		return "", err
	}
	defer headFile.Close()

	hash, err := io.ReadAll(headFile)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
