package util

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	zdutil "github.com/zerodoctor/zdgo-util"
)

var BIN_PATH string
var EXEC_PATH string

func init() {
	var err error

	EXEC_PATH, err = zdutil.GetExecPath()
	if err != nil {
		panic(err) // TODO: avoid panics
	}

	BIN_PATH = EXEC_PATH + "/bin"
	os.Mkdir(BIN_PATH, 0755)
}

func GetFile(filename string) (os.FileInfo, error) {
	return os.Stat(filename)
}

type File struct {
	Path string
	fs.FileInfo
}

func NewFileArray(root string, entry ...fs.DirEntry) []File {
	var f []File

	for _, dir := range entry {
		file, err := dir.Info()
		if err != nil {
			continue
		}

		f = append(f, File{
			Path:     root,
			FileInfo: file,
		})
	}

	return f
}

func GetAllFiles(file string) ([]File, error) {
	var result []File
	if !zdutil.FolderExists(file) {
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

	dir, err := os.ReadDir(file)
	if err != nil {
		return result, err
	}

	stack := zdutil.NewStack(NewFileArray(file, dir...)...)
	for stack.Len() > 0 {
		f := *stack.Pop()
		if f.IsDir() {
			root := f.Path + "/" + f.Name()
			dir, err = os.ReadDir(root)
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
