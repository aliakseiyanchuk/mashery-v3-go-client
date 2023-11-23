package main

import (
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestServiceEndpointShowTest(t *testing.T) {
	obj := ObjectWithExists[masherytypes.ServiceEndpointIdentifier, masherytypes.Endpoint]{
		Identifier: masherytypes.ServiceEndpointIdentifier{},
		Object:     createBaselineServiceEndpoint(),
		Exists:     true,
	}

	str, code := executeTemplate(subCmdServiceEndpointShow.Template, obj)
	assert.Equal(t, 0, code)
	fmt.Println(str)
}

func createBaselineServiceEndpoint() masherytypes.Endpoint {
	adapter := "customAdapter"
	auth := "sysAuth"

	return masherytypes.Endpoint{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:        "endpoint-id",
			Name:      "endpoint-name",
			Created:   nil,
			Updated:   nil,
			Retrieved: time.Time{},
		},
		AllowMissingApiKey:          false,
		ApiKeyValueLocationKey:      "",
		ApiKeyValueLocations:        nil,
		ApiMethodDetectionKey:       "",
		ApiMethodDetectionLocations: []string{"a", "b", "c"},
		Cache: &masherytypes.Cache{
			CacheTTLOverride:               1,
			ContentCacheKeyHeaders:         []string{"ch-a", "ch-b", "ch-c"},
			ClientSurrogateControlEnabled:  true,
			IncludeApiKeyInContentCacheKey: true,
			RespondFromStaleCacheEnabled:   true,
			ResponseCacheControlEnabled:    true,
			VaryHeaderEnabled:              true,
		},
		ConnectionTimeoutForSystemDomainRequest:  5,
		ConnectionTimeoutForSystemDomainResponse: 60,
		CookiesDuringHttpRedirectsEnabled:        false,
		Cors: &masherytypes.Cors{
			AllDomainsEnabled:        true,
			SubDomainMatchingAllowed: true,
			MaxAge:                   60,
			CookiesAllowed:           true,
			AllowedDomains:           []string{"a", "b", "c"},
			AllowedHeaders:           []string{"a", "b", "c"},
			ExposedHeaders:           []string{"a", "b", "c"},
		},
		CustomRequestAuthenticationAdapter:         &adapter,
		DropApiKeyFromIncomingCall:                 true,
		ForceGzipOfBackendCall:                     true,
		GzipPassthroughSupportEnabled:              false,
		HeadersToExcludeFromIncomingCall:           []string{"a", "b", "c"},
		HighSecurity:                               false,
		HostPassthroughIncludedInBackendCallHeader: false,
		InboundSslRequired:                         true,
		InboundMutualSslRequired:                   false,
		JsonpCallbackParameter:                     "",
		JsonpCallbackParameterValue:                "",
		ScheduledMaintenanceEvent:                  nil,
		ForwardedHeaders:                           []string{"fh-a", "fh-b", "fh-c"},
		ReturnedHeaders:                            []string{"rh-a", "rh-b", "rh-c"},
		Methods:                                    nil,
		NumberOfHttpRedirectsToFollow:              3,
		OutboundRequestTargetPath:                  "/abc/",
		OutboundRequestTargetQueryParameters:       "/def/",
		OutboundTransportProtocol:                  "https",
		Processor: &masherytypes.Processor{
			PreProcessEnabled:  true,
			PostProcessEnabled: false,
			PostInputs: map[string]string{
				"post-a": "post-b",
				"post-c": "post-d",
			},
			PreInputs: map[string]string{
				"pre-a": "pre-b",
				"pre-c": "pre-d",
			},
			Adapter: "process-adapter",
		},
		PublicDomains: []masherytypes.Domain{
			{
				Address: "abc.def",
			},
			{
				Address: "abc1.def1",
			},
		},
		RequestAuthenticationType: "oauth",
		RequestPathAlias:          "/ff/",
		RequestProtocol:           "https",
		OAuthGrantTypes:           []string{"gt"},
		StringsToTrimFromApiKey:   "a",
		SupportedHttpMethods:      []string{"get", "post"},
		SystemDomainAuthentication: &masherytypes.SystemDomainAuthentication{
			Type:        "vv",
			Username:    &auth,
			Certificate: nil,
			Password:    nil,
		},
		SystemDomains: []masherytypes.Domain{
			{Address: "abc.def"},
		},
		TrafficManagerDomain:           "a.b.c.",
		UseSystemDomainCredentials:     false,
		SystemDomainCredentialKey:      nil,
		SystemDomainCredentialSecret:   nil,
		ErrorSet:                       nil,
		UserControlledErrorLocation:    "11",
		UserControlledErrorLocationKey: "22",
	}
}
