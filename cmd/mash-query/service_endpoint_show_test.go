package main

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"testing"
	"time"
)

func TestServiceEndpointShowTemplate(t *testing.T) {

}

func createBaselineEndpoint() masherytypes.Endpoint {
	return masherytypes.Endpoint{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:        "endpoint-id",
			Name:      "endpoint-name",
			Created:   nil,
			Updated:   nil,
			Retrieved: time.Time{},
		},
		AllowMissingApiKey:                         false,
		ApiKeyValueLocationKey:                     "",
		ApiKeyValueLocations:                       []string{"loc"},
		ApiMethodDetectionKey:                      "",
		ApiMethodDetectionLocations:                nil,
		Cache:                                      nil,
		ConnectionTimeoutForSystemDomainRequest:    0,
		ConnectionTimeoutForSystemDomainResponse:   0,
		CookiesDuringHttpRedirectsEnabled:          false,
		Cors:                                       nil,
		CustomRequestAuthenticationAdapter:         nil,
		DropApiKeyFromIncomingCall:                 false,
		ForceGzipOfBackendCall:                     false,
		GzipPassthroughSupportEnabled:              false,
		HeadersToExcludeFromIncomingCall:           nil,
		HighSecurity:                               false,
		HostPassthroughIncludedInBackendCallHeader: false,
		InboundSslRequired:                         false,
		InboundMutualSslRequired:                   false,
		JsonpCallbackParameter:                     "",
		JsonpCallbackParameterValue:                "",
		ScheduledMaintenanceEvent:                  nil,
		ForwardedHeaders:                           nil,
		ReturnedHeaders:                            nil,
		Methods:                                    nil,
		NumberOfHttpRedirectsToFollow:              0,
		OutboundRequestTargetPath:                  "",
		OutboundRequestTargetQueryParameters:       "",
		OutboundTransportProtocol:                  "",
		Processor:                                  nil,
		PublicDomains:                              nil,
		RequestAuthenticationType:                  "",
		RequestPathAlias:                           "",
		RequestProtocol:                            "",
		OAuthGrantTypes:                            nil,
		StringsToTrimFromApiKey:                    "",
		SupportedHttpMethods:                       nil,
		SystemDomainAuthentication:                 nil,
		SystemDomains:                              nil,
		TrafficManagerDomain:                       "",
		UseSystemDomainCredentials:                 false,
		SystemDomainCredentialKey:                  nil,
		SystemDomainCredentialSecret:               nil,
		ErrorSet:                                   nil,
		UserControlledErrorLocation:                "",
		UserControlledErrorLocationKey:             "",
	}
}
