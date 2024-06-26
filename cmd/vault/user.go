package vault

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/zerodoctor/zdcli/generate"
	"github.com/zerodoctor/zdcli/logger"
	"github.com/zerodoctor/zdcli/util"
	"github.com/zerodoctor/zdtui/ui"
)

func (v *Vault) NewUserInput(userName string) error {
	form := ui.NewTextInputForm()
	pass := ui.NewTextInput(ui.WithTIPassword())
	pass.Input.Prompt = "Enter password: "
	pass.Input.Placeholder = "*********"
	pass.Focus()

	passConfirm := ui.NewTextInput(ui.WithTIPassword())
	passConfirm.Input.Prompt = "Enter confirm password: "
	passConfirm.Input.Placeholder = "*********"

	form.Inputs = append(form.Inputs, pass, passConfirm)

	errs := []error{errors.New("validate password")}
	for len(errs) > 0 {
		errs = []error{}

		if _, err := tea.NewProgram(form).Run(); err != nil {
			return fmt.Errorf("failed to start tea ui [error=%s]", err.Error())
		}
		if form.WasCancel {
			return nil
		}

		err := validatePasswords(pass.Input.Value(), passConfirm.Input.Value())
		if err != nil {
			errs = append(errs, err)
		}

		if len(errs) > 0 {
			logger.Warnf("[warnings=\n%s\n]", formatErrors(errs, "\t\n"))
		}
	}

	password := pass.Input.Value()

	respPolicy, err := v.client.System.PoliciesListAclPolicies(
		v.Ctx, vault.WithToken(v.GetToken()),
	)
	if err != nil {
		return fmt.Errorf("failed to list policies [error=%s]", err.Error())
	}

	policies := []string{}
	if len(respPolicy.Data.Keys) > 0 {
		area := ui.NewTextArea(
			ui.WithTitle(fmt.Sprintf("Add policies per line from list [%+v]", respPolicy.Data.Keys)),
		)
		if _, err := tea.NewProgram(area).Run(); err != nil {
			return fmt.Errorf("failed to start program [error=%s]", err.Error())
		}

		policies = strings.Fields(area.Value())
	}

	resp, err := v.NewUser(userName, password, policies)
	if err != nil {
		return err
	}

	str, err := util.StructString(resp)
	if err != nil {
		return err
	}
	fmt.Println(str)

	return err
}

