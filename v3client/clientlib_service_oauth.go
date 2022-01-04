package v3client

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
)

func GetServiceOAuthSecurityProfile(ctx context.Context, id string, c *transport.V3Transport) (*masherytypes.MasheryOAuth, error) {
	rv, err := c.GetObject(ctx, transport.FetchSpec{
		Resource:       fmt.Sprintf("/services/%s/securityProfile/oauth", id),
		Query:          nil,
		AppContext:     "service security profile oauth",
		ResponseParser: masherytypes.ParseMasheryServiceSecurityProfileOAuth,
	})

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(masherytypes.MasheryOAuth)
		return &retServ, nil
	}
}

// CreateServiceOAuthSecurityProfile Create a new service.
func CreateServiceOAuthSecurityProfile(ctx context.Context, id string, service masherytypes.MasheryOAuth, c *transport.V3Transport) (*masherytypes.MasheryOAuth, error) {
	rawResp, err := c.CreateObject(ctx, service, transport.FetchSpec{
		Resource:       fmt.Sprintf("/services/%s/securityProfile/oauth", id),
		AppContext:     "service security profile oauth",
		ResponseParser: masherytypes.ParseMasheryServiceSecurityProfileOAuth,
	})

	if err == nil {
		rv, _ := rawResp.(masherytypes.MasheryOAuth)
		return &rv, nil
	} else {
		return nil, err
	}
}

// UpdateServiceOAuthSecurityProfile Create a new service.
func UpdateServiceOAuthSecurityProfile(ctx context.Context, id string, service masherytypes.MasheryOAuth, c *transport.V3Transport) (*masherytypes.MasheryOAuth, error) {
	opContext := transport.FetchSpec{
		Resource:       fmt.Sprintf("/services/%s/securityProfile/oauth", id),
		AppContext:     "service security profile oauth",
		ResponseParser: masherytypes.ParseMasheryServiceSecurityProfileOAuth,
	}

	if d, err := c.UpdateObject(ctx, service, opContext); err == nil {
		rv, _ := d.(masherytypes.MasheryOAuth)
		return &rv, nil
	} else {
		return nil, err
	}
}

// DeleteServiceOAuthSecurityProfile Create a new service.
func DeleteServiceOAuthSecurityProfile(ctx context.Context, id string, c *transport.V3Transport) error {
	opContext := transport.FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/securityProfile/oauth", id),
		AppContext: "service security profile oauth",
	}

	return c.DeleteObject(ctx, opContext)
}
