package config

import (
	"fmt"
	"os"
	"reflect"

	"github.com/pelletier/go-toml/v2"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/zerodoctor/zdcli/logger"
	"github.com/zerodoctor/zdcli/util"
)

type Config struct {
	LuaCmd              string
	LuaDownloadURL      string
	RootLuaScriptDir    string
	PythonCmd           string
	RootPythonScriptDir string
	ScriptExec          string
	EditorCmd           string
	ServerEndPoint      string
	ShellCmd            string
	VaultEndpoint       string
	VaultUser           string
	VaultTokens         map[string]string
	SWFSMasterEndpoint  string
	SWFSFilerEndpoint   string

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
	return os.WriteFile(util.EXEC_PATH+"/zdconfig.toml", data, 0644)
}

func (c *Config) Load() error {
	info, _ := host.Info()
	c.OS = info.OS
	c.Arch = info.KernelArch
	if c.Arch == "x86_64" {
		c.Arch = "amd64"
	}

	dconf := Init()
	data, err := os.ReadFile(util.EXEC_PATH + "/zdconfig.toml")
	if err != nil {
		logger.Warnf("[error=%s] creating new config file", err.Error())
		c = dconf
		c.RootLuaScriptDir = util.EXEC_PATH + "/lua"
		return c.Save()
	}

	if err := toml.Unmarshal(data, c); err != nil {
		return err
	}

	c.SetDefaultValuesIfEmpty(dconf)
	return err
}

func (c *Config) SetDefaultValuesIfEmpty(defconf *Config) {
	dref := reflect.ValueOf(defconf).Elem()
	sref := reflect.ValueOf(c).Elem()
loop:
	for i := 0; i < sref.NumField(); i++ {
		switch sref.Field(i).Kind() {
		case reflect.String:
			if sref.Field(i).String() != "" {
				continue loop
			}

			dvalue := dref.FieldByName(sref.Type().Field(i).Name).String()
			sref.Field(i).SetString(dvalue)
			logger.Debugf("setting default [value=%s]", dvalue)
		case reflect.Map:
		}
	}
}

func (c *Config) String() string {
	return fmt.Sprintf(`
[LuaCmd=%s]
[RootLuaScriptDir=%s]
[PythonCmd=%s]
[RootPythonScriptDir=%s]
[EditorCmd=%s]
[ServerEndPoint=%s]
[ShellCmd=%s]
[VaultEndpoint=%s]
[VaultTokens=%s]
[SWFSMasterEndpoint=%s]
[SWFSFilerEndpoint=%s]
[OS=%s]
[Arch=%s]`,
		c.LuaCmd, c.RootLuaScriptDir,
		c.PythonCmd, c.RootPythonScriptDir, c.EditorCmd,
		c.ServerEndPoint, c.ShellCmd, c.VaultEndpoint,
		c.VaultTokens, c.SWFSMasterEndpoint, c.SWFSFilerEndpoint,
		c.OS, c.Arch,
	)
}
