package object

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fernandodona/xit/hash"
)

var DirPath string = filepath.Join(".", ".xit", "objects")

func CreateBlob(f *os.File) (*os.File, error) {
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	hash, err := hash.GetHashCode(content)
	if err != nil {
		return nil, err
	}

	folder := hash[:2]
	fileName := hash[2:]
	path := filepath.Join(DirPath, folder)

	if err := os.Mkdir(path, 0644); err != nil && !os.IsExist(err) {
		return nil, err
	}

	blob, err := os.OpenFile(filepath.Join(path, fileName), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	if err := compressObject(blob, content); err != nil {
		return nil, err
	}

	return blob, nil
}

func GetHashFromPath(filePath string) (string, error) {
	blobPathItems := strings.Split(filePath, string(os.PathSeparator))
	if len(blobPathItems) < 2 {
		return "", errors.New("este arquivo não é uma cópia versionada")
	}

	folder := blobPathItems[len(blobPathItems)-2]
	filename := blobPathItems[len(blobPathItems)-1]

	if len(folder) != 2 {
		return "", errors.New("este arquivo não é uma cópia versionada")
	}

	blobHash := folder + filename
	return blobHash, nil
}

func compressObject(object io.Writer, content []byte) error {
	var b bytes.Buffer
	compressor := gzip.NewWriter(&b)
	if _, err := compressor.Write(content); err != nil {
		return err
	}
	if err := compressor.Close(); err != nil {
		return err
	}
	if _, err := io.WriteString(object, b.String()); err != nil {
		return err
	}

	return nil
}
