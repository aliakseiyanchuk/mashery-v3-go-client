package v3client

import (
	"context"
	"fmt"
)

func (c *HttpTransport) GetServiceOAuthSecurityProfile(ctx context.Context, id string) (*MasheryOAuth, error) {
	rv, err := c.getObject(ctx, FetchSpec{
		Resource:       fmt.Sprintf("/services/%s/securityProfile/oauth", id),
		Query:          nil,
		AppContext:     "service security profile oauth",
		ResponseParser: ParseMasheryServiceSecurityProfileOAuth,
	})

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(MasheryOAuth)
		return &retServ, nil
	}
}

// Create a new service.
func (c *HttpTransport) CreateServiceOAuthSecurityProfile(ctx context.Context, id string, service MasheryOAuth) (*MasheryOAuth, error) {
	rawResp, err := c.createObject(ctx, service, FetchSpec{
		Resource:       fmt.Sprintf("/services/%s/securityProfile/oauth", id),
		AppContext:     "service security profile oauth",
		ResponseParser: ParseMasheryServiceSecurityProfileOAuth,
	})

	if err == nil {
		rv, _ := rawResp.(MasheryOAuth)
		return &rv, nil
	} else {
		return nil, err
	}
}

// Create a new service.
func (c *HttpTransport) UpdateServiceOAuthSecurityProfile(ctx context.Context, id string, service MasheryOAuth) (*MasheryOAuth, error) {
	opContext := FetchSpec{
		Resource:       fmt.Sprintf("/services/%s/securityProfile/oauth", id),
		AppContext:     "service security profile oauth",
		ResponseParser: ParseMasheryServiceSecurityProfileOAuth,
	}

	if d, err := c.updateObject(ctx, service, opContext); err == nil {
		rv, _ := d.(MasheryOAuth)
		return &rv, nil
	} else {
		return nil, err
	}
}

// Create a new service.
func (c *HttpTransport) DeleteServiceOAuthSecurityProfile(ctx context.Context, id string) error {
	opContext := FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/securityProfile/oauth", id),
		AppContext: "service security profile oauth",
	}

	return c.deleteObject(ctx, opContext)
}
