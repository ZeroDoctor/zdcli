package cmd

import (
	"errors"

	"github.com/urfave/cli/v2"
	"github.com/zerodoctor/zdcli/config"
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
		Usage: "save/load/read a env file using sqlite",
		Action: func(ctx *cli.Context) error {
			return nil
		},
	}
}
