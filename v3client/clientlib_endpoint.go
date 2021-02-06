package v3client

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

func (c *HttpTransport) ListEndpoints(ctx context.Context, serviceId string) ([]AddressableV3Object, error) {
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

// List endpoints with their extended information.
func (c *HttpTransport) ListEndpointsWithFullInfo(ctx context.Context, serviceId string) ([]MasheryEndpoint, error) {
	spec := FetchSpec{
		Pagination: PerPage,
		Resource:   fmt.Sprintf("/services/%s/endpoints", serviceId),
		Query: url.Values{
			"fields": {MasheryEndpointFieldsStr},
		},
		AppContext:     "endpoint of service",
		ResponseParser: ParseMasheryEndpointArray,
	}

	if d, err := c.fetchAll(ctx, spec); err != nil {
		return []MasheryEndpoint{}, err
	} else {
		// Convert individual fetches into the array of elements
		var rv []MasheryEndpoint
		for _, raw := range d {
			ms, ok := raw.([]MasheryEndpoint)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}

// Create a new service.
func (c *HttpTransport) CreateEndpoint(ctx context.Context, serviceId string, endp MasheryEndpoint) (*MasheryEndpoint, error) {
	rawResp, err := c.createObject(ctx, endp, FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/endpoints", serviceId),
		AppContext: "endpoint",
		Query: url.Values{
			"fields": {MasheryEndpointFieldsStr},
		},
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
func (c *HttpTransport) UpdateEndpoint(ctx context.Context, serviceId string, endp MasheryEndpoint) (*MasheryEndpoint, error) {
	if endp.Id == "" {
		return nil, errors.New("illegal argument: endpoint Id must be set and not nil")
	}

	opContext := FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s", serviceId, endp.Id),
		AppContext: "endpoint",
		Query: url.Values{
			"fields": {MasheryEndpointFieldsStr},
		},
		ResponseParser: ParseMasheryEndpoint,
	}

	if d, err := c.updateObject(ctx, endp, opContext); err == nil {
		rv, _ := d.(MasheryEndpoint)
		return &rv, nil
	} else {
		return nil, err
	}
}

func (c *HttpTransport) GetEndpoint(ctx context.Context, serviceId string, endpointId string) (*MasheryEndpoint, error) {
	fetchSpec := FetchSpec{
		Pagination: NotRequired,
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s", serviceId, endpointId),
		Query: url.Values{
			"fields": {MasheryEndpointFieldsStr},
		},
		AppContext:     "endpoint",
		ResponseParser: ParseMasheryEndpoint,
	}

	if raw, err := c.getObject(ctx, fetchSpec); err != nil {
		return nil, err
	} else {
		if rv, ok := raw.(MasheryEndpoint); ok {
			return &rv, nil
		} else {
			return nil, errors.New("invalid return type")
		}
	}
}

func (c *HttpTransport) DeleteEndpoint(ctx context.Context, serviceId, endpointId string) error {
	return c.deleteObject(ctx, FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s", serviceId, endpointId),
		AppContext: "endpoint",
	})
}

// Count the number of services that would match this criteria
func (c *HttpTransport) CountEndpointsOf(ctx context.Context, serviceId string) (int64, error) {
	opCtx := FetchSpec{
		Pagination: NotRequired,
		Resource:   fmt.Sprintf("/services/%s/endpoints", serviceId),
		AppContext: "service endpoints",
	}

	return c.count(ctx, opCtx)
}
