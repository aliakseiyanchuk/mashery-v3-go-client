package main

import (
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestServiceShowCommandWillBeFound(t *testing.T) {
	f := locateSubCommandExecutor([]string{"service", "show", "-service-id", "abc"})
	assert.NotNil(t, f)
}

func TestServiceShowTemplate(t *testing.T) {
	obj := createBaselineServiceObject()

	str, code := executeTemplate(subCmdShowService.Template, obj)
	assert.Equal(t, 0, code)
	fmt.Println(str)
}

func TestServiceShowTemplateWithoutForwardedOAuthHeaders(t *testing.T) {
	obj := createBaselineServiceObject()
	obj.Object.SecurityProfile.OAuth.ForwardedHeaders = nil

	str, code := executeTemplate(subCmdShowService.Template, obj)
	assert.Equal(t, 0, code)
	fmt.Println(str)
}

func createBaselineServiceObject() ObjectWithExists[masherytypes.ServiceIdentifier, masherytypes.Service] {
	return ObjectWithExists[masherytypes.ServiceIdentifier, masherytypes.Service]{
		Identifier: masherytypes.ServiceIdentifier{ServiceId: "srvId"},
		Exists:     true,
		Object: masherytypes.Service{
			AddressableV3Object: masherytypes.AddressableV3Object{
				Id:        "srvId",
				Name:      "srvName",
				Created:   nil,
				Updated:   nil,
				Retrieved: time.Time{},
			},
			Cache:             nil,
			Endpoints:         nil,
			EditorHandle:      "abc	",
			RevisionNumber:    43,
			RobotsPolicy:      "",
			CrossdomainPolicy: "",
			Description:       "Description",
			ErrorSets:         nil,
			QpsLimitOverall:   nil,
			RFC3986Encode:     false,
			SecurityProfile: &masherytypes.MasherySecurityProfile{
				OAuth: &masherytypes.MasheryOAuth{
					AccessTokenTtlEnabled:       true,
					AccessTokenTtl:              0,
					AccessTokenType:             "",
					AllowMultipleToken:          false,
					AuthorizationCodeTtl:        0,
					ForwardedHeaders:            []string{"a", "b", "c"},
					MasheryTokenApiEnabled:      false,
					RefreshTokenEnabled:         false,
					EnableRefreshTokenTtl:       false,
					TokenBasedRateLimitsEnabled: false,
					ForceOauthRedirectUrl:       false,
					ForceSslRedirectUrlEnabled:  false,
					GrantTypes:                  nil,
					MACAlgorithm:                "",
					QPSLimitCeiling:             0,
					RateLimitCeiling:            0,
					RefreshTokenTtl:             0,
					SecureTokensEnabled:         false,
				},
			},
			Version:      "Version",
			Roles:        nil,
			Organization: nil,
		},
	}
}
