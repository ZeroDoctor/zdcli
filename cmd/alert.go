package cmd

import (
	"context"
	"errors"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/zerodoctor/zdcli/alert"
)

type AlertCmd struct{}

func NewAlertCmd() *cli.Command {
	alert := &AlertCmd{}

	return &cli.Command{
		Name:  "alert",
		Usage: "notifies user when an event happens",
		Subcommands: []*cli.Command{
			alert.EndpointSubCmd(),
		},

		Action: func(ctx *cli.Context) error {
			cli.ShowAppHelp(ctx)
			return errors.New("must provide additional subcommand(s)")
		},
	}
}

func (a *AlertCmd) TimerSubCmd() *cli.Command {
	return &cli.Command{
		// TODO: implement timer code
	}
}

func (a *AlertCmd) EndpointSubCmd() *cli.Command {
	return &cli.Command{
		Name:    "endpoint",
		Aliases: []string{"e"},
		Usage:   "create an alert when endpoint fails or returns status code not between 200-299",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "route",
				Usage:    "an endpoint i.e. https://google.com",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "message",
				Usage:    "a message to display when route/endpoint fails",
				Required: true,
			},
			&cli.IntFlag{
				Name:    "check-duration",
				Aliases: []string{"c"},
				Usage:   "period (in seconds) alert checks endpoint",
				Value:   5,
			},
			&cli.BoolFlag{
				Name:    "once",
				Aliases: []string{"o"},
				Usage:   "only alert once then exit",
			},
		},

		Action: func(ctx *cli.Context) error {
			c, cancel := context.WithCancel(ctx.Context)
			defer cancel()

			sec := ctx.Int("check-duration")
			checkDur := time.Duration(sec) * time.Second

			params := alert.WatchEndpointParams{
				Ctx:         c,
				HealthRoute: ctx.String("route"),
				Message:     ctx.String("message"),
				CheckDur:    checkDur,
				Once:        ctx.Bool("once"),
			}

			a := alert.WatchEndpoint(params)
			a.Wait()

			return nil
		},
	}
}
