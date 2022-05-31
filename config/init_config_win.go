//go:build windows && !linux && !freebsd && !netbsd && !openbsd && !darwin && !js
// +build windows,!linux,!freebsd,!netbsd,!openbsd,!darwin,!js

package config

import "github.com/zerodoctor/zdcli/util"

func Init() *Config {
	return &Config{
		ShellCmd:      "cmd /c",
		EditorCmd:     "nvim",
		LuaCmd:        "lua",
		RootScriptDir: util.EXEC_PATH + "/lua",
	}
}
