package v3client

// Implementation of a fixed access token provider
type FixedTokenProvider struct {
	token string
}

func (f FixedTokenProvider) AccessToken() (string, error) {
	return f.token, nil
}

func NewFixedTokenProvider(tkn string) V3AccessTokenProvider {
	rv := FixedTokenProvider{token: tkn}
	return &rv
}
