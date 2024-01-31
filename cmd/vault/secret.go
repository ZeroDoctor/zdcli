package vault

import (
	"fmt"
	"os"
	"path"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/zerodoctor/zdcli/logger"
	"github.com/zerodoctor/zdcli/util"
	"github.com/zerodoctor/zdtui/ui"
)

func (v *Vault) NewKey() error {
	tiMount := ui.NewTextInput()
	tiMount.Input.Prompt = "Enter mount: "
	tiMount.Input.Placeholder = "key"
	tiMount.Focus()

	tiPath := ui.NewTextInput()
	tiPath.Input.Prompt = "Enter path: "
	tiPath.Input.Placeholder = "github"

	form := ui.NewTextInputForm(tiMount, tiPath)
	if _, err := tea.NewProgram(form).Run(); err != nil {
		return fmt.Errorf("failed to start tea ui [error=%s]", err.Error())
	}
	if form.WasCancel {
		return nil
	}

	if tiMount.Input.Err != nil {
		return fmt.Errorf("failed to get input [mount_error=%s]",
			tiMount.Input.Err.Error(),
		)
	}

	if tiPath.Input.Err != nil {
		return fmt.Errorf("failed to get input [path_error=%s]",
			tiPath.Input.Err.Error(),
		)
	}

	dir, err := os.ReadDir(".")
	if err != nil {
		return fmt.Errorf("failed to read current directory [error=%s]", err.Error())
	}

	items := []list.Item{}

	for i := range dir {
		if dir[i].IsDir() {
			continue
		}

		fileName := dir[i].Name()
		fileType := path.Ext(fileName)
		if fileType != ".json" && fileType != ".env" {
			continue
		}

		items = append(items, ui.NewItem(fileName, "", nil))
	}

	if len(items) <= 0 {
		return fmt.Errorf("json or env files not found")
	}

	li := ui.NewList(items, 0, 0)
	p := tea.NewProgram(li)
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to start program [error=%s]", err.Error())
	}
	if li.WasCancel {
		return nil
	}

	selectedItem := (li.List.Items()[li.List.Index()]).(*ui.Item)
	fileName := selectedItem.Title()
	logger.Infof("uploading [file=%s] to [path=%s]...", fileName, tiPath.Input.Value())
	data, err := util.ConvertEnvFile(fileName)
	if err != nil {
		return fmt.Errorf("failed to read [file=%s] [error=%s]", fileName, err.Error())
	}

	req := schema.KvV2WriteRequest{Data: data}
	resp, err := v.client.Secrets.KvV2Write(
		v.Ctx, tiPath.Input.Value(), req,
		vault.WithMountPath(tiMount.Input.Value()),
		vault.WithToken(
			v.cfg.VaultTokens[v.cfg.VaultUser],
		),
	)
	if err != nil {
		return fmt.Errorf("failed to write [path=%s] [file=%s] [error=%s]",
			tiPath.Input.Value(), fileName, err.Error(),
		)
	}

	str, err := util.StructString(resp)
	if err != nil {
		return err
	}
	fmt.Println(str)

	return nil
}

func (v *Vault) GetKey() error {
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
		v.Ctx, path.Input.Value(),
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

func (v *Vault) ListKey() error {
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
		v.Ctx, path.Input.Value(),
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
