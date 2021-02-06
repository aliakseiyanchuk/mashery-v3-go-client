package v3client

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

var serviceAllFieldsQuery = url.Values{
	"fields": {MasheryServiceFullFieldsStr},
}

func (c *HttpTransport) GetService(ctx context.Context, id string) (*MasheryService, error) {
	rv, err := c.getObject(ctx, FetchSpec{
		Resource:       fmt.Sprintf("/services/%s", id),
		Query:          serviceAllFieldsQuery,
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
func (c *HttpTransport) CreateService(ctx context.Context, service MasheryService) (*MasheryService, error) {
	rawResp, err := c.createObject(ctx, service, FetchSpec{
		Resource:   "/services",
		AppContext: "services",
		Query: url.Values{
			"fields": {MasheryServiceFullFieldsStr},
		},
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
func (c *HttpTransport) UpdateService(ctx context.Context, service MasheryService) (*MasheryService, error) {
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

// Delete a service.
func (c *HttpTransport) DeleteService(ctx context.Context, serviceId string) error {
	opContext := FetchSpec{
		Resource:   fmt.Sprintf("/services/%s", serviceId),
		AppContext: "service",
	}

	return c.deleteObject(ctx, opContext)
}

// List services that are filtered according to the condition that is V3-supported and containing the fields
// that the requester specifies
func (c *HttpTransport) ListServicesFiltered(ctx context.Context, params map[string]string, fields []string) ([]MasheryService, error) {
	return c.listServicesWithQuery(ctx, c.v3FilteringParams(params, fields))
}

func (c *HttpTransport) ListServices(ctx context.Context) ([]MasheryService, error) {
	return c.listServicesWithQuery(ctx, nil)
}

func (c *HttpTransport) listServicesWithQuery(ctx context.Context, qs url.Values) ([]MasheryService, error) {
	opCtx := FetchSpec{
		Pagination:     PerItem,
		Resource:       "/services",
		Query:          qs,
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

// Count the number of services that would match this criteria
func (c *HttpTransport) CountServices(ctx context.Context, params map[string]string) (int64, error) {
	opCtx := FetchSpec{
		Pagination: NotRequired,
		Resource:   "/services",
		Query: url.Values{
			"filter": {toV3FilterExpression(params)},
		},
		AppContext:     "all service count",
		ResponseParser: ParseMasheryServiceArray,
	}

	return c.count(ctx, opCtx)
}
