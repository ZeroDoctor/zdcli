package main

import (
	"fmt"

	"github.com/zerodoctor/zdcli/command"
	"github.com/zerodoctor/zdcli/logger"
)

func PrintCommand(msg []byte) (int, error) {
	fmt.Printf("%s", string(msg))
	return len(msg), nil
}

func main() {
	logger.Init()
	info := command.Info{
		Command: "curl https://www.zerodoc.dev",

		OutBufSize: 100,
		OutFunc:    PrintCommand,
	}

	command.Exec(&info)
}
