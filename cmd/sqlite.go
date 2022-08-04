package cmd

import (
	"errors"

	"github.com/urfave/cli/v2"
	"github.com/zerodoctor/zdcli/config"
	"github.com/zerodoctor/zdcli/db"
)

type SqliteCmd struct{}

func NewSqlitecmd(cfg *config.Config) *cli.Command {
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

			&cli.StringFlag{
				Name:  "file_path",
				Usage: "name of env file",
			},

			&cli.StringFlag{
				Name:  "project_name",
				Usage: "name of project env file belongs to",
			},

			&cli.StringSliceFlag{
				Name:  "save",
				Usage: "save env file[s] and store into sqlite db [name=unix_timestamp.env.db]",
			},

			&cli.StringFlag{
				Name:  "read",
				Usage: "output content from env file[s]",
			},

			&cli.BoolFlag{
				Name:  "list",
				Usage: "outputs a list of env files",
			},

			&cli.BoolFlag{
				Name:  "write",
				Usage: "writes a env file from db",
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

			return nil
		},
	}
}
