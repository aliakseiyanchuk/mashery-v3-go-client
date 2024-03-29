package v2client

import "context"

type V2Authorizer struct {
	apiKey    string
	signature string
}

func (v *V2Authorizer) HeaderAuthorization(_ context.Context) (map[string]string, error) {
	return nil, nil
}

func (v *V2Authorizer) QueryStringAuthorization(_ context.Context) (map[string]string, error) {
	return map[string]string{
		"apikey": v.apiKey,
		"sig":    v.signature,
	}, nil
}

func (v *V2Authorizer) UpdateSignature(sig string) {
	v.signature = sig
}

func (v *V2Authorizer) Close() {
	// Nothing to do
}

// NewV2Authorizer Creates a new instance of the V2 Authorizer.
func NewV2Authorizer(key string) *V2Authorizer {
	return &V2Authorizer{
		apiKey: key,
	}
}
