package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/TwiN/go-pastebin"
	"github.com/zerodoctor/zdcli/command"
	"github.com/zerodoctor/zdcli/config"
	"github.com/zerodoctor/zdcli/logger"
	"github.com/zerodoctor/zdcli/tui/ui"
	"github.com/zerodoctor/zdcli/util"
)

func CreateLua(name string, cfg *config.Config) {
	temp := `
local app = require('lib.app')

local script = app:extend()

function script:hello_world(env_type)
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

	StartEdit(name, cfg)
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

func StartLua(cmd string, cfg *config.Config) {
	info := command.Info{
		Command: cfg.ShellCmd, // TODO: allow user to set lua endpoint
		Args:    []string{cfg.LuaCmd + " build-app.lua " + cmd},
		Dir:     cfg.RootScriptDir, // TODO: allow user to set lua direcoty
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

func StartEdit(cmd string, cfg *config.Config) {
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
		Command: cfg.EditorCmd + " " + strings.Join(cmdArr, " "),
		Dir:     cfg.RootScriptDir + "/scripts/",
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

func StartLs(cfg *config.Config) {
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

func PasteBinUpload(paths []string) {
	dir, err := os.Getwd()
	if err != nil {
		logger.Errorf("failed to get working directory [error=%s]", err.Error())
		return
	}

	var files []*os.File

	for _, path := range paths {
		path = dir + "/" + path
		file, err := os.OpenFile(path, os.O_RDONLY, 0644)
		if err != nil {
			logger.Errorf("failed to read [file=%s] [error=%s]", path, err.Error())
			continue
		}

		files = append(files, file)
	}

	// TODO: integrate vault
	client, err := pastebin.NewClient(os.Getenv("PASTE_BIN_USER"), os.Getenv("PASTE_BIN_PASS"), os.Getenv("PASTE_BIN_KEY"))
	if err != nil {
		logger.Errorf("failed to create paste bin client [error=%s]", err.Error())
		return
	}

	for _, file := range files {
		content, err := ioutil.ReadAll(file)
		if err != nil {
			logger.Errorf("failed to read [file=%s] [error=%s]", file.Name(), err.Error())
			continue
		}

		name := file.Name()
		ftype := ""
		index := strings.LastIndex(name, ".")
		if index != -1 {
			ftype = name[index:]
			name = name[:index]
		}

		key, err := client.CreatePaste(
			pastebin.NewCreatePasteRequest(name, string(content), pastebin.ExpirationNever, pastebin.VisibilityPrivate, ftype),
		)
		if err != nil {
			logger.Errorf("failed to upload [file=%s] to pastebin [error=%s]", file.Name, err.Error())
			continue
		}

		logger.Info("created paste [key=%s]", key)
		file.Close()
	}
}

func PasteBinUpdate(paths []string) {
	dir, err := os.Getwd()
	if err != nil {
		logger.Errorf("failed to get working directory [error=%s]", err.Error())
		return
	}

	var files []*os.File

	for _, path := range paths {
		path = dir + "/" + path
		file, err := os.OpenFile(path, os.O_RDONLY, 0644)
		if err != nil {
			logger.Errorf("failed to read [file=%s] [error=%s]", path, err.Error())
			continue
		}

		files = append(files, file)
	}

	// TODO: integrate vault
	client, err := pastebin.NewClient(os.Getenv("PASTE_BIN_USER"), os.Getenv("PASTE_BIN_PASS"), os.Getenv("PASTE_BIN_KEY"))
	if err != nil {
		logger.Errorf("failed to create paste bin client [error=%s]", err.Error())
		return
	}

	content, err := client.GetAllUserPastes()
	if err != nil {
		logger.Errorf("failed to get all users pastes [error=%s]", err.Error())
		return
	}

	for _, paste := range content {
		logger.Infof("[paste=%+v]\n", paste)
	}
}
