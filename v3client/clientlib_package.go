package v3client

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

func (c *Client) GetPackage(ctx context.Context, id string) (*MasheryPackage, error) {
	qs := url.Values{
		"fields": {
			"id", "name", "created", "updated", "description", "notifyDeveloperPeriod",
			"notifyDeveloperNearQuota", "notifyDeveloperOverQuota", "notifyDeveloperOverThrottle", "notifyAdminNearQuota",
			"notifyAdminOverQuota", "notifyAdminOverThrottle", "notifyAdminEmails", "nearQuotaThreshold", "eav", "keyAdapter", "keyLength",
			"sharedSecretLength", "plans",
		},
	}

	rv, err := c.getObject(ctx, FetchSpec{
		Resource:       fmt.Sprintf("/packages/%s", id),
		Query:          qs,
		AppContext:     "service",
		ResponseParser: ParseMasheryPackage,
	})

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(MasheryPackage)
		return &retServ, nil
	}
}

// Create a new service.
func (c *Client) CreatePackage(ctx context.Context, service MasheryService) (*MasheryPackage, error) {
	rawResp, err := c.createObject(ctx, service, FetchSpec{
		Resource:       "/packages",
		AppContext:     "package",
		ResponseParser: ParseMasheryPackage,
	})

	if err == nil {
		rv, _ := rawResp.(MasheryPackage)
		return &rv, nil
	} else {
		return nil, err
	}
}

// Create a new service.
func (c *Client) UpdatePackage(ctx context.Context, service MasheryPackage) (*MasheryPackage, error) {
	if service.Id == "" {
		return nil, errors.New("illegal argument: package Id must be set and not nil")
	}

	opContext := FetchSpec{
		Resource:       fmt.Sprintf("/packages/%s", service.Id),
		AppContext:     "package",
		ResponseParser: ParseMasheryService,
	}

	if d, err := c.updateObject(ctx, service, opContext); err == nil {
		rv, _ := d.(MasheryPackage)
		return &rv, nil
	} else {
		return nil, err
	}
}

func (c *Client) ListPackages(ctx context.Context) ([]MasheryPackage, error) {
	opCtx := FetchSpec{
		Pagination:     PerItem,
		Resource:       "/packages",
		Query:          nil,
		AppContext:     "all service",
		ResponseParser: ParseMasheryPackageArray,
	}

	if d, err := c.fetchAll(ctx, opCtx); err != nil {
		return []MasheryPackage{}, err
	} else {
		// Convert individual fetches into the array of elements
		var rv []MasheryPackage
		for _, raw := range d {
			ms, ok := raw.([]MasheryPackage)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}
