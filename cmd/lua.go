package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/zerodoctor/zdcli/command"
	"github.com/zerodoctor/zdcli/config"
	"github.com/zerodoctor/zdcli/logger"
	"github.com/zerodoctor/zdcli/tui/ui"
	"github.com/zerodoctor/zdcli/util"
)

func NewLuaCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:    "new",
		Aliases: []string{"n"},
		Usage:   "create a new lua script",
		Action: func(ctx *cli.Context) error {
			for _, arg := range ctx.Args().Slice() {
				CreateLua(arg, cfg)
			}

			return nil
		},
	}
}

func CreateLua(name string, cfg *config.Config) {
	temp := `
local app = require('lib.app')
local util = require('lib.util')

local script = app:extend()

function script:hello_world()
	print('hello world!')
end

return script
`

	if i := strings.LastIndex(name, ".lua"); i == -1 {
		name += ".lua"
	}

	err := ioutil.WriteFile(cfg.RootScriptDir+"/scripts/"+name, []byte(temp), 0644)
	if err != nil {
		logger.Errorf("failed to write lua script template [error=%s]", err.Error())
		return
	}

	EditLua(name, cfg)
}
func RemoveLuaCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:    "remove",
		Aliases: []string{"rm"},
		Usage:   "remove a lua script or a directory",
		Action: func(ctx *cli.Context) error {
			for _, arg := range ctx.Args().Slice() {
				RemoveLua(arg, cfg)
			}

			return nil
		},
	}
}
func RemoveLua(name string, cfg *config.Config) {

	if util.FolderExists(cfg.RootScriptDir + "/scripts/" + name) {
		os.RemoveAll(cfg.RootScriptDir + "/scripts/" + name)

		return
	}

	if i := strings.LastIndex(name, ".lua"); i == -1 {
		name += ".lua"
	}

	if err := os.Remove(cfg.RootScriptDir + "/scripts/" + name); err != nil {
		logger.Errorf("[error=%s]", err.Error())
	}
}

func EditLuaCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:    "edit",
		Aliases: []string{"e"},
		Usage:   "edits a lua script",
		Action: func(ctx *cli.Context) error {
			EditLua(ctx.Args().Get(0), cfg)
			return nil
		},
	}
}

func EditLua(cmd string, cfg *config.Config) {
	var cmdArr []string
	split := strings.Split(cmd, " ")
	for _, str := range split {
		if len(str) >= 4 && str[len(str)-4:] != ".lua" {
			cmdArr = append(cmdArr, str+".lua")
			continue
		} else if len(str) < 4 {
			cmdArr = append(cmdArr, str+".lua")
			continue
		}

		cmdArr = append(cmdArr, str)
	}

	info := command.Info{
		Command: cfg.EditorCmd + " ./scripts/" + strings.Join(cmdArr, " "),
		Dir:     cfg.RootScriptDir,
		Ctx:     context.Background(),
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Stdin:   os.Stdin,
	}

	err := command.Exec(&info)
	if err != nil {
		logger.Errorf("failed command [error=%s]", err.Error())
	}
}

func ListLuaCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "list current lua scripts",
		Action: func(ctx *cli.Context) error {
			ListLua(cfg)
			return nil
		},
	}
}

func ListLua(cfg *config.Config) {
	path := cfg.RootScriptDir + "/scripts"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		logger.Errorf("failed ls [error=%s]", err.Error())
		return
	}

	var data [][]interface{}

	for _, file := range files {
		data = append(data, []interface{}{file.Mode(), file.Name(), file.Size(), file.ModTime()})
	}

	table, err := ui.NewTable([]string{"Mode", "Name", "Size", "Modify Time"}, data, 0, 0)
	if err != nil {
		logger.Errorf("failed to create ls table [error=%s]", err.Error())
		return
	}

	fmt.Println(table.View())
}
