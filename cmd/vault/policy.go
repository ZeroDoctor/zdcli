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

func (v *Vault) NewPolicy(name string, file string) (interface{}, error) {
	if name == "" {
		tiPolicy := ui.NewTextInput()
		tiPolicy.Input.Prompt = "Enter policy name: "
		tiPolicy.Input.Placeholder = "user"
		tiPolicy.Focus()

		form := ui.NewTextInputForm(tiPolicy)
		if _, err := tea.NewProgram(form).Run(); err != nil {
			return nil, fmt.Errorf("failed to start tea ui [error=%s]", err.Error())
		}
		if form.WasCancel {
			return nil, nil
		}

		if tiPolicy.Input.Err != nil {
			return nil, fmt.Errorf("failed to get input path [policy_error=%s]",
				tiPolicy.Input.Err.Error(),
			)
		}

		name = tiPolicy.Input.Value()
	}

	if file == "" {
		dir, err := os.ReadDir(".")
		if err != nil {
			return nil, fmt.Errorf("failed to read current directory [error=%s]", err.Error())
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
			return nil, fmt.Errorf("failed to start program [error=%s]", err.Error())
		}

		if li.WasCancel {
			return nil, nil
		}

		selectedItem := (li.List.Items()[li.List.Index()]).(*ui.Item)
		file = selectedItem.Title()
	}

	return v.newPolicy(name, file)
}

func (v *Vault) newPolicy(name, file string) (interface{}, error) {
	logger.Infof("uploading [file=%s] to [policy=%s]...", file, name)
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read [file=%s] [error=%s]", file, err.Error())
	}

	req := schema.PoliciesWriteAclPolicyRequest{Policy: string(data)}
	resp, err := v.client.System.PoliciesWriteAclPolicy(
		v.Ctx, name, req, vault.WithToken(
			v.cfg.VaultTokens[v.cfg.VaultUser],
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to write [policy=%s] [file=%s] [error=%s]",
			name, file, err.Error(),
		)
	}

	return resp, err
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
		v.Ctx,
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
		v.Ctx, vault.WithToken(
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
