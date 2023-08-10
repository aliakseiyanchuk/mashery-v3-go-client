package v3client

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContextProviderWillReturnTokenFromKey(t *testing.T) {
	provider := NewContextTokenProvider()

	passCtx := ContextWithAccessToken(context.TODO(), "abc")
	at, err := provider.AccessToken(passCtx)

	assert.Equal(t, "abc", at)
	assert.Nil(t, err)
}

func TestContextProviderWillReturnHeaderAuthorization(t *testing.T) {
	provider := NewContextTokenProvider()

	passCtx := ContextWithAccessToken(context.TODO(), "abc")
	at, err := provider.HeaderAuthorization(passCtx)

	assert.Equal(t, "Bearer abc", at["Authorization"])
	assert.Nil(t, err)
}

func TestContextProviderWillRejectIfTokenIsMissing(t *testing.T) {
	provider := NewContextTokenProvider()

	at, err := provider.AccessToken(context.TODO())

	assert.Equal(t, "", at)
	assert.NotNil(t, err)
	assert.Equal(t, "token is not supplied in the context", err.Error())
}

func TestContextProviderWillRejectHeaderAuthorizationOnMissingToken(t *testing.T) {
	provider := NewContextTokenProvider()

	at, err := provider.HeaderAuthorization(context.TODO())

	assert.Nil(t, at)
	assert.NotNil(t, err)
	assert.Equal(t, "token is not supplied in the context", err.Error())
}
