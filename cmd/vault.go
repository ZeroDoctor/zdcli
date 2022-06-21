package cmd

import (
	"errors"
	"fmt"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
	"github.com/zerodoctor/zdcli/config"
	"github.com/zerodoctor/zdcli/logger"
	"github.com/zerodoctor/zdcli/tui/ui"
	"github.com/zerodoctor/zdvault"
	"golang.org/x/term"
)

const MAIN_TOKEN string = "main_token"

type VFlag int

func (b VFlag) Has(f VFlag) bool { return b&f != 0 }

const (
	VEndpoint VFlag = 1 << iota
	VToken
)

var ErrMissingVaultEndpoint error = errors.New("missing vault endpoint")
var ErrMissingVaultMainToken error = errors.New("missing vault main login token")

func validate(flag VFlag, cfg *config.Config) error {
	var errs []any

	if flag.Has(VEndpoint) && cfg.VaultEndpoint == "" {
		return ErrMissingVaultEndpoint
	}

	if flag.Has(VToken) {
		if _, ok := cfg.VaultTokens[MAIN_TOKEN]; !ok {
			errs = append(errs, ErrMissingVaultMainToken)
		}
	}

	var format string
	for range errs {
		format += "[error=%w] "
	}

	var err error
	if len(errs) > 0 {
		err = fmt.Errorf(format, errs...)
	}

	return err
}

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
			VaultNewSubCmd(cfg),
			VaultListSubCmd(cfg),
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
			if err := validate(VEndpoint, cfg); err != nil {
				return err
			}

			if _, ok := cfg.VaultTokens[MAIN_TOKEN]; ok {
				if _, err := zdvault.RevokeSelfToken(MAIN_TOKEN); err != nil {
					logger.Errorf("failed to revoke current token [error=%s]", err.Error())
				}
			}

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
			if err := validate(VEndpoint|VToken, cfg); err != nil {
				return err
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

func VaultNewSubCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:    "new",
		Aliases: []string{"n"},
		Usage:   "used to create auth methods or secrets",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "key",
				Aliases: []string{"k"},
				Usage:   "create new secret in key/value engine",
			},
		},
		Action: func(ctx *cli.Context) error {
			if err := validate(VEndpoint|VToken, cfg); err != nil {
				return err
			}

			if ctx.Bool("key") {
				return VaultNewKeySubCmd(cfg)
			}

			cli.ShowAppHelp(ctx)
			return errors.New("must provide additional flag(s)")
		},
	}
}

func VaultNewKeySubCmd(cfg *config.Config) error {
	return nil
}

func VaultListSubCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "list current folders",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "key",
				Aliases: []string{"k"},
				Usage:   "list folders in key/value engine",
			},
		},
		Action: func(ctx *cli.Context) error {
			if err := validate(VEndpoint|VToken, cfg); err != nil {
				return err
			}

			if ctx.Bool("key") {
				return VaultListKeySubCmd(cfg)
			}

			cli.ShowAppHelp(ctx)
			return errors.New("must provide additional flag(s)")
		},
	}
}

func VaultListKeySubCmd(cfg *config.Config) error {

	// TODO: fix text input to accept multiple tea.models
	text := ui.NewTextInput()
	text.Input.Prompt = "Enter path: "
	text.Input.Placeholder = "kv/secret"

	p := tea.NewProgram(text)

	if err := p.Start(); err != nil {
		logger.Errorf("failed to start tea ui [error=%s]", err.Error())
		return nil
	}

	if text.Input.Err != nil {
		logger.Errorf("failed to get input [error=%s]", text.Input.Err.Error())
		return nil
	}

	data, err := zdvault.ListSecret(MAIN_TOKEN, text.Input.Value())
	if err != nil {
		logger.Errorf("failed to list secret folders [error=%s]", err.Error())
	}

	fmt.Println(string(data))

	return nil
}
