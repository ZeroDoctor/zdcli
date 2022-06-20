package config

import (
	"io/ioutil"

	"github.com/pelletier/go-toml/v2"
	"github.com/zerodoctor/zdcli/util"
)

type Config struct {
	LuaCmd         string
	EditorCmd      string
	RootScriptDir  string
	ServerEndPoint string
	ShellCmd       string
	VaultEndpoint  string
	VaultTokens    map[string]string
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

	return toml.Unmarshal(data, c)
}
