package temp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

type Temp struct {
	Endpoint string
}

func NewTempClient(endpoint string) *Temp {
	return &Temp{
		Endpoint: endpoint,
	}
}

type AppRoleWriteSecretIdResponse struct {
	SecretId         string `json:"secret_id,omitempty"`
	SecretIdAccessor string `json:"secret_id_accessor,omitempty"`
	SecretIdNumUses  int32  `json:"secret_id_num_uses,omitempty"`
	SecretIdTtl      int32  `json:"secret_id_ttl,omitempty"`
}

func (t *Temp) AppRoleWriteSecretId(ctx context.Context, roleName string, request schema.AppRoleWriteSecretIdRequest, token string) (*vault.Response[AppRoleWriteSecretIdResponse], error) {

	requestPath := "/v1/auth/{approle_mount_path}/role/{role_name}/secret-id"
	requestPath = strings.Replace(requestPath, "{"+"approle_mount_path"+"}", url.PathEscape("approle"), -1)
	requestPath = strings.Replace(requestPath, "{"+"role_name"+"}", url.PathEscape(roleName), -1)

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return nil, fmt.Errorf("could not encode request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, t.Endpoint+requestPath, bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Vault-Token", token)

	client := http.DefaultClient

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &vault.Response[AppRoleWriteSecretIdResponse]{}
	if err := json.Unmarshal(data, result); err != nil {
		return nil, err
	}

	return result, nil
}