func (v *Vault) NewUser(name, password string, policies []string) (interface{}, error) {
	if err := validateUserName(name); err != nil {
		return nil, err
	}

	reqUserpass := schema.UserpassWriteUserRequest{
		Password:      password,
		TokenPolicies: policies,
	}

	respUserpass, err := v.client.Auth.UserpassWriteUser(
		v.Ctx, name, reqUserpass, vault.WithToken(v.GetToken()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create [user=%s] [error=%s]", name, err.Error())
	}

	return respUserpass, nil
}

func (v *Vault) UpdateUserPolicies(userName string) error {
	respPolicy, err := v.client.System.PoliciesListAclPolicies(
		v.Ctx, vault.WithToken(v.GetToken()),
	)
	if err != nil {
		return fmt.Errorf("failed to list policies [error=%s]", err.Error())
	}

	policies := []string{}
	if len(respPolicy.Data.Keys) > 0 {
		area := ui.NewTextArea(
			ui.WithTitle(fmt.Sprintf("Add policies per line from list [%+v]", respPolicy.Data.Keys)),
		)
		if _, err := tea.NewProgram(area).Run(); err != nil {
			return fmt.Errorf("failed to start program [error=%s]", err.Error())
		}

		policies = strings.Fields(area.Value())
	}

	req := schema.UserpassUpdatePoliciesRequest{
		TokenPolicies: policies,
	}
	resp, err := v.client.Auth.UserpassUpdatePolicies(
		v.Ctx, userName, req, vault.WithToken(v.GetToken()),
	)
	if err != nil {
		return fmt.Errorf("failed to update [user=%s] policies [error=%s]", userName, err.Error())
	}

	str, err := util.StructString(resp)
	if err != nil {
		return err
	}
	fmt.Println(str)

	return nil
}

func (v *Vault) GetUser(userName string) error {
	resp, err := v.client.Auth.UserpassReadUser(
		v.Ctx, userName, vault.WithToken(v.GetToken()),
	)
	if err != nil {
		return fmt.Errorf("failed to read [user=%s] [error=%s]",
			userName, err.Error(),
		)
	}

	str, err := util.StructString(resp)
	if err != nil {
		return err
	}
	fmt.Println(str)

	return nil
}

func (v *Vault) ListUsers() error {
	respList, err := v.client.Auth.UserpassListUsers(
		v.Ctx, vault.WithToken(v.GetToken()),
	)
	if err != nil {
		return fmt.Errorf("failed to print user list [error=%s]", err.Error())
	}

	str, err := util.StructString(respList)
	if err != nil {
		return err
	}
	fmt.Println(str)

	return nil
}

type Alias struct {
	accessor string
	userName string
	entityID string
	aliasID  string
}

func (v *Vault) NewAliasInput(userName string, withMeta bool) (Alias, error) {
	if userName == "" {
		user := ui.NewTextInput()
		user.Input.Prompt = "Enter username: "
		user.Input.Placeholder = "username"
		user.Input.Focus()

		form := ui.NewTextInputForm(user)
		if _, err := tea.NewProgram(form).Run(); err != nil {
			return Alias{}, fmt.Errorf("failed to start tea ui [error=%s]", err.Error())
		}
		if form.WasCancel {
			return Alias{}, nil
		}

		userName = user.Input.Value()
	}

	metaData := map[string]interface{}{}
	if withMeta {
		meta := ui.NewTextArea(
			ui.WithTitle("Enter MetaData Key/Value Pair split by ':'"),
		)
		if _, err := tea.NewProgram(meta).Run(); err != nil {
			return Alias{}, fmt.Errorf("failed to start tea ui [error=%s]", err.Error())
		}

		metaSplit := strings.Split(meta.Value(), "\n")
		for i := range metaSplit {
			data := strings.Split(metaSplit[i], ":")
			metaData[data[0]] = ""
			if len(data) > 1 {
				metaData[data[0]] = data[1]
			}
		}
	}

	alias, err := v.NewAlias(userName, metaData)
	if err != nil {
		return Alias{}, err
	}

	str, err := util.StructString(alias)
	if err != nil {
		return alias, err
	}
	fmt.Println(str)

	return alias, nil
}

func (v *Vault) NewAlias(userName string, metaData map[string]interface{}) (Alias, error) {
	alias := Alias{userName: userName}

	userpassConfig, err := v.client.System.AuthReadConfiguration(
		v.Ctx, "userpass", vault.WithToken(v.GetToken()),
	)
	if err != nil {
		return Alias{}, fmt.Errorf("failed to read auth config [error=%s]", err.Error())
	}
	alias.accessor = userpassConfig.Data.Accessor

	req := schema.EntityCreateRequest{
		Name:     alias.userName + "-alias",
		Metadata: metaData,
	}
	resp, err := v.client.Identity.EntityCreate(
		v.Ctx, req, vault.WithToken(
			v.cfg.VaultTokens[v.cfg.VaultUser],
		),
	)
	if err != nil {
		return alias, fmt.Errorf("failed to create entity [username=%s] [error=%s]",
			alias.userName, err.Error(),
		)
	}

	str, _ := util.StructString(resp)

	var ok bool
	alias.entityID, ok = resp.Data["id"].(string)
	if !ok {
		return alias, fmt.Errorf("failed to create alias for [user=%s] [resp=%s]",
			userName, str,
		)
	}

	logger.Infof("creating alias with [userpass_accessor_id=%s] and [entity_id=%s]...",
		alias.accessor, alias.entityID,
	)

	reqAlias := schema.EntityCreateAliasRequest{
		Name:          alias.userName,
		CanonicalId:   alias.entityID,
		MountAccessor: alias.accessor,
	}
	entityAlias, err := v.client.Identity.EntityCreateAlias(
		v.Ctx, reqAlias, vault.WithToken(v.GetToken()),
	)
	if err != nil {
		return alias, fmt.Errorf("failed to create entity alias [error=%s]", err.Error())
	}

	if entityAlias != nil && entityAlias.Data != nil {
		alias.aliasID = entityAlias.Data["id"].(string)
	}

	return alias, nil
}

func (v *Vault) EnableTOTPInput(userName string, withMeta bool) error {
	alias, err := v.NewAliasInput(userName, withMeta)
	if err != nil {
		return err
	}

	resp, err := v.EnableTOTP(alias)
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

func (v *Vault) EnableTOTP(alias Alias) (interface{}, error) {
	list, err := v.client.Identity.MfaListTotpMethods(
		v.Ctx, vault.WithToken(v.GetToken()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list methods [error=%s]\n failed to find mfa method id. ask admin to create one", err.Error())
	}

	if len(list.Data.Keys) <= 0 {
		return nil, fmt.Errorf("failed to find mfa method id. ask admin to create one")
	}

	methodID := list.Data.Keys[0]

	logger.Infof("creating totp mfa with [method_id=%s] and [entity_id=%s]...", methodID, alias.entityID)
	adminGenReq := schema.MfaAdminGenerateTotpSecretRequest{
		EntityId: alias.entityID,
		MethodId: methodID,
	}
	adminGenResp, err := v.client.Identity.MfaAdminGenerateTotpSecret(
		v.Ctx, adminGenReq, vault.WithToken(v.GetToken()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate totp secret [error=%s]", err.Error())
	}

	err = generate.TOTP(alias.userName, adminGenResp.Data["barcode"].(string))
	if err != nil {
		return nil, fmt.Errorf("failed to generate qr code [error=%s]", err.Error())
	}

	loginEnfReq := schema.MfaWriteLoginEnforcementRequest{
		AuthMethodAccessors: []string{alias.accessor},
		MfaMethodIds:        []string{methodID},
	}
	loginEnfResp, err := v.client.Identity.MfaWriteLoginEnforcement(
		v.Ctx, alias.userName+"-mfa", loginEnfReq, vault.WithToken(v.GetToken()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create login enforcement [error=%s]", err.Error())
	}

	return loginEnfResp, nil
}
