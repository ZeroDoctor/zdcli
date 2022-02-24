package util

import (
	"os"
	"path/filepath"
)

var EXEC_PATH string

func init() {
	var err error

	EXEC_PATH, err = GetExecPath()
	if err != nil {
		panic(err)
	}
}

func GetExecPath() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return path, err
	}

	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		return path, err
	}

	return path, err
}
