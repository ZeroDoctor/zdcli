package main

import (
	"context"
	"os"
	"time"

	"github.com/zerodoctor/zdcli/command"
	"github.com/zerodoctor/zdcli/logger"
)

func StartCmd(cmd string) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()

	info := command.Info{
		Command: "lua build-app.lua " + cmd, // TODO: allow user to set lua endpoint
		Dir:     "./lua/",                   // TODO: allow user to set lua direcoty
		Ctx:     ctx,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Stdin:   os.Stdin,
	}

	err := command.Exec(&info)
	if err != nil {
		logger.Errorf("failed command [error=%s]", err.Error())
	}
}
