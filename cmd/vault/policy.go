package vault

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/zerodoctor/zdcli/logger"
	"github.com/zerodoctor/zdcli/util"
	"github.com/zerodoctor/zdtui/ui"
)

func (v *Vault) NewPolicy() error {
	tiPolicy := ui.NewTextInput()
	tiPolicy.Input.Prompt = "Enter policy name: "
	tiPolicy.Input.Placeholder = "user"
	tiPolicy.Focus()

	form := ui.NewTextInputForm(tiPolicy)
	if _, err := tea.NewProgram(form).Run(); err != nil {
		return fmt.Errorf("failed to start tea ui [error=%s]", err.Error())
	}
	if form.WasCancel {
		return nil
	}

	if tiPolicy.Input.Err != nil {
		return fmt.Errorf("failed to get input path [policy_error=%s]",
			tiPolicy.Input.Err.Error(),
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
		if split := strings.Split(fileName, "."); len(split) <= 1 || split[1] != "hcl" {
			continue
		}

		items = append(items, ui.NewItem(fileName, "", nil))
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
	logger.Infof("uploading [file=%s] to [policy=%s]...", fileName, tiPolicy.Input.Value())
	data, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("failed to read [file=%s] [error=%s]", fileName, err.Error())
	}

	req := schema.PoliciesWriteAclPolicyRequest{Policy: string(data)}
	resp, err := v.client.System.PoliciesWriteAclPolicy(
		v.ctx, tiPolicy.Input.Value(), req, vault.WithToken(
			v.cfg.VaultTokens[v.cfg.VaultUser],
		),
	)
	if err != nil {
		return fmt.Errorf("failed to write [policy=%s] [file=%s] [error=%s]",
			tiPolicy.Input.Value(), fileName, err.Error(),
		)
	}

	str, err := util.StructString(resp)
	if err != nil {
		return err
	}
	fmt.Println(str)

	return nil
}

func (v *Vault) GetPolicy(policyName string) error {
	if policyName == "" {
		tiPolicy := ui.NewTextInput()
		tiPolicy.Input.Prompt = "Enter policy name: "
		tiPolicy.Input.Placeholder = "user"
		tiPolicy.Focus()

		form := ui.NewTextInputForm(tiPolicy)
		if _, err := tea.NewProgram(form).Run(); err != nil {
			return fmt.Errorf("failed to start tea ui [error=%s]", err.Error())
		}
		if form.WasCancel {
			return nil
		}

		if tiPolicy.Input.Err != nil {
			return fmt.Errorf("failed to get input path [policy_error=%s]",
				tiPolicy.Input.Err.Error(),
			)
		}

		policyName = tiPolicy.Input.Value()
	}

	data, err := v.client.System.PoliciesReadAclPolicy(
		v.ctx,
		policyName,
		vault.WithToken(
			v.cfg.VaultTokens[v.cfg.VaultUser],
		),
	)
	if err != nil {
		return fmt.Errorf("failed to read policy [error=%s]", err.Error())
	}

	hclClean := strings.Join(strings.Split(data.Data.Policy, "\n{"), " {")
	fmt.Println(hclClean)

	return nil
}

func (v *Vault) ListPolicies() error {
	data, err := v.client.System.PoliciesListAclPolicies(
		v.ctx, vault.WithToken(
			v.cfg.VaultTokens[v.cfg.VaultUser],
		),
	)
	if err != nil {
		return fmt.Errorf("failed to list policies [error=%s]", err.Error())
	}

	str, err := util.StructString(data)
	if err != nil {
		return err
	}
	fmt.Println(str)

	return nil
}
