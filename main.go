package main

import (
	"os"
	"strings"
	"time"

	"github.com/zerodoctor/zdcli/logger"
	"github.com/zerodoctor/zdcli/tui/comp"
)

func main() {
	logger.Init()

	vargs := os.Args
	if len(vargs) > 2 {
		cmd := strings.Join(vargs[1:], " ")
		logger.Info("exec:", cmd)

		if len(vargs) > 3 {
			switch vargs[1] {
			case "--edit":
				cmd = strings.Join(vargs[2:], " ")
				StartEdit(cmd)

				return
			}
		}

		StartLua(cmd)
		return
	}

	running := true
	for running {
		exit := StartTui()

		switch exit.Code {
		case comp.EXIT_EDT:
			StartEdit(exit.Msg)
			time.Sleep(100 * time.Millisecond)
			continue

		case comp.EXIT_CMD:

		case comp.EXIT_LUA:
			StartLua(exit.Msg)
			time.Sleep(100 * time.Millisecond)
			continue

		}

		running = false
	}

	logger.Print("Good Bye.")
}
