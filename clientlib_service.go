package mashery_v3_go_client

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

func (c *Client) GetService(ctx context.Context, id string) (*MasheryService, error) {
	qs := url.Values{
		"fields": {
			"id", "name", "created", "updated", "endpoints", "editorHandle",
			"revisionNumber", "robotsPolicy", "crossdomainPolicy",
			"description", "errorSets", "qpsLimitOverall",
			"rfc3986Encode", "securityProfile", "version",
		},
	}

	resource := fmt.Sprintf("/services/%s?fields=%s", id, qs.Encode())

	rv, err := c.getObject(ctx, FetchSpec{
		Resource:       resource,
		Query:          qs,
		AppContext:     "service",
		ResponseParser: ParseMasheryService,
	})

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(MasheryService)
		return &retServ, nil
	}
}

// Create a new service.
func (c *Client) CreateService(ctx context.Context, service MasheryService) (*MasheryService, error) {
	rawResp, err := c.createObject(ctx, service, FetchSpec{
		Resource:       "/services",
		AppContext:     "services",
		ResponseParser: ParseMasheryService,
	})

	if err == nil {
		rv, _ := rawResp.(MasheryService)
		return &rv, nil
	} else {
		return nil, err
	}
}

// Create a new service.
func (c *Client) UpdateService(ctx context.Context, service MasheryService) (*MasheryService, error) {
	if service.Id == "" {
		return nil, errors.New("illegal argument: service Id must be set and not nil")
	}

	opContext := FetchSpec{
		Resource:       fmt.Sprintf("/services/%s", service.Id),
		AppContext:     "service",
		ResponseParser: ParseMasheryService,
	}

	if d, err := c.updateObject(ctx, service, opContext); err == nil {
		rv, _ := d.(MasheryService)
		return &rv, nil
	} else {
		return nil, err
	}
}

func (c *Client) ListServices(ctx context.Context) ([]MasheryService, error) {
	opCtx := FetchSpec{
		Pagination:     PerItem,
		Resource:       "/services",
		Query:          nil,
		AppContext:     "all service",
		ResponseParser: ParseMasheryServiceArray,
	}

	if d, err := c.fetchAll(ctx, opCtx); err != nil {
		return []MasheryService{}, nil
	} else {
		// Convert individual fetches into the array of elements
		var rv []MasheryService
		for _, raw := range d {
			ms, ok := raw.([]MasheryService)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}
