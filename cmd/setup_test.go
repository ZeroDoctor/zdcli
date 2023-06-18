package cmd

import (
	"context"
	"log"
	"testing"

	"github.com/zerodoctor/go-logging"
	"github.com/zerodoctor/zdcli/config"
	"github.com/zerodoctor/zdcli/logger"
	"github.com/zerodoctor/zdcli/util"
)

func TestDownloadLua(t *testing.T) {
	cfg := config.Init()
	logger.Init(logging.DEBUG)

	setup := &SetupCmd{}

	util.BIN_PATH = "."
	if err := setup.DownloadLua(context.Background(), cfg); err != nil {
		log.Fatal(err.Error())
	}
}
