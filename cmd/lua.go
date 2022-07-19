package cmd

import (
	"context"
	"fmt"
	"io/fs"
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
				NewLua(arg, cfg)
			}

			return nil
		},
	}
}

func NewLua(name string, cfg *config.Config) {
	temp := `
local app = require('lib.app')
local util = require('lib.util')

local script = app:extend()

function script:hello_world()
	print('hello world!')
end

return script
`

	name = strings.ReplaceAll(name, ".", "/")
	if i := strings.LastIndex(name, ".lua"); i == -1 {
		name += ".lua"
	}

	path := cfg.RootScriptDir + "/scripts/" + name

	index := strings.LastIndex(path, "/")
	if !util.FolderExists(path[:index]) {
		os.MkdirAll(path[:index], 0644)
	}

	err := ioutil.WriteFile(path, []byte(temp), 0644)
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

	name = strings.ReplaceAll(name, ".", "/")
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
			str = strings.ReplaceAll(str, ".", "/")
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

type File struct {
	path string
	rel  string
	fs.FileInfo
}

func NewFiles(path string, rel string, files ...fs.FileInfo) []File {
	var fs []File

	for i := range files {
		fs = append(fs, File{path: path, rel: rel, FileInfo: files[i]})

	}

	return fs
}

func ListLua(cfg *config.Config) {
	path := cfg.RootScriptDir + "/scripts"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		logger.Errorf("failed ls [error=%s]", err.Error())
		return
	}

	var data [][]interface{} // data needed to preset as a tui table

	// recursively loop through all directories
	folder := util.NewStack(NewFiles(path, "", files...)...)
	for folder.Len() > 0 {
		file := *folder.Pop()

		if file.IsDir() {
			p := path + "/" + file.Name()
			if files, err = ioutil.ReadDir(p); err != nil {
				logger.Errorf("failed to read [dir=%s] [error=%s]", path, err.Error())
				continue
			}

			r := file.Name()
			if file.rel != "" {
				r = file.rel + "." + file.Name()
			}

			folder.Push(NewFiles(p, r, files...)...)
			continue
		}

		fileName := file.Name()
		if file.rel != "" {
			fileName = file.rel + "." + file.Name()
		}

		data = append(data, []interface{}{file.Mode(), fileName, file.Size(), file.ModTime()})
	}

	table, err := ui.NewTable([]string{"Mode", "Name", "Size", "Modify Time"}, data, 0, 0)
	if err != nil {
		logger.Errorf("failed to create ls table [error=%s]", err.Error())
		return
	}

	fmt.Println(table.View())
}
