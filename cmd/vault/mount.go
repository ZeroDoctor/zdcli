package vault

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/zerodoctor/zdcli/util"
	"github.com/zerodoctor/zdtui/ui"
)

func (v *Vault) EnableMount() error {
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

	req := schema.MountsEnableSecretsEngineRequest{
		Description: desc.Input.Value(),
		Type:        mtype.Input.Value(),
	}
	if mtype.Input.Value() == "kv" {
		req.Options = make(map[string]interface{})
		req.Options["version"] = "2"
	}

	resp, err := v.client.System.MountsEnableSecretsEngine(
		v.Ctx, mount.Input.Value(), req, vault.WithToken(v.GetToken()),
	)
	if err != nil {
		return fmt.Errorf("failed to enable secret engine [path=%s] [error=%s]", mount.Input.Value(), err.Error())
	}

	str, err := util.StructString(resp)
	if err != nil {
		return err
	}
	fmt.Println(str)

	return nil
}

func (v *Vault) DisableMount() error {
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

	resp, err := v.client.System.MountsDisableSecretsEngine(
		v.Ctx, mount.Input.Value(), vault.WithToken(v.GetToken()),
	)
	if err != nil {
		return fmt.Errorf("failed to disable secret engine [mount=%s] [error=%s]", mount.Input.Value(), err.Error())
	}

	str, err := util.StructString(resp)
	if err != nil {
		return err
	}
	fmt.Println(str)

	return nil
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
