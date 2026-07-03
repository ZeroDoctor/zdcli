package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/zerodoctor/zdcli/logger"
)

type FileServerCmd struct{}

func NewFileServerCmd() *cli.Command {
	fs := &FileServerCmd{}

	return &cli.Command{
		Name:  "fs",
		Usage: "starts a simple file server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "dir",
				Value: ".",
				Usage: "directory to serve",
			},
			&cli.StringFlag{
				Name:  "port",
				Value: "8080",
				Usage: "port to serve on",
			},
		},
		Action: func(ctx *cli.Context) error {
			dir := ctx.String("dir")
			port := ctx.String("port")

			if err := fs.StartFileServer(dir, port); err != nil {
				return err
			}

			return nil
		},
	}
}

func (fs *FileServerCmd) StartFileServer(dir, port string) error {
	if dir == "." {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory [error=%s]", err.Error())
		}
		dir = cwd
	}

	fileServer := http.FileServer(http.Dir(dir))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Disable caching
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		// Log the request
		logger.Infof("serving [method=%s] [path=%s] [remote=%s]", r.Method, r.URL.Path, r.RemoteAddr)

		// Serve the file
		fileServer.ServeHTTP(w, r)
	})
	http.Handle("/", handler)

	addr := fmt.Sprintf(":%s", port)
	logger.Infof("starting file server [dir=%s] [port=%s]", dir, port)
	return http.ListenAndServe(addr, nil)
}
