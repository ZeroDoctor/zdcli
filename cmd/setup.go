package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
	"github.com/zerodoctor/zdcli/command"
	"github.com/zerodoctor/zdcli/config"
	"github.com/zerodoctor/zdcli/logger"
	"github.com/zerodoctor/zdcli/util"
	zdutil "github.com/zerodoctor/zdgo-util"
	"github.com/zerodoctor/zdtui/ui"
)

type SetupCmd struct{}

func NewSetupCmd(cfg *config.Config) *cli.Command {
	setup := &SetupCmd{}

	return &cli.Command{
		Name:  "setup",
		Usage: "setup lua, editor, and dir configs",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "list",
				Aliases: []string{"ls"},
			},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.Bool("list") {
				configPath := util.EXEC_PATH + "/zdconfig.toml"
				if !zdutil.FileExists(configPath) {
					logger.Warnf("failed to find config file [path=%s]", configPath)
					return nil
				}

				data, err := os.ReadFile(configPath)
				if err != nil {
					logger.Errorf("failed to read config file [path=%s] [error=%s]", configPath, err.Error())
					return nil
				}
				fmt.Println(string(data))

				return nil
			}

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

			if _, err := tea.NewProgram(form).Run(); err != nil {
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
				setup.AskToDownloadLua(ctx.Context, cfg)
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

func (s *SetupCmd) AskToDownloadLua(ctx context.Context, cfg *config.Config) error {
	logger.Infof("would you like to download and install lua? (Y/n):")

	downloadLuaInput := ui.NewTextInput()
	downloadLuaInput.Input.Placeholder = "y"
	downloadLuaInput.Focus()

	form := ui.NewTextInputForm(downloadLuaInput)
	if _, err := tea.NewProgram(form).Run(); err != nil {
		logger.Errorf("failed to start tea ui [error=%s]", err.Error())
		return nil
	}
	if form.WasCancel {
		return nil
	}

	if strings.ToLower(downloadLuaInput.Input.Value()) == "y" {
		if err := s.DownloadLua(ctx, cfg); err != nil {
			logger.Errorf("failed download and install lua [error=%s]", err.Error())
			return nil
		}
	}

	return nil
}

func (s *SetupCmd) DownloadLua(ctx context.Context, cfg *config.Config) error {
	logger.Infof("download lua from [url=%s]", cfg.LuaDownloadURL)

	req, err := http.NewRequest(http.MethodGet, cfg.LuaDownloadURL, nil)
	if err != nil {
		return err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	cfg.LuaCmd = "lua54"
	file := "lua.tar.gz"
	if runtime.GOOS == "windows" {
		file = "lua.zip"
	}

	if err := util.ExtractFromHttpResponse(ctx, file, util.BIN_PATH, resp.Body); err != nil {
		logger.Warnf("failed to extract response [error=%s]", err.Error())
	}

	if _, err := util.FollowDownloadRedirection(cfg.LuaDownloadURL, resp, func(resp *http.Response) error {
		if err := util.ExtractFromHttpResponse(ctx, file, util.BIN_PATH, resp.Body); err != nil {
			logger.Warnf("failed to extract response [error=%s]", err.Error())
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
