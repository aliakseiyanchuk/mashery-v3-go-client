package mashery_v3_go_client_test

import (
	mashery_v3_go_client "github.com/aliakseiyanchuk/mashery-v3-go-client"
	"testing"
)

func TestFixedTokenProvider(t *testing.T) {
	p := mashery_v3_go_client.NewFixedTokenProvider("ABCD")
	if p == nil {
		t.Errorf("Fixed token provider must be returned")
		t.FailNow()
	}

	if token, err := p.AccessToken(); err == nil {
		if "ABCD" != token {
			t.Errorf("Unexpected return token: %s", token)
		}
	} else {
		t.Errorf("A non-nil error was returned from the fixed token provider")
	}

}
