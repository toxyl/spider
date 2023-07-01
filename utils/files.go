package utils

import (
	"fmt"
	"os"
)

func FileDelete(path string) error {
	return os.Remove(path)
}

func FileExists(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	_, err = file.Stat()
	return err == nil
}

func FileRead(path string) (string, error) {
	if !FileExists(path) {
		return "", fmt.Errorf("file %s does not exist", path)
	}
	bytes, err := os.ReadFile(path)
	return string(bytes), err
}

func FileWrite(path, content string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create %s file", path)
	}
	return os.WriteFile(f.Name(), []byte(content), 0644)
}
