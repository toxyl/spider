package main

import (
	"os"
	"path/filepath"
)

func getTempFilePath(name string) string {
	return filepath.Join(os.TempDir(), name+".tmp")
}
