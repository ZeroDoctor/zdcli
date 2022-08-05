package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/zerodoctor/zdcli/config"
	"github.com/zerodoctor/zdcli/db"
	"github.com/zerodoctor/zdcli/logger"
)

type SqliteCmd struct{}

func NewSqliteCmd(cfg *config.Config) *cli.Command {
	sql := &SqliteCmd{}

	return &cli.Command{
		Name:  "lite",
		Usage: "interacts with a sqlite database",
		Subcommands: []*cli.Command{
			sql.EnvSubCmd(cfg),
		},

		Action: func(ctx *cli.Context) error {
			cli.ShowAppHelp(ctx)
			return errors.New("must provide additional subcommand(s)")
		},
	}
}

func (s *SqliteCmd) EnvSubCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "env",
		Usage: "save/read/list/write a env file using sqlite",
		Flags: []cli.Flag{
			// &cli.BoolFlag{
			// 	Name:    "recursive",
			// 	Aliases: []string{"r"},
			// 	Usage:   "includes sub folders",
			// },

			// &cli.StringFlag{
			// 	Name:  "file_path",
			// 	Usage: "name of env file",
			// },

			&cli.StringFlag{
				Name:  "project_name",
				Usage: "name of project env file belongs to",
			},

			&cli.BoolFlag{
				Name:  "list",
				Usage: "outputs a list of env files",
			},

			&cli.StringSliceFlag{
				Name:  "save",
				Usage: "save env file[s] and store into sqlite db [name=unix_timestamp.env.db]",
			},

			&cli.StringSliceFlag{
				Name:  "read",
				Usage: "output content from env file[s]",
			},

			&cli.StringSliceFlag{
				Name:  "write",
				Usage: "writes a env file[s] from db",
			},
		},
		Action: func(ctx *cli.Context) error {

			dbh, err := db.NewHandler()
			if err != nil {
				return err
			}

			project := ctx.String("project_name")

			files := ctx.StringSlice("save")
			for i := range files {
				err := dbh.SaveEnvFile(files[i], project)
				if err != nil {
					return err
				}
			}

			if len(files) > 0 {
				logger.Info("finished saving env files")
			}

			var envs []db.Env
			reads := ctx.StringSlice("read")
			for i := range reads {
				envs, err = dbh.ReadEnvFile(project, reads[i])
				if err != nil {
					return err
				}

				fmt.Printf("[project=%s] [file=%s] [created_at=%s]\n[content=%s]\n\n",
					envs[i].ProjectName, envs[i].FileName,
					envs[i].CreatedAt.Format(time.RFC3339),
					string(envs[i].FileContent),
				)
			}

			writes := ctx.StringSlice("write")
			for i := range writes {
				envs, err = dbh.ReadEnvFile(project, writes[i])
				if err != nil {
					return err
				}

				if err = ioutil.WriteFile(envs[i].FileName, []byte(envs[i].FileContent), 0644); err != nil {
					return err
				}
			}

			if !ctx.Bool("list") {
				return nil
			}

			if project != "" {
				envs, err = dbh.ReadEnvProjectFiles(project)
			} else {
				envs, err = dbh.ReadAllEnv()
			}
			if err != nil {
				return err
			}

			for i := range envs {
				fmt.Printf("[project=%s] [file=%s] [created_at=%+v]\n",
					envs[i].ProjectName, envs[i].FileName,
					envs[i].CreatedAt,
				)
			}

			return nil
		},
	}
}
