package vault

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/zerodoctor/zdcli/util"
	"github.com/zerodoctor/zdtui/ui"
)

func (v *Vault) EnableMountInput() error {
	mount := ui.NewTextInput()
	mount.Input.Prompt = "Enter mount name: "
	mount.Input.Placeholder = "key"
	mount.Focus()

	mtype := ui.NewTextInput()
	mtype.Input.Prompt = "Enter type: "
	mtype.Input.Placeholder = "kv"

	desc := ui.NewTextInput()
	desc.Input.Prompt = "Enter description: "
	desc.Input.Placeholder = ""

	form := ui.NewTextInputForm(mount, mtype, desc)
	if _, err := tea.NewProgram(form).Run(); err != nil {
		return fmt.Errorf("failed to start tea ui [error=%s]", err.Error())
	}
	if form.WasCancel {
		return nil
	}

	if mount.Input.Err != nil {
		return fmt.Errorf("failed to get input [path_error=%s]",
			mount.Input.Err.Error(),
		)
	}

	if mtype.Input.Err != nil {
		return fmt.Errorf("failed to get input [type_error=%s]",
			mtype.Input.Err.Error(),
		)
	}

	if desc.Input.Err != nil {
		return fmt.Errorf("failed to get input [desc_error=%s]",
			mount.Input.Err.Error(),
		)
	}

	resp, err := v.EnableMount(mount.Input.Value(), desc.Input.Value(), mtype.Input.Value())
	if err != nil {
		return err
	}

	str, err := util.StructString(resp)
	if err != nil {
		return err
	}
	fmt.Println(str)

	return nil
}

func (v *Vault) EnableMount(name, desc, mtype string) (interface{}, error) {
	req := schema.MountsEnableSecretsEngineRequest{
		Description: desc,
		Type:        mtype,
	}
	if mtype == "kv" {
		req.Options = make(map[string]interface{})
		req.Options["version"] = "2"
	}

	resp, err := v.client.System.MountsEnableSecretsEngine(
		v.Ctx, name, req, vault.WithToken(v.GetToken()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to enable secret engine [path=%s] [error=%s]", name, err.Error())
	}

	return resp, nil
}

func (v *Vault) DisableMountInput() error {
	mount := ui.NewTextInput()
	mount.Input.Prompt = "Enter mount name: "
	mount.Input.Placeholder = "key"
	mount.Focus()

	form := ui.NewTextInputForm(mount)
	if _, err := tea.NewProgram(form).Run(); err != nil {
		return fmt.Errorf("failed to start tea ui [error=%s]", err.Error())
	}
	if form.WasCancel {
		return nil
	}

	if mount.Input.Err != nil {
		return fmt.Errorf("failed to get input [path_error=%s]",
			mount.Input.Err.Error(),
		)
	}

	resp, err := v.DisableMount(mount.Input.Value())
	if err != nil {
		return err
	}

	str, err := util.StructString(resp)
	if err != nil {
		return err
	}
	fmt.Println(str)

	return nil
}

func (v *Vault) DisableMount(name string) (interface{}, error) {
	resp, err := v.client.System.MountsDisableSecretsEngine(
		v.Ctx, name, vault.WithToken(v.GetToken()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to disable secret engine [mount=%s] [error=%s]", name, err.Error())
	}

	return resp, nil
}

func (v *Vault) ListMounts() error {
	resp, err := v.client.System.MountsListSecretsEngines(
		v.Ctx, vault.WithToken(v.GetToken()),
	)
	if err != nil {
		return err
	}

	str, err := util.StructString(resp)
	if err != nil {
		return err
	}
	fmt.Println(str)

	return nil
}
