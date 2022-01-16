package v3client

import "fmt"

// FixedTokenProvider Implementation of a fixed access token provider
type FixedTokenProvider struct {
	token  string
	header string
}

func (f FixedTokenProvider) HeaderAuthorization() (map[string]string, error) {
	return map[string]string{
		"Authorization": f.header,
	}, nil
}

func (f FixedTokenProvider) QueryStringAuthorization() (map[string]string, error) {
	var noop map[string]string
	return noop, nil
}

func (f FixedTokenProvider) AccessToken() (string, error) {
	return f.token, nil
}

func (f *FixedTokenProvider) UpdateToken(tkn string) {
	f.token = tkn
	f.header = fmt.Sprintf("Bearer %s", tkn)
}

func (f FixedTokenProvider) Close() {
	// Do nothing.
}

func NewFixedTokenProvider(tkn string) V3AccessTokenProvider {
	rv := FixedTokenProvider{}
	rv.UpdateToken(tkn)
	return rv
}
