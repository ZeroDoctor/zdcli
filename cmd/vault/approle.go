package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/zerodoctor/zdcli/cmd/vault/temp"
	"github.com/zerodoctor/zdcli/logger"
	"github.com/zerodoctor/zdtui/ui"
)

type AppRole struct {
	RoleID     string                                            `json:"role_id"`
	SecretID   string                                            `json:"secret_id"`
	ReadRole   *vault.Response[temp.AppRoleReadRoleResponse]     `json:"read_role"`
	ReadRoleId *vault.Response[schema.AppRoleReadRoleIdResponse] `json:"read_role_id"`
}

func (v *Vault) GetApprole(roleName string) (interface{}, error) {
	var appRole AppRole

	respRole, err := v.tempClient.AppRoleReadRole(
		v.Ctx, roleName, v.GetToken(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to read [role_name=%s] [error=%s]", roleName, err.Error())
	}
	appRole.ReadRole = respRole

	respRoleID, err := v.client.Auth.AppRoleReadRoleId(
		v.Ctx, roleName, vault.WithToken(v.GetToken()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to read id [role_name=%s] [error=%s]", roleName, err.Error())
	}
	appRole.ReadRoleId = respRoleID

	return &appRole, nil
}

type ApproleRequest struct {
	Name    *string                             `json:"name,omitempty"`
	Approle *schema.AppRoleWriteRoleRequest     `json:"approle,omitempty"`
	Secret  *schema.AppRoleWriteSecretIdRequest `json:"secret,omitempty"`
}

func (v *Vault) NewApprole(roleName string, withTokenSettings, withSecretSettings, createSecret bool, file string) (interface{}, error) {
	var approleRequest ApproleRequest
	writeRequest := schema.AppRoleWriteRoleRequest{}
	if withTokenSettings && file == "" {
		respPolicy, err := v.client.System.PoliciesListAclPolicies(
			v.Ctx, vault.WithToken(v.GetToken()),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to list policies [error=%s]", err.Error())
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
			return nil, fmt.Errorf("failed to start tea ui [error=%s]", err.Error())
		}
		if form.WasCancel {
			return nil, nil
		}

		writeRequest.TokenMaxTtl = maxTtl.Input.Value()
		writeRequest.TokenPolicies = strings.Split(policies.Input.Value(), ",")
		uses, err := strconv.Atoi(numUses.Input.Value())
		if err != nil {
			return nil, fmt.Errorf("failed to parse num of [uses=%s] [error=%s]",
				numUses.Input.Value(), err.Error(),
			)
		}
		writeRequest.TokenNumUses = int32(uses)
		writeRequest.TokenType = tokenType.Input.Value()

		approleRequest = ApproleRequest{
			Approle: &writeRequest,
		}
	}

	if file != "" {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(data, &approleRequest); err != nil {
			return nil, err
		}
	}

	role, err := v.newApprole(approleRequest)
	if createSecret && file == "" {
		role.SecretID, err = v.NewSecretID(roleName, withSecretSettings)
		if err != nil {
			return nil, err
		}
	}

	return role, err
}

func (v *Vault) newApprole(approleRequest ApproleRequest) (*AppRole, error) {
	_, err := v.client.Auth.AppRoleWriteRole(
		v.Ctx, *approleRequest.Name, *approleRequest.Approle, vault.WithToken(v.GetToken()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create new approle [error=%s]", err.Error())
	}

	respRoleID, err := v.client.Auth.AppRoleReadRoleId(
		v.Ctx, *approleRequest.Name, vault.WithToken(v.GetToken()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to read role id [error=%s]", err.Error())
	}
	appRole := &AppRole{
		RoleID: respRoleID.Data.RoleId,
	}

	if approleRequest.Secret != nil {
		appRole.SecretID, err = v.newSecretID(*approleRequest.Name, *approleRequest.Secret)
		if err != nil {
			return appRole, err
		}
	}

	return appRole, nil
}

func (v *Vault) NewSecretID(roleName string, withSecretSettings bool) (string, error) {
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

	return v.newSecretID(roleName, req)
}

func (v *Vault) newSecretID(roleName string, secretRequest schema.AppRoleWriteSecretIdRequest) (string, error) {
	respSecretID, err := v.tempClient.AppRoleWriteSecretId(
		v.Ctx, roleName, secretRequest, v.GetToken(),
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate new secret id [error=%s]", err.Error())
	}

	return respSecretID.Data.SecretId, nil
}

func (v *Vault) ListApprole() (interface{}, error) {
	resp, err := v.client.Auth.AppRoleListRoles(
		v.Ctx, vault.WithToken(v.GetToken()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list approles [error=%s]", err.Error())
	}

	return resp, nil
}

func (v *Vault) ListApproleSecretAccessors(approle string) (interface{}, error) {
	resp, err := v.client.Auth.AppRoleListSecretIds(
		v.Ctx, approle, vault.WithToken(v.GetToken()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list [approle=%s] secrets [error=%s]", approle, err.Error())
	}

	return resp, nil
}

func (v *Vault) RemoveApprole(approle string) (interface{}, error) {
	resp, err := v.client.Auth.AppRoleDeleteRole(
		v.Ctx, approle, vault.WithToken(v.GetToken()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to delete [approle=%s] [error=%s]", approle, err.Error())
	}

	return resp, nil
}
