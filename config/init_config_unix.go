//go:build linux || freebsd || netbsd || openbsd
// +build linux freebsd netbsd openbsd

package config

import "github.com/zerodoctor/zdcli/util"

func Init() *Config {
	return &Config{
		ShellCmd:       "bash -c", // might need to use -i
		EditorCmd:      "nvim",
		LuaCmd:         "lua",
		RootScriptDir:  util.EXEC_PATH + "/lua",
		ServerEndPoint: "https://api.zerodoc.dev",
		VaultEndpoint:  "https://vault.zerodoc.dev",
	}
}
