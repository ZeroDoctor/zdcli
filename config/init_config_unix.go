//go:build linux || freebsd || netbsd || openbsd
// +build linux freebsd netbsd openbsd

package config

import "github.com/zerodoctor/zdcli/util"

func Init() *Config {
	return &Config{
		LuaCmd:              "lua",
		RootLuaScriptDir:    util.EXEC_PATH + "/lua",
		LuaDownloadURL:      "https://sourceforge.net/projects/luabinaries/files/5.4.2/Tools%20Executables/lua-5.4.2_Linux54_64_bin.tar.gz",
		PythonCmd:           "python",
		RootPythonScriptDir: util.EXEC_PATH + "/python",
		EditorCmd:           "nvim",
		ShellCmd:            "bash -c",
		ServerEndPoint:      "https://api.zerodoc.dev",
		VaultEndpoint:       "https://vault.zerodoc.dev",
	}
}
