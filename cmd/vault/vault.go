package vault

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/urfave/cli/v2"
	"github.com/zerodoctor/zdcli/config"
)

type VaultCmd struct {
	cfg    *config.Config
	client *vault.Client
	ctx    context.Context
}

func NewVaultCmd(cfg *config.Config) *cli.Command {
	if err := validate(VEndpoint, cfg); err != nil {
		cli.Exit(err.Error(), 1)
		return nil
	}

	client, err := vault.New(
		vault.WithAddress(cfg.VaultEndpoint),
		vault.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		msg := fmt.Sprintf("failed to create vault client [error=%s]", err.Error())
		cli.Exit(msg, 1)
		return nil
	}

	vault := &VaultCmd{
		cfg:    cfg,
		client: client,
	}

	return &cli.Command{
		Name:    "vault",
		Aliases: []string{"v"},
		Usage:   "commands that communicates with a vault server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "switch",
				Aliases: []string{"s"},
				Usage:   "switch current vault user",
			},
		},
		Subcommands: []*cli.Command{
			vault.LoginSubCmd(),
			vault.RevokeSelfSubCmd(),
			vault.NewSubCmd(),
			vault.GetSubCmd(),
			vault.ListSubCmd(),
			vault.EnableSubCmd(),
			vault.DisableSubCmd(),
		},
		Action: func(ctx *cli.Context) error {
			vault.ctx = ctx.Context

			if ctx.String("switch") != "" {
				if _, ok := vault.cfg.VaultTokens[ctx.String("switch")]; !ok {
					return fmt.Errorf(
						"token for [user=%s] does not exist. please login as user",
						ctx.String("switch"),
					)
				}

				vault.cfg.VaultUser = ctx.String("switch")
			}

			cli.ShowAppHelp(ctx)
			return errors.New("must provide additional subcommand(s)")
		},
	}
}

func (v *VaultCmd) GetToken() string      { return v.cfg.VaultTokens[v.cfg.VaultUser] }
func (v *VaultCmd) SetToken(token string) { v.cfg.VaultTokens[v.cfg.VaultUser] = token }

func (v *VaultCmd) LoginSubCmd() *cli.Command {
	return &cli.Command{
		Name:  "login",
		Usage: "used to login to vault server",
		Action: func(ctx *cli.Context) error {
			if err := validate(VEndpoint, v.cfg); err != nil {
				return err
			}

			v.ctx = ctx.Context
			return v.LoginUser()
		},
	}
}

func (v *VaultCmd) RevokeSelfSubCmd() *cli.Command {
	return &cli.Command{
		Name:    "revoke-self",
		Aliases: []string{"rs"},
		Usage:   "used to remove current session token",
		Action: func(ctx *cli.Context) error {
			if _, ok := v.cfg.VaultTokens[v.cfg.VaultUser]; !ok {
				return nil
			}

			if err := validate(VEndpoint|VToken, v.cfg); err != nil {
				return err
			}

			v.ctx = ctx.Context
			return v.RevokeSelf()
		},
	}
}

func (v *VaultCmd) NewSubCmd() *cli.Command {
	return &cli.Command{
		Name:  "new",
		Usage: "used to create auth methods or secrets",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "key",
				Aliases: []string{"k"},
				Usage:   "create new/update secret in key/value engine",
			},
			&cli.BoolFlag{
				Name:    "policy",
				Aliases: []string{"p"},
				Usage:   "create new/update policy in sys/policies",
			},
			&cli.StringFlag{
				Name:    "user",
				Aliases: []string{"u"},
				Usage:   "create a new user",
			},
		},
		Subcommands: []*cli.Command{
			&cli.Command{
				Name:    "alias",
				Aliases: []string{"a"},
				Usage:   "creates a new entity alias",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "with-meta",
						Aliases: []string{"m"},
						Usage:   "an alias with meta data",
					},
				},
				Action: func(ctx *cli.Context) error {
					if err := validate(VEndpoint|VToken, v.cfg); err != nil {
						return err
					}

					v.ctx = ctx.Context

					userName := ""
					if ctx.Args().Len() > 0 {
						userName = ctx.Args().Get(0)
					}

					_, err := v.NewAlias(userName, ctx.Bool("with-meta"))
					return err
				},
			},
			&cli.Command{
				Name:    "approle",
				Aliases: []string{"ar"},
				Usage:   "creates a new approle",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "name",
						Aliases:  []string{"n"},
						Usage:    "name of approle",
						Required: true,
					},
					&cli.BoolFlag{
						Name:    "token",
						Aliases: []string{"t"},
						Usage:   "use custom token settings",
					},
					&cli.BoolFlag{
						Name:    "secret",
						Aliases: []string{"s"},
						Usage:   "use custom secret settings",
					},
					&cli.BoolFlag{
						Name:    "create-secret",
						Aliases: []string{"cs"},
						Usage:   "use create new sceret",
					},
				},
				Action: func(ctx *cli.Context) error {
					if err := validate(VEndpoint|VToken, v.cfg); err != nil {
						return err
					}

					v.ctx = ctx.Context

					return v.NewApprole(
						ctx.String("name"), ctx.Bool("token"),
						ctx.Bool("secret"), ctx.Bool("create-secret"),
					)
				},
			},
		},
		Action: func(ctx *cli.Context) error {
			if err := validate(VEndpoint|VToken, v.cfg); err != nil {
				return err
			}

			v.ctx = ctx.Context

			if ctx.Bool("policy") {
				return v.NewPolicy()
			}

			if ctx.String("user") != "" {
				return v.NewUser(ctx.String("user"))
			}

			if ctx.Bool("key") {
				return v.NewKey()
			}

			cli.ShowAppHelp(ctx)
			return errors.New("must provide additional flag(s)")
		},
	}
}

