package vault

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/zerodoctor/zdcli/logger"
	"github.com/zerodoctor/zdcli/util"
	"github.com/zerodoctor/zdtui/ui"
)

func (v *Vault) LoginUser() error {
	user := ui.NewTextInput()
	user.Input.Prompt = "Enter username: "
	user.Input.Placeholder = "username"
	user.Input.Focus()

	pass := ui.NewTextInput(ui.WithTIPassword())
	pass.Input.Prompt = "Enter password: "
	pass.Input.Placeholder = "********"

	form := ui.NewTextInputForm(user, pass)
	if _, err := tea.NewProgram(form).Run(); err != nil {
		return fmt.Errorf("failed to start tea ui [error=%s]", err.Error())
	}
	if form.WasCancel {
		return nil
	}
	userName := user.Input.Value()

	if _, ok := v.cfg.VaultTokens[userName]; ok {
		if _, err := v.client.Auth.TokenRevoke(
			v.Ctx,
			schema.TokenRevokeRequest{
				Token: v.cfg.VaultTokens[userName],
			},
			vault.WithToken(v.cfg.VaultTokens[userName]),
		); err != nil {
			logger.Warnf("failed to revoke current token [error=%s]", err.Error())
			v.cfg.VaultTokens["failed-revoke-"+userName+"-"+util.RandString(8)] = v.cfg.VaultTokens[userName]
		}

		delete(v.cfg.VaultTokens, userName)
	}

	resp, err := v.client.Auth.UserpassLogin(
		v.Ctx,
		user.Input.Value(),
		schema.UserpassLoginRequest{
			Password: pass.Input.Value(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create vault token [error=%s]", err.Error())
	}

	if resp.Auth.MFARequirement == nil {
		v.cfg.VaultUser = userName
		v.SetToken(resp.Auth.ClientToken)
		if err := v.cfg.Save(); err != nil {
			return fmt.Errorf("failed to save vault [token=%s] and key [error=%s]", resp.Auth.ClientToken, err.Error())
		}

		return nil
	}

	fmt.Println()
	for i := range resp.Warnings {
		logger.Warnf("[warn=%s]", resp.Warnings[i])
	}

	mfaSelected := 0
	mfaRequirment := resp.Auth.MFARequirement
	if len(resp.Auth.MFARequirement.MFAConstraints[userName+"-mfa"].Any) > 1 {
		// TODO: implement mfa type choice

		logger.Warn("choosing mfa type not supported. Selected [type=%s]",
			mfaRequirment.MFAConstraints[userName+"-mfa"].Any[mfaSelected].Type,
		)
	}

	code := ui.NewTextInput()
	code.Input.Prompt = "Enter code: "
	code.Focus()

	for code.Input.Value() == "" {
		form = ui.NewTextInputForm(code)
		if _, err := tea.NewProgram(form).Run(); err != nil {
			return fmt.Errorf("failed to start tea ui [error=%s]", err.Error())
		}
		if form.WasCancel {
			return nil
		}
	}
	fmt.Println()

	reqMFAValidate := schema.MfaValidateRequest{
		MfaRequestId: mfaRequirment.MFARequestID,
		MfaPayload: map[string]interface{}{
			mfaRequirment.MFAConstraints[userName+"-mfa"].Any[mfaSelected].ID: []string{code.Input.Value()},
		},
	}
	respMFAValidate, err := v.client.System.MfaValidate(v.Ctx, reqMFAValidate)
	if err != nil {
		return fmt.Errorf("failed to validate mfa [error=%s]", err.Error())
	}

	v.cfg.VaultUser = userName
	v.SetToken(respMFAValidate.Auth.ClientToken)
	if err := v.cfg.Save(); err != nil {
		return fmt.Errorf("failed to save vault [token=%s] and key [error=%s]", respMFAValidate.Auth.ClientToken, err.Error())
	}

	return nil
}

func (v *Vault) RevokeSelf() error {
	if _, err := v.client.Auth.TokenRevoke(
		v.Ctx,
		schema.TokenRevokeRequest{
			Token: v.GetToken(),
		},
		vault.WithToken(v.GetToken()),
	); err != nil {
		v.cfg.VaultTokens["failed-revoke-"+util.RandString(8)] = v.cfg.VaultTokens[v.cfg.VaultUser]
		return fmt.Errorf("failed to revoke current token [error=%s]", err.Error())
	}

	delete(v.cfg.VaultTokens, v.cfg.VaultUser)

	if err := v.cfg.Save(); err != nil {
		return fmt.Errorf("failed to deleted vault token [error=%s]", err.Error())
	}

	return nil
}

func (v *Vault) EnableAuthMethod(path, desc, mtype string) (interface{}, error) {
	request := schema.AuthEnableMethodRequest{
		Description: desc,
		Type:        mtype,
	}

	resp, err := v.client.System.AuthEnableMethod(
		v.Ctx, path, request, vault.WithToken(v.GetToken()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to enable auth method [error=%w]", err)
	}

	return resp, nil
}
