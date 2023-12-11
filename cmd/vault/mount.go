package vault

import (
	"fmt"

	"github.com/hashicorp/vault-client-go"
	"github.com/zerodoctor/zdcli/util"
)

func (v *VaultCmd) ListMounts() error {
	resp, err := v.client.System.MountsListSecretsEngines(
		v.ctx, vault.WithToken(v.GetToken()),
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
