package cmd

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
	"github.com/zerodoctor/zdcli/command"
	"github.com/zerodoctor/zdcli/config"
	"github.com/zerodoctor/zdcli/logger"
	"github.com/zerodoctor/zdtui/ui"
)

type SetupCmd struct{}

func NewSetupCmd(cfg *config.Config) *cli.Command {
	setup := &SetupCmd{}

	return &cli.Command{
		Name:  "setup",
		Usage: "setup lua, editor, and dir configs",
		Action: func(ctx *cli.Context) error {
			luaCmd := ui.NewTextInput()
			luaCmd.Input.Prompt = "Enter lua command: "
			luaCmd.Input.Placeholder = cfg.LuaCmd
			luaCmd.Input.Focus()

			editorCmd := ui.NewTextInput()
			editorCmd.Input.Prompt = "Enter editor command: "
			editorCmd.Input.Placeholder = cfg.EditorCmd

			shellCmd := ui.NewTextInput()
			shellCmd.Input.Prompt = "Enter shell command: "
			shellCmd.Input.Placeholder = cfg.ShellCmd

			serverEndpoint := ui.NewTextInput()
			serverEndpoint.Input.Prompt = "Enter server endpoint command: "
			serverEndpoint.Input.Placeholder = cfg.ServerEndPoint

			vaultEndpoint := ui.NewTextInput()
			vaultEndpoint.Input.Prompt = "Enter vault endpoint command: "
			vaultEndpoint.Input.Placeholder = cfg.VaultEndpoint

			form := ui.NewTextInputForm(
				luaCmd, editorCmd, shellCmd, serverEndpoint, vaultEndpoint,
			)

			if err := tea.NewProgram(form).Start(); err != nil {
				logger.Errorf("failed to start tea ui [error=%s]", err.Error())
				return nil
			}
			if form.WasCancel {
				return nil
			}

			cfg.LuaCmd = luaCmd.Input.Value()
			cfg.EditorCmd = editorCmd.Input.Value()
			cfg.ShellCmd = shellCmd.Input.Value()
			cfg.ServerEndPoint = serverEndpoint.Input.Value()
			cfg.VaultEndpoint = vaultEndpoint.Input.Value()

			logger.Infof("checking lua path...\n%s", cfg)
			if err := setup.CheckingLuaPath(cfg.LuaCmd); err != nil {
				logger.Errorf("failed to execute lua [error=%s]", err.Error())
				logger.Infof("would you like to download and install lua? (y/n):")
				// code for ui options
				// code for download and install if yes
			}

			logger.Infof("saving...\n%s", cfg)
			if err := cfg.Save(); err != nil {
				logger.Errorf("failed to save config [error=%s]", err.Error())
			}

			return nil
		},
	}
}

func (s *SetupCmd) CheckingLuaPath(lua string) error {
	info := &command.Info{
		Command: lua + " --version",
	}

	err := command.Exec(info)
	if err != nil {
		return err
	}

	if info.ErrBuffer != "" {
		return errors.New(info.ErrBuffer)
	}

	return nil
}
