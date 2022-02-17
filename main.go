package main

import (
	"github.com/zerodoctor/zdcli/logger"
)

func main() {
	logger.Init()
	StartTui()
	logger.Print("Good Bye.")
}
