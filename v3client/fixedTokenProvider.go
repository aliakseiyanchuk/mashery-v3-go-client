package v3client

// FixedTokenProvider Implementation of a fixed access token provider
type FixedTokenProvider struct {
	token string
}

func (f FixedTokenProvider) AccessToken() (string, error) {
	return f.token, nil
}

func (f *FixedTokenProvider) UpdateToken(tkn string) {
	f.token = tkn
}

func (f FixedTokenProvider) Close() {
	// Do nothing.
}

func NewFixedTokenProvider(tkn string) V3AccessTokenProvider {
	return FixedTokenProvider{token: tkn}
}
