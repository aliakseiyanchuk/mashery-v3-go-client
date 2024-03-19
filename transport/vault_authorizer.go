package transport

import (
	"context"
	"fmt"
)

type VaultAuthorizer struct {
	Authorizer

	auth map[string]string
}

func (va VaultAuthorizer) HeaderAuthorization(_ context.Context) (map[string]string, error) {
	return va.auth, nil
}
func (va VaultAuthorizer) QueryStringAuthorization(_ context.Context) (map[string]string, error) {
	return nil, nil
}

func (va VaultAuthorizer) Close() {
	// Do nothing
}

// NewVaultAuthorizer Create HashiCorp vault authorizer
func NewVaultAuthorizer(token VaultToken) Authorizer {
	rv := &VaultAuthorizer{
		auth: map[string]string{
			"X-Vault-Token": string(token),
		},
	}

	return rv
}

// NewBearerAuthorizer Bearer token authorization
func NewBearerAuthorizer(token string) Authorizer {
	rv := &VaultAuthorizer{
		auth: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", token),
		},
	}

	return rv
}
