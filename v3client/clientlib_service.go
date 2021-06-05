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

func GetService(ctx context.Context, id string, c *HttpTransport) (*MasheryService, error) {
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

// CreateService Create a new service.
func CreateService(ctx context.Context, service MasheryService, c *HttpTransport) (*MasheryService, error) {
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

// UpdateService updates a Mashery service
func UpdateService(ctx context.Context, service MasheryService, c *HttpTransport) (*MasheryService, error) {
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
func DeleteService(ctx context.Context, serviceId string, c *HttpTransport) error {
	opContext := FetchSpec{
		Resource:   fmt.Sprintf("/services/%s", serviceId),
		AppContext: "service",
	}

	return c.deleteObject(ctx, opContext)
}

// ListServicesFiltered List services that are filtered according to the condition that is V3-supported and containing the fields
// that the requester specifies
func ListServicesFiltered(ctx context.Context, params map[string]string, fields []string, c *HttpTransport) ([]MasheryService, error) {
	return listServicesWithQuery(ctx, c.v3FilteringParams(params, fields), c)
}

func ListServices(ctx context.Context, c *HttpTransport) ([]MasheryService, error) {
	return listServicesWithQuery(ctx, nil, c)
}

func listServicesWithQuery(ctx context.Context, qs url.Values, c *HttpTransport) ([]MasheryService, error) {
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
func CountServices(ctx context.Context, params map[string]string, c *HttpTransport) (int64, error) {
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
