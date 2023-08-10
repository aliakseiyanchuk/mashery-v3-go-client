package v3client_test

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"testing"
)

func TestFixedTokenProvider(t *testing.T) {
	p := v3client.NewFixedTokenProvider("ABCD")
	if p == nil {
		t.Errorf("Fixed token provider must be returned")
		t.FailNow()
	}

	if token, err := p.AccessToken(context.TODO()); err == nil {
		if "ABCD" != token {
			t.Errorf("Unexpected return token: %s", token)
		}
	} else {
		t.Errorf("A non-nil error was returned from the fixed token provider")
	}

}
