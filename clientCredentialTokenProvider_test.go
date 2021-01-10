package mashery_v3_go_client_test

import (
	mashery_v3_go_client "github.com/aliakseiyanchuk/mashery-v3-go-client"
	"testing"
)

func TestMasheryV3Credentials_FullySpecified(t *testing.T) {
	creds := mashery_v3_go_client.MasheryV3Credentials{}
	if creds.FullySpecified() {
		t.Errorf("Credentials cannot be reported as fully specified where no field was set")
	}

	creds = mashery_v3_go_client.MasheryV3Credentials{
		AreaId:   "A",
		ApiKey:   "K",
		Secret:   "S",
		Username: "U",
		Password: "P",
	}
	if !creds.FullySpecified() {
		t.Errorf("Credentials must be reported as fully specified when all fields are set")
	}
}

func TestMasheryV3Credentials_Inherit(t *testing.T) {
	creds := mashery_v3_go_client.MasheryV3Credentials{}
	inh := mashery_v3_go_client.MasheryV3Credentials{
		AreaId:   "aid",
		ApiKey:   "key",
		Secret:   "sec",
		Username: "uname",
		Password: "pwd",
	}

	creds.Inherit(&inh)

	if creds.AreaId != "aid" {
		t.Errorf("Area Id is not inherited")
	}

	if creds.ApiKey != "key" {
		t.Errorf("Key is not inherited")
	}

	if creds.Secret != "sec" {
		t.Errorf("Secret is not inherited")
	}

	if creds.Username != "uname" {
		t.Errorf("Username is not inherited")
	}

	if creds.Password != "pwd" {
		t.Errorf("Password is not inherited")
	}

}