func (v *VaultCmd) GetSubCmd() *cli.Command {
	return &cli.Command{
		Name:  "get",
		Usage: "get secrets from vault",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "key",
				Aliases: []string{"k"},
				Usage:   "fetch secret from key/value engine",
			},
			&cli.BoolFlag{
				Name:    "policy",
				Aliases: []string{"p"},
				Usage:   "read a policy",
			},
			&cli.StringFlag{
				Name:    "user",
				Aliases: []string{"u"},
				Usage:   "read user info",
			},
		},
		Action: func(ctx *cli.Context) error {
			if err := validate(VEndpoint|VToken, v.cfg); err != nil {
				return err
			}

			v.ctx = ctx.Context

			if ctx.Bool("key") {
				return v.GetKey()
			}

			if ctx.Bool("policy") {
				policy := ctx.Args().First()
				return v.GetPolicy(policy)
			}

			if ctx.String("user") != "" {
				return v.GetUser(ctx.String("user"))
			}

			cli.ShowAppHelp(ctx)
			return errors.New("must provide additional flag(s)")
		},
	}
}

func (v *VaultCmd) ListSubCmd() *cli.Command {
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
			&cli.BoolFlag{
				Name:    "policy",
				Aliases: []string{"p"},
				Usage:   "list all policies",
			},
			&cli.BoolFlag{
				Name:    "user",
				Aliases: []string{"u"},
				Usage:   "list all users",
			},
			&cli.BoolFlag{
				Name:    "approle",
				Aliases: []string{"ar"},
				Usage:   "list all approles",
			},
			&cli.BoolFlag{
				Name:    "mount",
				Aliases: []string{"m"},
				Usage:   "list all secret mounts",
			},
			&cli.StringFlag{
				Name:    "approle-secrets",
				Aliases: []string{"ars"},
				Usage:   "list all secrets for approle",
			},
		},
		Action: func(ctx *cli.Context) error {
			if err := validate(VEndpoint|VToken, v.cfg); err != nil {
				return err
			}

			v.ctx = ctx.Context

			if ctx.Bool("key") {
				return v.ListKey()
			}

			if ctx.Bool("policy") {
				return v.ListPolicies()
			}

			if ctx.Bool("user") {
				return v.ListUsers()
			}

			if ctx.Bool("approle") {
				return v.ListApprole()
			}

			if ctx.Bool("mount") {
				return v.ListMounts()
			}

			if ctx.String("approle-secrets") != "" {
				return v.ListApproleSecrets("approle-secrets")
			}

			cli.ShowAppHelp(ctx)
			return errors.New("must provide additional flag(s)")
		},
	}
}

func (v *VaultCmd) EnableSubCmd() *cli.Command {
	return &cli.Command{
		Name:  "enable",
		Usage: "enable methods or systems",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "secret",
				Aliases: []string{"s"},
				Usage:   "enable secret engine",
			},
		},
		Subcommands: []*cli.Command{
			{
				Name: "totp",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "username",
						Aliases: []string{"u"},
					},
					&cli.BoolFlag{
						Name:    "with-meta",
						Aliases: []string{"m"},
					},
				},
				Action: func(ctx *cli.Context) error {
					if err := validate(VEndpoint|VToken, v.cfg); err != nil {
						return err
					}

					v.ctx = ctx.Context

					return v.EnableTOTP(ctx.String("username"), ctx.Bool("with-meta"))
				},
			},
		},
		Action: func(ctx *cli.Context) error {
			if err := validate(VEndpoint|VToken, v.cfg); err != nil {
				return err
			}

			v.ctx = ctx.Context

			if ctx.Bool("secret") {
				return v.EnableMount()
			}

			cli.ShowAppHelp(ctx)
			return errors.New("must provide additional subcommand(s)")
		},
	}
}

func (v *VaultCmd) DisableSubCmd() *cli.Command {
	return &cli.Command{
		Name:  "disable",
		Usage: "disable methods or systems",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "secret",
				Aliases: []string{"s"},
				Usage:   "disable secret engine",
			},
		},
		Action: func(ctx *cli.Context) error {
			if err := validate(VEndpoint|VToken, v.cfg); err != nil {
				return err
			}

			v.ctx = ctx.Context

			if ctx.Bool("secret") {
				return v.DisableMount()
			}

			cli.ShowAppHelp(ctx)
			return errors.New("must provide additional subcommand(s)")
		},
	}
}

func (v *VaultCmd) RemoveSubCmd() *cli.Command {
	return &cli.Command{
		Name:  "remove",
		Usage: "remove methods or secrets",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "approle",
				Aliases: []string{"ar"},
				Usage:   "remove vault approle",
			},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.String("approle") != "" {
				return v.RemoveApprole(ctx.String("approle"))
			}

			return nil
		},
	}
}
