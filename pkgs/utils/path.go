package utils

import (
	"path/filepath"
	"runtime"
	"strings"
)

func init() {
	initRelative()
}

var prefixPath string

func initRelative() {
	_, fileName, _, _ := runtime.Caller(0)
	prefixPath = fileName[:len(fileName)-len("utils/relative.go")]
}

// Relative returns the relative path to a file from the caller's file.
func relative(path string) string {
	return strings.TrimPrefix(filepath.ToSlash(path), prefixPath)
}

// AbsolutePath returns the absolute path to a file relative to the caller's file.
func AbsolutePath(rel string) string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Join(filepath.Dir(filename), rel)
}
