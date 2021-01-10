package v3client_test

import (
	"encoding/json"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"io/ioutil"
	"testing"
	"time"
)

const savedFileName = "./out/sampleSavedAccessToken.json"

func saveTestFile(inp *v3client.TimedAccessTokenResponse) bool {
	if data, err := json.Marshal(inp); err == nil {
		err = ioutil.WriteFile(savedFileName, data, 0644)
		return err == nil
	} else {
		return false
	}

}

func TestNewFileSystemTokenProvider(t *testing.T) {
	ref := v3client.TimedAccessTokenResponse{
		Obtained: time.Now(),
		AccessTokenResponse: v3client.AccessTokenResponse{
			TokenType:    "bearer",
			ApiKey:       "apiKey",
			AccessToken:  "accessToken",
			ExpiresIn:    3600,
			RefreshToken: "refreshToken",
			Scope:        "areaId",
		},
	}

	saved := saveTestFile(&ref)
	if !saved {
		t.Log("Test file could not be saved")
		t.FailNow()
	}

	if p, err := v3client.NewFileSystemTokenProviderFrom(savedFileName); err == nil {
		token, tokenInvalidError := p.AccessToken()
		if tokenInvalidError != nil {
			t.Errorf("The token must be valid")
		}
		if token != "accessToken" {
			t.Errorf("Unexpected access token value: %s", token)
		}
	} else {
		t.Errorf("File system provider produced an error: %s", err)
	}
}

func TestNewFileSystemTokenProviderWithExpiredToken(t *testing.T) {
	ref := v3client.TimedAccessTokenResponse{
		Obtained: time.Unix(time.Now().Unix()-int64(7200), 0),
		AccessTokenResponse: v3client.AccessTokenResponse{
			TokenType:    "bearer",
			ApiKey:       "apiKey",
			AccessToken:  "accessToken",
			ExpiresIn:    3600,
			RefreshToken: "refreshToken",
			Scope:        "areaId",
		},
	}

	saved := saveTestFile(&ref)
	if !saved {
		t.Log("Test file could not be saved")
		t.FailNow()
	}

	if p, err := v3client.NewFileSystemTokenProviderFrom(savedFileName); err == nil {
		_, tokenInvalidError := p.AccessToken()
		if tokenInvalidError == nil {
			t.Errorf("Token MUST be declared invalid")
		}
	} else {
		t.Errorf("File system provider produced an error: %s", err)
	}
}
