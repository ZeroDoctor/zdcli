package config

import (
	"fmt"
	"io/ioutil"

	"github.com/pelletier/go-toml/v2"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/zerodoctor/zdcli/logger"
	"github.com/zerodoctor/zdcli/util"
)

type Config struct {
	LuaCmd             string
	EditorCmd          string
	RootScriptDir      string
	ServerEndPoint     string
	ShellCmd           string
	VaultEndpoint      string
	VaultTokens        map[string]string
	SWFSMasterEndpoint string
	SWFSFilerEndpoint  string

	OS   string
	Arch string
}

func (c *Config) Save() error {
	info, _ := host.Info()
	c.OS = info.OS
	c.Arch = info.KernelArch
	if c.Arch == "x86_64" {
		c.Arch = "amd64"
	}

	data, err := toml.Marshal(c)
	if err != nil {
		return err
	}

	logger.Infof("updating config [file=%s]", util.EXEC_PATH+"/zdconfig.toml")
	return ioutil.WriteFile(util.EXEC_PATH+"/zdconfig.toml", data, 0644)
}

func (c *Config) Load() error {
	info, _ := host.Info()
	c.OS = info.OS
	c.Arch = info.KernelArch
	if c.Arch == "x86_64" {
		c.Arch = "amd64"
	}

	data, err := ioutil.ReadFile(util.EXEC_PATH + "/zdconfig.toml")
	if err != nil {
		logger.Warnf("[error=%s] creating new config file", err.Error())
		c = Init()
		c.RootScriptDir = util.EXEC_PATH + "/lua"
		return c.Save()
	}

	return toml.Unmarshal(data, c)
}

func (c *Config) String() string {
	return fmt.Sprintf(`[LuaCmd=%s]
[EditorCmd=%s]
[RootScriptDir=%s]
[ServerEndPoint=%s]
[ShellCmd=%s]
[VaultEndpoint=%s]
[VaultTokens=%s]
[SWFSMasterEndpoint=%s]
[SWFSFilerEndpoint=%s]
[OS=%s]
[Arch=%s]`,
		c.LuaCmd, c.EditorCmd, c.RootScriptDir,
		c.ServerEndPoint, c.ShellCmd, c.VaultEndpoint,
		c.VaultTokens, c.SWFSMasterEndpoint, c.SWFSFilerEndpoint,
		c.OS, c.Arch,
	)
}
