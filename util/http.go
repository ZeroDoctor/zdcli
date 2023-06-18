package util

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"

	"github.com/mholt/archiver/v4"
	"github.com/zerodoctor/zdcli/logger"
	zdutil "github.com/zerodoctor/zdgo-util"
)

func ExtractFromHttpResponse(ctx context.Context, file string, targetDir string, reader io.ReadCloser) error {
	format, input, err := archiver.Identify(file, reader)
	if err != nil {
		return err
	}
	fmt.Printf("[INFO] found compress file format [type=%s]\n", format.Name())

	if err := saveCompressFile(file, input); err != nil {
		return err
	}
	reader.Close()

	fmt.Printf("[INFO] extracting archive [file=%s]\n", file)
	fsys, err := archiver.FileSystem(ctx, file)
	if err != nil {
		return err
	}

	return extractAllFromFileSystem(fsys, targetDir)
}

func saveCompressFile(file string, reader io.Reader) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	return os.WriteFile(file, data, 0777)
}

func extractAllFromFileSystem(fsys fs.FS, targetDir string) error {
	if dir, ok := fsys.(fs.ReadDirFS); ok {
		stack := zdutil.NewStack(".")

		for stack.Len() > 0 {
			path := *stack.Pop()
			entries, err := dir.ReadDir(path)
			if err != nil {
				return err
			}

			for i := range entries {
				path := path + "/" + entries[i].Name()
				fmt.Printf("[INFO] extracting [file=%s]...\n", path)

				targetPath := targetDir + "/" + path
				if entries[i].IsDir() {
					os.MkdirAll(targetPath, 0755)
					stack.Push(path)
					continue
				}

				if path[:2] == "./" {
					path = path[2:]
				}

				file, err := fsys.Open(path)
				if err != nil {
					fmt.Printf("[ERROR] failed to decompress [file=%s] [error=%s]\n", path, err.Error())
					continue
				}
				defer file.Close()

				data, err := io.ReadAll(file)
				if err != nil {
					fmt.Printf("[ERROR] failed to read all [file=%s] [error=%s]\n", path, err.Error())
					continue
				}

				if err := os.WriteFile(targetPath, data, 0777); err != nil {
					fmt.Printf("[ERROR] failed to write [file=%s] [error=%s]\n", path, err.Error())
					continue
				}
			}
		}

		return nil
	}

	return errors.New("file system not found")
}

func FollowDownloadRedirection(URL string, resp *http.Response, handleResponse func(resp *http.Response) error) (*http.Response, error) {
	for URL != resp.Request.URL.String() {
		logger.Infof("following download redirection [url=%s]", resp.Request.URL.String())
		URL = resp.Request.URL.String()

		if err := handleResponse(resp); err != nil {
			return resp, err
		}

		req, err := http.NewRequest(http.MethodGet, resp.Request.URL.String(), nil)
		if err != nil {
			return resp, err
		}

		client := http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			return resp, err
		}
	}

	return resp, nil
}
