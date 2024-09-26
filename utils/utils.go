package utils

import "path/filepath"

func GetRepoBasePath() string {
	return filepath.Join(".", ".xit")
}
