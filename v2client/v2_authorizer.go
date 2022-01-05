package v2client

type V2Authorizer struct {
	apiKey    string
	signature string
}

func (v *V2Authorizer) Authorization() (map[string]string, error) {
	return map[string]string{
		"apikey": v.apiKey,
		"sig":    v.signature,
	}, nil
}

func (v *V2Authorizer) UpdateSignature(sig string) {
	v.signature = sig
}

func (v V2Authorizer) Close() {
	// Nothing to do
}

// NewV2Authorizer Creates a new instance of the V2 Authorizer.
func NewV2Authorizer(key string) *V2Authorizer {
	return &V2Authorizer{
		apiKey: key,
	}
}
