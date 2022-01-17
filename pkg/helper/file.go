package helper

import (
	"os"
)

// EnsureDir ensure directory exists
func EnsureDir(path string) (string, error) {
	if _, e := os.Stat(path); os.IsNotExist(e) {
		e = os.MkdirAll(path, os.ModePerm)
		if e != nil {
			return "", e
		}
	}
	return path, nil
}

// EnsureFile ensure directory exists
func EnsureFile(path string) (string, error) {
	if !FileExists(path) {
		f, e := os.Create(path)
		if e != nil {
			return "", e
		}
		f.Close()
	}
	return path, nil
}

// DirExists ensure if file exists
func DirExists(filePath string) bool {
	info, e := os.Stat(filePath)
	if os.IsNotExist(e) {
		return false
	}
	return info.IsDir()
}

// FileExists ensure if file exists
func FileExists(filePath string) bool {
	_, e := os.Stat(filePath)
	return os.IsExist(e)
}
