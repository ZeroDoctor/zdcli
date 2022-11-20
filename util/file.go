package util

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	zdgoutil "github.com/zerodoctor/zdgo-util"
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
		index = strings.LastIndex(path, "\\")
		if index == -1 {
			return path, fmt.Errorf("exec path is messed up [path=%s]", path)
		}
	}
	path = path[:index]

	return path, err
}

func GetFile(filename string) (os.FileInfo, error) {
	return os.Stat(filename)
}

type File struct {
	Path string
	fs.FileInfo
}

func NewFileArray(root string, files ...fs.FileInfo) []File {
	var f []File

	for _, file := range files {
		f = append(f, File{
			Path:     root,
			FileInfo: file,
		})
	}

	return f
}

func GetAllFiles(file string) ([]File, error) {
	var result []File
	if !zdgoutil.FolderExists(file) {
		var f File
		var err error

		f.FileInfo, err = GetFile(file)
		if err != nil {
			return result, err
		}

		index := strings.LastIndex(file, "/")
		if index < 0 {
			index = strings.LastIndex(file, "\\")
		}

		f.Path = "."
		if index > 0 {
			f.Path = file[:index]
		}

		result = append(result, f)

		return result, nil
	}

	dir, err := ioutil.ReadDir(file)
	if err != nil {
		return result, err
	}

	stack := NewStack(NewFileArray(file, dir...)...)
	for stack.Len() > 0 {
		f := *stack.Pop()
		if f.IsDir() {
			root := f.Path + "/" + f.Name()
			dir, err = ioutil.ReadDir(root)
			if err != nil {
				return result, err
			}

			stack.Push(NewFileArray(root, dir...)...)
			continue
		}

		fmt.Printf("[found=%s]\n", f.Name())
		result = append(result, f)
	}

	return result, nil
}
