package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/TwiN/go-pastebin"
	"github.com/urfave/cli/v2"
	"github.com/zerodoctor/zdcli/logger"
)

func PasteCmd() *cli.Command {
	return &cli.Command{
		Name:  "paste",
		Usage: "common commands to interact with pastebin.com. May need to login via this cli before use.",
		Subcommands: []*cli.Command{
			{
				Name:  "upload",
				Usage: "upload files in pastebin.com while keep the same pastebin key",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "folder",
					},
				},
				Action: func(ctx *cli.Context) error {
					var paths []string

					if ctx.String("folder") != "" {
						folder := ctx.String("folder")
						files, err := ioutil.ReadDir(folder)
						if err != nil {
							logger.Errorf("failed to read [folder=%s] [error=%s]", folder, err.Error())
							return nil
						}

						for _, file := range files {
							if file.IsDir() {
								continue
							}

							paths = append(paths, folder+"/"+file.Name())
						}
					}

					paths = append(paths, ctx.Args().Slice()...)
					PasteBinUpload(paths)

					return nil
				},
			},
		},
	}
}

func PasteBinUpload(paths []string) {
	fileMap := make(map[string]*os.File)
	for _, path := range paths {
		file, err := os.OpenFile(path, os.O_RDONLY, 0644)
		if err != nil {
			logger.Errorf("failed to read [file=%s] [error=%s]", path, err.Error())
			continue
		}

		name := file.Name()
		index := strings.LastIndex(name, ".")
		if index != -1 {
			name = name[:index]
			index = strings.LastIndex(name, "/")
			if index != -1 {
				name = name[index+1:]
			}
		}
		fileMap[name] = file
	}

	// TODO: integrate vault
	client, err := pastebin.NewClient(os.Getenv("PASTE_BIN_USER"), os.Getenv("PASTE_BIN_PASS"), os.Getenv("PASTE_BIN_KEY"))
	if err != nil {
		logger.Errorf("failed to create paste bin client [error=%s]", err.Error())
		return
	}

	content, err := client.GetAllUserPastes()
	if err != nil {
		logger.Errorf("failed to get all users pastes [error=%s]", err.Error())
		return
	}

	var minecraftCli strings.Builder
	for _, paste := range content {
		if file, ok := fileMap[paste.Title]; ok {
			client.DeletePaste(paste.Key)

			content, err := ioutil.ReadAll(file)
			if err != nil {
				logger.Errorf("failed to read [file=%s] [error=%s]", file.Name(), err.Error())
				continue
			}

			name := file.Name()
			ftype := ""
			index := strings.LastIndex(name, ".")
			if index != -1 {
				ftype = name[index+1:]
			}

			key, err := client.CreatePaste(
				pastebin.NewCreatePasteRequest(paste.Title, string(content), pastebin.ExpirationNever, pastebin.VisibilityPrivate, ftype),
			)
			if err != nil {
				logger.Errorf("failed to upload [file=%s] to pastebin [error=%s]", paste.Title, err.Error())
				continue
			}

			file.Close()
			logger.Infof("update paste [file=%s] [key=%s]", paste.Title, key)
			minecraftCli.WriteString("pastebin get " + key + " " + paste.Title + ".lua && ")
			delete(fileMap, paste.Title)
		}
	}

	for title, file := range fileMap {
		content, err := ioutil.ReadAll(file)
		if err != nil {
			logger.Errorf("failed to read [file=%s] [error=%s]", file.Name(), err.Error())
			continue
		}

		name := file.Name()
		ftype := ""
		index := strings.LastIndex(name, ".")
		if index != -1 {
			ftype = name[index+1:]
		}

		key, err := client.CreatePaste(
			pastebin.NewCreatePasteRequest(title, string(content), pastebin.ExpirationNever, pastebin.VisibilityPrivate, ftype),
		)
		if err != nil {
			logger.Errorf("failed to upload [file=%s] to pastebin [error=%s]", title, err.Error())
			continue
		}

		file.Close()
		logger.Infof("create paste [file=%s] [key=%s]", title, key)
		minecraftCli.WriteString("pastebin get " + key + " " + title + ".lua && ")
	}

	fmt.Println(minecraftCli.String()[:minecraftCli.Len()-len(" && ")])
}