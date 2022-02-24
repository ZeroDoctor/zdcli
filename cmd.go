package main

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/zerodoctor/zdcli/command"
	"github.com/zerodoctor/zdcli/logger"
)

func StartLua(cmd string) {
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

func StartEdit(cmd string) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()

	var cmdArr []string
	split := strings.Split(cmd, " ")
	for _, str := range split {
		if len(str) > 4 && str[len(str)-4:] != ".lua" {
			cmdArr = append(cmdArr, str+".lua")
			continue
		} else if len(str) < 4 {
			cmdArr = append(cmdArr, str+".lua")
			continue
		}

		cmdArr = append(cmdArr, str)
	}

	info := command.Info{
		Command: "nvim " + strings.Join(cmdArr, " "),
		Dir:     "./lua/scripts/",
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
