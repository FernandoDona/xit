package index

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/fernandodona/xit/object"
)

var FilePath string = filepath.Join(".", ".xit", "index")

type Index struct {
	Objects map[string]string `json:"objects"`
}

func Build() (Index, error) {
	var index Index
	filePath := filepath.Join(".", ".xit", "index")
	if _, err := os.Stat(filePath); err != nil {
		return index, err
	}
	file, err := os.Open(filePath)
	if err != nil {
		return index, err
	}

	reader := bufio.NewReader(file)
	content, err := io.ReadAll(reader)
	if err != nil {
		return index, err
	}

	if err := json.Unmarshal(content, &index); err != nil {
		return index, err
	}

	return index, nil
}

func Add(file, blob *os.File) error {
	idx, err := Build()
	if err != nil {
		return err
	}

	blobHash, err := object.GetHashFromPath(blob)
	if err != nil {
		return err
	}

	relativePath, _ := filepath.Rel(".", file.Name())

	idx.Objects[relativePath] = blobHash
	if err := Update(idx); err != nil {
		return err
	}

	return nil
}

func Update(index Index) error {
	content, err := json.Marshal(index)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(FilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil && !os.IsExist(err) {
		return err
	}
	io.WriteString(file, string(content))
	return nil
}
