package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

	index := strings.LastIndex(path, "/")
	if index == -1 {
		return path, fmt.Errorf("exec path is messed up [path=%s]", path)
	}
	path = path[:index]

	return path, err
}