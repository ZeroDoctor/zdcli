package config

import (
	"fmt"
	"io/ioutil"

	"github.com/pelletier/go-toml/v2"
	"github.com/shirou/gopsutil/v3/host"
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
	data, err := toml.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(util.EXEC_PATH+"/zdconfig.toml", data, 0644)
}

func (c *Config) Load() error {
	data, err := ioutil.ReadFile(util.EXEC_PATH + "/zdconfig.toml")
	if err != nil {
		c = Init()
		return c.Save()
	}

	info, _ := host.Info()
	c.OS = info.OS
	c.Arch = info.KernelArch
	if c.Arch == "x86_64" {
		c.Arch = "amd64"
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
