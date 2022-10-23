package v3client

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
)

func GetServiceOAuthSecurityProfile(ctx context.Context, id masherytypes.ServiceIdentifier, c *transport.V3Transport) (*masherytypes.MasheryOAuth, error) {
	rv, err := c.GetObject(ctx, transport.FetchSpec{
		Resource:       fmt.Sprintf("/services/%s/securityProfile/oauth", id.ServiceId),
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
func CreateServiceOAuthSecurityProfile(ctx context.Context, oauth masherytypes.MasheryOAuth, c *transport.V3Transport) (*masherytypes.MasheryOAuth, error) {
	rawResp, err := c.CreateObject(ctx, oauth, transport.FetchSpec{
		Resource:       fmt.Sprintf("/services/%s/securityProfile/oauth", oauth.ParentService.ServiceId),
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
func UpdateServiceOAuthSecurityProfile(ctx context.Context, oauth masherytypes.MasheryOAuth, c *transport.V3Transport) (*masherytypes.MasheryOAuth, error) {
	opContext := transport.FetchSpec{
		Resource:       fmt.Sprintf("/services/%s/securityProfile/oauth", oauth.ParentService.ServiceId),
		AppContext:     "service security profile oauth",
		ResponseParser: masherytypes.ParseMasheryServiceSecurityProfileOAuth,
	}

	if d, err := c.UpdateObject(ctx, oauth, opContext); err == nil {
		rv, _ := d.(masherytypes.MasheryOAuth)
		return &rv, nil
	} else {
		return nil, err
	}
}

// DeleteServiceOAuthSecurityProfile Create a new service.
func DeleteServiceOAuthSecurityProfile(ctx context.Context, id masherytypes.ServiceIdentifier, c *transport.V3Transport) error {
	opContext := transport.FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/securityProfile/oauth", id.ServiceId),
		AppContext: "service security profile oauth",
	}

	return c.DeleteObject(ctx, opContext)
}
