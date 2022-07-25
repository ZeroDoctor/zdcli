package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/linxGnu/goseaweedfs"
	"github.com/urfave/cli/v2"
	"github.com/zerodoctor/zdcli/config"
	"github.com/zerodoctor/zdcli/logger"
)

type SeaweedFS struct{}

func NewSeaweedFS(cfg *config.Config) *cli.Command {
	s := &SeaweedFS{}

	return &cli.Command{
		Name:        "weed",
		Aliases:     []string{"fs"},
		Description: "store folder and files to seaweed file system",
		Subcommands: []*cli.Command{
			s.UploadFilesCmd(cfg),
		},
		Action: func(ctx *cli.Context) error {
			cli.ShowAppHelp(ctx)
			return fmt.Errorf("additional sub command required")
		},
	}
}

func (s *SeaweedFS) ConnectFS(cfg *config.Config) (*goseaweedfs.Seaweed, error) {

	sw, error := goseaweedfs.NewSeaweed(
		cfg.SWFSMasterEndpoint,
		[]string{cfg.SWFSFilerEndpoint},
		4096,
		&http.Client{Timeout: 5 * time.Minute},
	)

	return sw, error
}

func (s *SeaweedFS) UploadFilesCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:        "upload",
		Aliases:     []string{"u"},
		Description: "upload file(s) to seaweed server",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "sources",
				Aliases:  []string{"src"},
				Required: true,
				Usage:    "a list of files to upload",
			},
			&cli.StringFlag{
				Name:    "destination",
				Aliases: []string{"dst"},
				Usage:   "what collection to put the source files",
			},
		},
		Action: func(ctx *cli.Context) error {
			sw, err := s.ConnectFS(cfg)
			if err != nil {
				logger.Errorf("failed to connect to file server  [error=%s]", err.Error())
				return nil
			}
			defer sw.Close()

			if err := s.UploadFiles(sw, ctx.StringSlice("sources"), ctx.String("destination")); err != nil {
				logger.Errorf("failed to upload [file(s)=%+v] [error=%s]", ctx.StringSlice("sources"), err.Error())
			}

			return nil
		},
	}
}

func (s *SeaweedFS) UploadFiles(sw *goseaweedfs.Seaweed, files []string, dest string) error {
	if len(files) <= 0 {
		return fmt.Errorf("source file(s) not found")
	}

	if len(files) == 1 {
		cm, fp, err := sw.UploadFile(files[0], dest, "")
		if err != nil {
			return err
		}

		logger.Infof("saved [file=%+v] [meta=%+v]", fp, cm)

		return nil
	}

	results, err := sw.BatchUploadFiles(files, dest, "")
	if err != nil {
		return err
	}

	for i := range results {
		logger.Infof("saved [file=%+v]", results[i])
	}

	return nil
}
