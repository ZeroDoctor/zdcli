package vault

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/zerodoctor/zdcli/util"
	"github.com/zerodoctor/zdtui/ui"
)

func (v *VaultCmd) NewKey(fileName string) error {
	path := ui.NewTextInput()
	path.Input.Prompt = "Enter path: "
	path.Input.Placeholder = "/secret/data/github"
	path.Focus()

	// TODO: create view editor in tui

	return nil
}

func (v *VaultCmd) GetKey() error {
	mount := ui.NewTextInput()
	mount.Input.Prompt = "Enter mount: "
	mount.Input.Placeholder = "sys/mount"
	mount.Input.Focus()

	path := ui.NewTextInput()
	path.Input.Prompt = "Enter path: "
	path.Input.Placeholder = "/secret/github"

	form := ui.NewTextInputForm(mount, path)
	if _, err := tea.NewProgram(form).Run(); err != nil {
		return fmt.Errorf("failed to start tea ui [error=%s]", err.Error())
	}

	if path.WasCancel {
		return nil
	}

	if path.Input.Err != nil || mount.Input.Err != nil {
		return fmt.Errorf("failed to get input path [path_error=%s] [mount_error=%s]",
			path.Input.Err.Error(), mount.Input.Err.Error(),
		)
	}

	var data *vault.Response[schema.KvV2ReadResponse]
	var err error
	if data, err = v.client.Secrets.KvV2Read(
		v.ctx, path.Input.Value(),
		vault.WithToken(v.cfg.VaultTokens[v.cfg.VaultUser]),
		vault.WithMountPath(mount.Input.Value()),
	); err != nil {
		return fmt.Errorf("failed to get secret [path=%s] [error=%s]", path.Input.Value(), err.Error())
	}

	str, err := util.StructString(data)
	if err != nil {
		return err
	}
	fmt.Println(str)

	return nil
}

func (v *VaultCmd) ListKey() error {
	mount := ui.NewTextInput()
	mount.Input.Prompt = "Enter mount: "
	mount.Input.Placeholder = "keys"
	mount.Input.Focus()

	path := ui.NewTextInput()
	path.Input.Prompt = "Enter path: "
	path.Input.Placeholder = "secret/github"

	form := ui.NewTextInputForm(mount, path)
	if _, err := tea.NewProgram(form).Run(); err != nil {
		return fmt.Errorf("failed to start tea ui [error=%s]", err.Error())
	}
	if form.WasCancel {
		return nil
	}

	if path.Input.Err != nil || mount.Input.Err != nil {
		return fmt.Errorf("failed to get input path [path_error=%s] [mount_error=%s]",
			path.Input.Err.Error(), mount.Input.Err.Error(),
		)
	}

	data, err := v.client.Secrets.KvV2List(
		v.ctx, path.Input.Value(),
		vault.WithToken(v.cfg.VaultTokens[v.cfg.VaultUser]),
		vault.WithMountPath(mount.Input.Value()),
	)
	if err != nil {
		return fmt.Errorf("failed to list secret folders [error=%s]", err.Error())
	}

	str, err := util.StructString(data)
	if err != nil {
		return err
	}
	fmt.Println(str)

	return nil
}
