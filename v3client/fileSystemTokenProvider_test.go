package v3client_test

import (
	"context"
	"encoding/json"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"os"
	"testing"
	"time"
)

const savedFileName = "../out/sampleSavedAccessToken.json"

func saveTestFile(inp *masherytypes.TimedAccessTokenResponse) bool {
	if data, err := json.Marshal(inp); err == nil {
		err = os.WriteFile(savedFileName, data, 0644)
		return err == nil
	} else {
		return false
	}

}

func TestNewFileSystemTokenProvider(t *testing.T) {
	ref := masherytypes.TimedAccessTokenResponse{
		Obtained: time.Now(),
		AccessTokenResponse: masherytypes.AccessTokenResponse{
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

	p := v3client.NewFileSystemTokenProviderFrom(savedFileName)
	token, tokenInvalidError := p.AccessToken(context.TODO())
	if tokenInvalidError != nil {
		t.Errorf("The token must be valid")
	}
	if token != "accessToken" {
		t.Errorf("Unexpected access token value: %s", token)
	}
}

func TestNewFileSystemTokenProviderWithExpiredToken(t *testing.T) {
	ref := masherytypes.TimedAccessTokenResponse{
		Obtained: time.Unix(time.Now().Unix()-int64(7200), 0),
		AccessTokenResponse: masherytypes.AccessTokenResponse{
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

	p := v3client.NewFileSystemTokenProviderFrom(savedFileName)
	_, tokenInvalidError := p.AccessToken(context.TODO())
	if tokenInvalidError == nil {
		t.Errorf("Token MUST be declared invalid")
	}

}
