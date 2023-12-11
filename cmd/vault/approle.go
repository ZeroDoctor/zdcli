package vault

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/zerodoctor/zdcli/logger"
	"github.com/zerodoctor/zdcli/util"
	"github.com/zerodoctor/zdtui/ui"
)

type AppRole struct {
	RoleID   string `json:"role_id"`
	SecretID string `json:"secret_id"`
}

func (v *VaultCmd) NewApprole(roleName string, withTokenSettings, withSecretSettings, createSecret bool) error {
	req := schema.AppRoleWriteRoleRequest{}
	if withTokenSettings {
		respPolicy, err := v.client.System.PoliciesListAclPolicies(
			v.ctx, vault.WithToken(v.GetToken()),
		)
		if err != nil {
			return fmt.Errorf("failed to list policies [error=%s]", err.Error())
		}

		if len(respPolicy.Data.Keys) > 0 {
			logger.Infof("Available policies %+v", respPolicy.Data.Keys)
		}

		maxTtl := ui.NewTextInput()
		maxTtl.Input.Prompt = "Token Max TTL: "
		maxTtl.Input.Placeholder = "0"
		maxTtl.Input.SetValue("0")
		maxTtl.Input.Focus()

		policies := ui.NewTextInput()
		policies.Input.Prompt = "Policies (sep. ','): "
		policies.Input.Placeholder = ""
		policies.Input.SetValue("")

		numUses := ui.NewTextInput()
		numUses.Input.Prompt = "Number of Uses:"
		numUses.Input.Placeholder = "0"
		numUses.Input.SetValue("0")

		tokenType := ui.NewTextInput()
		tokenType.Input.Prompt = "Type:"
		tokenType.Input.Placeholder = "default"
		tokenType.Input.SetValue("default")

		form := ui.NewTextInputForm(maxTtl, policies, numUses, tokenType)
		if _, err := tea.NewProgram(form).Run(); err != nil {
			return fmt.Errorf("failed to start tea ui [error=%s]", err.Error())
		}
		if form.WasCancel {
			return nil
		}

		req.TokenMaxTtl = maxTtl.Input.Value()
		req.TokenPolicies = strings.Split(policies.Input.Value(), ",")
		uses, err := strconv.Atoi(numUses.Input.Value())
		if err != nil {
			return fmt.Errorf("failed to parse num of [uses=%s] [error=%s]",
				numUses.Input.Value(), err.Error(),
			)
		}
		req.TokenNumUses = int32(uses)
		req.TokenType = tokenType.Input.Value()
	}

	respAppRole, err := v.client.Auth.AppRoleWriteRole(
		v.ctx, roleName, req, vault.WithToken(v.GetToken()),
	)
	if err != nil {
		return fmt.Errorf("failed to create new approle [error=%s]", err.Error())
	}

	str, err := util.StructString(respAppRole)
	if err != nil {
		return err
	}
	fmt.Println(str)

	respRoleID, err := v.client.Auth.AppRoleReadRoleId(
		v.ctx, roleName, vault.WithToken(v.GetToken()),
	)
	if err != nil {
		return fmt.Errorf("failed to read role id [error=%s]", err.Error())
	}
	appRole := AppRole{
		RoleID: respRoleID.Data.RoleId,
	}

	if createSecret {
		appRole.SecretID, err = v.NewSecretID(roleName, withSecretSettings)
		if err != nil {
			return err
		}
	}

	str, err = util.StructString(appRole)
	if err != nil {
		return err
	}
	fmt.Println(str)

	return nil
}

func (v *VaultCmd) NewSecretID(roleName string, withSecretSettings bool) (string, error) {
	req := schema.AppRoleWriteSecretIdRequest{}
	if withSecretSettings {
		ttl := ui.NewTextInput()
		ttl.Input.Prompt = "Secret TTL: "
		ttl.Input.Placeholder = "0"
		ttl.Input.Focus()

		numUses := ui.NewTextInput()
		numUses.Input.Prompt = "0"

		form := ui.NewTextInputForm(ttl, numUses)
		if _, err := tea.NewProgram(form).Run(); err != nil {
			return "", fmt.Errorf("failed to start tea ui [error=%s]", err.Error())
		}
		if form.WasCancel {
			return "", nil
		}

		req.Ttl = ttl.Input.Value()
		uses, err := strconv.Atoi(numUses.Input.Value())
		if err != nil {
			return "", fmt.Errorf("failed to parse num of [uses=%s] [error=%s]",
				numUses.Input.Value(), err.Error(),
			)
		}
		req.NumUses = int32(uses)
	}

	respSecretID, err := v.client.Auth.AppRoleWriteSecretId(
		v.ctx, roleName, req, vault.WithToken(v.GetToken()),
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate new secret id [error=%s]", err.Error())
	}

	return respSecretID.Data.SecretId, nil
}

func (v *VaultCmd) ListApprole() error {
	resp, err := v.client.Auth.AppRoleListRoles(
		v.ctx, vault.WithToken(
			v.cfg.VaultTokens[v.cfg.VaultUser],
		),
	)
	if err != nil {
		return fmt.Errorf("failed to list approles [error=%s]", err.Error())
	}

	str, err := util.StructString(resp)
	if err != nil {
		return err
	}
	fmt.Println(str)

	return nil
}

func (v *VaultCmd) ListApproleSecrets(approle string) error {
	resp, err := v.client.Auth.AppRoleListSecretIds(
		v.ctx, approle, vault.WithToken(v.GetToken()),
	)
	if err != nil {
		return fmt.Errorf("failed to list [approle=%s] secrets [error=%s]", approle, err.Error())
	}

	str, err := util.StructString(resp)
	if err != nil {
		return err
	}
	fmt.Println(str)

	return nil
}

func (v *VaultCmd) RemoveApprole(approle string) error {
	resp, err := v.client.Auth.AppRoleDeleteRole(
		v.ctx, approle, vault.WithToken(v.GetToken()),
	)
	if err != nil {
		return fmt.Errorf("failed to delete [approle=%s] [error=%s]", approle, err.Error())
	}

	str, err := util.StructString(resp)
	if err != nil {
		return err
	}
	fmt.Println(str)

	return nil
}
