package transport

import "context"

type VaultAuthorizer struct {
	Authorizer

	vaultAuth map[string]string
}

func (va VaultAuthorizer) HeaderAuthorization(_ context.Context) (map[string]string, error) {
	return va.vaultAuth, nil
}
func (va VaultAuthorizer) QueryStringAuthorization(_ context.Context) (map[string]string, error) {
	return nil, nil
}

func (va VaultAuthorizer) Close() {
	// Do nothing
}

// NewVaultAuthorizer Create HashiCorp vault authorizer
func NewVaultAuthorizer(token string) Authorizer {
	rv := &VaultAuthorizer{
		vaultAuth: map[string]string{
			"X-Vault-Token": token,
		},
	}

	return rv
}
