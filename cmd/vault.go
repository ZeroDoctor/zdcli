package cmd

import (
	"errors"
	"fmt"
	"syscall"

	"github.com/urfave/cli/v2"
	"github.com/zerodoctor/zdcli/config"
	"github.com/zerodoctor/zdcli/logger"
	"github.com/zerodoctor/zdvault"
	"golang.org/x/term"
)

const MAIN_TOKEN string = "main_token"

func VaultCmd(cfg *config.Config) *cli.Command {
	if cfg.VaultEndpoint != "" {
		zdvault.SetEndpoint(cfg.VaultEndpoint)
	}

	return &cli.Command{
		Name:    "vault",
		Aliases: []string{"v"},
		Usage:   "commands that communicates with a vault server",
		Subcommands: []*cli.Command{
			VaultSetSubCmd(cfg),
			VaultLoginSubCmd(cfg),
			VaultRevokeSelfSubCmd(cfg),
		},
		Action: func(ctx *cli.Context) error {
			cli.ShowAppHelp(ctx)
			return errors.New("must provide additional subcommand(s)")
		},
	}
}

func VaultSetSubCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "set",
		Usage: "set various options needed for vault operations",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "endpoint",
				Aliases: []string{"e"},
			},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.String("endpoint") != "" {
				cfg.VaultEndpoint = ctx.String("endpoint")
				if err := cfg.Save(); err != nil {
					logger.Errorf("failed to save endpoint [error=%s]", err.Error())
				}
			}

			return nil
		},
	}
}

func VaultLoginSubCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "login",
		Usage: "used to login to vault server",
		Action: func(ctx *cli.Context) error {
			fmt.Print("enter username:")
			var username string
			if _, err := fmt.Scanln(&username); err != nil {
				logger.Errorf("failed to read user input [error=%s]", err.Error())
				return nil
			}

			fmt.Print("enter password:")
			bpass, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				logger.Errorf("failed to read password [error=%s]", err.Error())
			}
			password := string(bpass)
			fmt.Println()

			cred := zdvault.Cred{
				AppRole:  false,
				Username: username,
				Password: password,
				Key:      MAIN_TOKEN,
			}

			if err := zdvault.CreateNewToken(cred); err != nil {
				logger.Errorf("failed to create vault token [error=%s]", err.Error())
				return nil
			}

			if cfg.VaultTokens == nil {
				cfg.VaultTokens = make(map[string]string)
			}

			cfg.VaultTokens[cred.Key] = zdvault.GetToken(cred.Key)
			if err := cfg.Save(); err != nil {
				logger.Errorf("failed to save vault token and key [error=%s]", err.Error())
				return nil
			}

			return nil
		},
	}
}

func VaultRevokeSelfSubCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:    "revoke-self",
		Aliases: []string{"rs"},
		Usage:   "used to remove current session token",
		Action: func(ctx *cli.Context) error {
			if cfg.VaultTokens == nil {
				logger.Warn("vault tokens map is empty")
				return nil
			}

			if _, err := zdvault.RevokeSelfToken(MAIN_TOKEN); err != nil {
				logger.Errorf("failed to self revoken token [error=%s]", err.Error())
				return nil
			}

			delete(cfg.VaultTokens, MAIN_TOKEN)

			if err := cfg.Save(); err != nil {
				logger.Errorf("failed to deleted vault token [error=%s]", err.Error())
				return nil
			}

			return nil
		},
	}
}
