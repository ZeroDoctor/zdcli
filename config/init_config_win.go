//go:build windows && !linux && !freebsd && !netbsd && !openbsd && !darwin && !js
// +build windows,!linux,!freebsd,!netbsd,!openbsd,!darwin,!js

package config

import "github.com/zerodoctor/zdcli/util"

func Init() *Config {
	return &Config{
		ShellCmd:       "cmd /c",
		EditorCmd:      "nvim",
		LuaCmd:         "lua",
		LuaDownloadURL: "https://sourceforge.net/projects/luabinaries/files/5.4.2/Tools%20Executables/lua-5.4.2_Win64_bin.zip",
		RootScriptDir:  util.EXEC_PATH + "/lua",
		ServerEndPoint: "https://api.zerodoc.dev",
		VaultEndpoint:  "https://vault.zerodoc.dev",
	}
}
