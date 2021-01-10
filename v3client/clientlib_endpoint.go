package v3client

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

func (c *Client) ListEndpoints(ctx context.Context, serviceId string) ([]AddressableV3Object, error) {
	spec := FetchSpec{
		Pagination:     PerPage,
		Resource:       fmt.Sprintf("/services/%s/endpoints", serviceId),
		Query:          nil,
		AppContext:     "endpoint of service",
		ResponseParser: ParseMasheryEndpointArray,
	}

	if d, err := c.fetchAll(ctx, spec); err != nil {
		return []AddressableV3Object{}, err
	} else {
		// Convert individual fetches into the array of elements
		var rv []MasheryEndpoint
		for _, raw := range d {
			ms, ok := raw.([]MasheryEndpoint)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return AddressableEndpoints(rv), nil
	}
}

// Create a new service.
func (c *Client) CreateEndpoint(ctx context.Context, serviceId string, endp MasheryEndpoint) (*MasheryEndpoint, error) {
	rawResp, err := c.createObject(ctx, endp, FetchSpec{
		Resource:       fmt.Sprintf("/services/%s/endpoints", serviceId),
		AppContext:     "endpoint",
		ResponseParser: ParseMasheryEndpoint,
	})

	if err == nil {
		rv, _ := rawResp.(MasheryEndpoint)
		return &rv, nil
	} else {
		return nil, err
	}
}

// Create a new service.
func (c *Client) UpdateEndpoint(ctx context.Context, serviceId string, endp MasheryEndpoint) (*MasheryEndpoint, error) {
	if endp.Id == "" {
		return nil, errors.New("illegal argument: endpoint Id must be set and not nil")
	}

	opContext := FetchSpec{
		Resource:       fmt.Sprintf("/services/%s/endpoints/%s", serviceId, endp.Id),
		AppContext:     "endpoint",
		ResponseParser: ParseMasheryEndpoint,
	}

	if d, err := c.updateObject(ctx, endp, opContext); err == nil {
		rv, _ := d.(MasheryEndpoint)
		return &rv, nil
	} else {
		return nil, err
	}
}

func (c *Client) GetEndpoint(ctx context.Context, serviceId string, endpointId string) (*MasheryEndpoint, error) {
	qs := url.Values{
		"fields": {
			"id", "allowMissingApiKey", "apiKeyValueLocationKey", "apiKeyValueLocations",
			"apiMethodDetectionKey", "apiMethodDetectionLocations", "cache", "connectionTimeoutForSystemDomainRequest",
			"connectionTimeoutForSystemDomainResponse", "cookiesDuringHttpRedirectsEnabled", "cors",
			"created", "customRequestAuthenticationAdapter", "dropApiKeyFromIncomingCall", "forceGzipOfBackendCall",
			"gzipPassthroughSupportEnabled", "headersToExcludeFromIncomingCall", "highSecurity",
			"hostPassthroughIncludedInBackendCallHeader", "inboundSslRequired", "jsonpCallbackParameter",
			"jsonpCallbackParameterValue", "scheduledMaintenanceEvent", "forwardedHeaders", "returnedHeaders",
			"methods", "name", "numberOfHttpRedirectsToFollow", "outboundRequestTargetPath", "outboundRequestTargetQueryParameters",
			"outboundTransportProtocol", "processor", "publicDomains", "requestAuthenticationType",
			"requestPathAlias", "requestProtocol", "oauthGrantTypes", "stringsToTrimFromApiKey", "supportedHttpMethods",
			"systemDomainAuthentication", "systemDomains", "trafficManagerDomain", "updated", "useSystemDomainCredentials",
			"systemDomainCredentialKey", "systemDomainCredentialSecret",
		},
	}

	fetchSpec := FetchSpec{
		Pagination:     NotRequired,
		Resource:       fmt.Sprintf("/services/%s/endpoints/%s", serviceId, endpointId),
		Query:          qs,
		AppContext:     "endpoint",
		ResponseParser: ParseMasheryEndpoint,
	}

	if raw, err := c.getObject(ctx, fetchSpec); err != nil {
		return nil, &WrappedError{
			Context: "get endpoint",
			Cause:   err,
		}
	} else {
		if rv, ok := raw.(MasheryEndpoint); ok {
			return &rv, nil
		} else {
			return nil, errors.New("invalid return type")
		}
	}
}
