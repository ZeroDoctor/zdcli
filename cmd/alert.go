package cmd

import (
	"context"
	"errors"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/zerodoctor/zdcli/alert"
)

func TimerSubCmd() *cli.Command {
	return &cli.Command{
		// TODO: implement timer code
	}
}

func EndpointSubCmd() *cli.Command {
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
			},
		},

		Action: func(ctx *cli.Context) error {
			c, cancel := context.WithCancel(ctx.Context)
			defer cancel()

			checkDur := 5 * time.Second
			if sec := ctx.Int("check-duration"); sec > 0 {
				checkDur = time.Duration(sec) * time.Second
			}

			a := alert.WatchEndpoint(
				c,
				ctx.String("route"),
				ctx.String("message"),
				checkDur,
			)
			a.Wait()

			return nil
		},
	}
}

func AlertCmd() *cli.Command {
	return &cli.Command{
		Name:  "alert",
		Usage: "notifies user when an event happens",
		Subcommands: []*cli.Command{
			EndpointSubCmd(),
		},

		Action: func(ctx *cli.Context) error {
			cli.ShowAppHelp(ctx)
			return errors.New("must provide additional subcommand(s)")
		},
	}
}
