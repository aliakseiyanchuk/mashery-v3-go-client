package v3client

import (
	"context"
	"errors"
	"fmt"
)

type ContextTokenProvider struct {
	header string
}

type contextKeyType string

const contextKey = contextKeyType("accesstoken.v3.client.mashery")

func (f *ContextTokenProvider) HeaderAuthorization(ctx context.Context) (map[string]string, error) {
	if tkn, err := f.AccessToken(ctx); err != nil {
		return nil, err
	} else {
		headerVal := fmt.Sprintf("%s %s", f.header, tkn)
		return map[string]string{
			"Authorization": headerVal,
		}, nil
	}
}

func (f *ContextTokenProvider) QueryStringAuthorization(_ context.Context) (map[string]string, error) {
	var noop map[string]string
	return noop, nil
}

func (f *ContextTokenProvider) AccessToken(ctx context.Context) (string, error) {
	if v := ctx.Value(contextKey); v != nil {
		return v.(string), nil
	} else {
		return "", errors.New("token is not supplied in the context")
	}
}

func (f *ContextTokenProvider) Close() {
	// Do nothing.
}

func ContextWithAccessToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, contextKey, token)
}

func AccessTokenFromContext(ctx context.Context) string {
	return ctx.Value(contextKey).(string)
}

func NewContextTokenProvider() V3AccessTokenProvider {
	return &ContextTokenProvider{
		header: "Bearer",
	}
}
