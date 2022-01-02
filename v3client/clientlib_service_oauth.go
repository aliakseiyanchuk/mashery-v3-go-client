package v3client

import (
	"context"
	"fmt"
)

func GetServiceOAuthSecurityProfile(ctx context.Context, id string, c *HttpTransport) (*MasheryOAuth, error) {
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

// CreateServiceOAuthSecurityProfile Create a new service.
func CreateServiceOAuthSecurityProfile(ctx context.Context, id string, service MasheryOAuth, c *HttpTransport) (*MasheryOAuth, error) {
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

// UpdateServiceOAuthSecurityProfile Create a new service.
func UpdateServiceOAuthSecurityProfile(ctx context.Context, id string, service MasheryOAuth, c *HttpTransport) (*MasheryOAuth, error) {
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

// DeleteServiceOAuthSecurityProfile Create a new service.
func DeleteServiceOAuthSecurityProfile(ctx context.Context, id string, c *HttpTransport) error {
	opContext := FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/securityProfile/oauth", id),
		AppContext: "service security profile oauth",
	}

	return c.deleteObject(ctx, opContext)
}
