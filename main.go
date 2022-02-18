package main

import (
	"os"
	"strings"
	"time"

	"github.com/zerodoctor/zdcli/logger"
	"github.com/zerodoctor/zdcli/tui"
)

func main() {
	logger.Init()

	vargs := os.Args
	if len(vargs) > 2 {
		cmd := strings.Join(vargs[1:], " ")
		logger.Info("exec:", cmd)
		StartCmd(cmd)
		return
	}

	running := true
	for running {
		exit := StartTui()

		switch exit.Code {
		case tui.EXIT_CMD:
			StartCmd(exit.Msg)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		running = false
	}

	logger.Print("Good Bye.")
}
