package v3client

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"net/url"
)

func ListEndpoints(ctx context.Context, serviceId string, c *transport.V3Transport) ([]masherytypes.AddressableV3Object, error) {
	spec := transport.FetchSpec{
		Pagination:     transport.PerPage,
		Resource:       fmt.Sprintf("/services/%s/endpoints", serviceId),
		Query:          nil,
		AppContext:     "endpoint of service",
		ResponseParser: masherytypes.ParseMasheryEndpointArray,
	}

	if d, err := c.FetchAll(ctx, spec); err != nil {
		return []masherytypes.AddressableV3Object{}, err
	} else {
		// Convert individual fetches into the array of elements
		var rv []masherytypes.MasheryEndpoint
		for _, raw := range d {
			ms, ok := raw.([]masherytypes.MasheryEndpoint)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return masherytypes.AddressableEndpoints(rv), nil
	}
}

// ListEndpointsWithFullInfo List endpoints with their extended information.
func ListEndpointsWithFullInfo(ctx context.Context, serviceId string, c *transport.V3Transport) ([]masherytypes.MasheryEndpoint, error) {
	spec := transport.FetchSpec{
		Pagination: transport.PerPage,
		Resource:   fmt.Sprintf("/services/%s/endpoints", serviceId),
		Query: url.Values{
			"fields": {MasheryEndpointFieldsStr},
		},
		AppContext:     "endpoint of service",
		ResponseParser: masherytypes.ParseMasheryEndpointArray,
	}

	if d, err := c.FetchAll(ctx, spec); err != nil {
		return []masherytypes.MasheryEndpoint{}, err
	} else {
		// Convert individual fetches into the array of elements
		var rv []masherytypes.MasheryEndpoint
		for _, raw := range d {
			ms, ok := raw.([]masherytypes.MasheryEndpoint)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}

// CreateEndpoint Create a new endpoint of the service.
func CreateEndpoint(ctx context.Context, serviceId string, endp masherytypes.MasheryEndpoint, c *transport.V3Transport) (*masherytypes.MasheryEndpoint, error) {
	rawResp, err := c.CreateObject(ctx, endp, transport.FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/endpoints", serviceId),
		AppContext: "endpoint",
		Query: url.Values{
			"fields": {MasheryEndpointFieldsStr},
		},
		ResponseParser: masherytypes.ParseMasheryEndpoint,
	})

	if err == nil {
		rv, _ := rawResp.(masherytypes.MasheryEndpoint)
		return &rv, nil
	} else {
		return nil, err
	}
}

// UpdateEndpoint updates an endpoint
func UpdateEndpoint(ctx context.Context, serviceId string, endp masherytypes.MasheryEndpoint, c *transport.V3Transport) (*masherytypes.MasheryEndpoint, error) {
	if endp.Id == "" {
		return nil, errors.New("illegal argument: endpoint Id must be set and not nil")
	}

	opContext := transport.FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s", serviceId, endp.Id),
		AppContext: "endpoint",
		Query: url.Values{
			"fields": {MasheryEndpointFieldsStr},
		},
		ResponseParser: masherytypes.ParseMasheryEndpoint,
	}

	if d, err := c.UpdateObject(ctx, endp, opContext); err == nil {
		rv, _ := d.(masherytypes.MasheryEndpoint)
		return &rv, nil
	} else {
		return nil, err
	}
}

func GetEndpoint(ctx context.Context, serviceId string, endpointId string, c *transport.V3Transport) (*masherytypes.MasheryEndpoint, error) {
	fetchSpec := transport.FetchSpec{
		Pagination: transport.NotRequired,
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s", serviceId, endpointId),
		Query: url.Values{
			"fields": {MasheryEndpointFieldsStr},
		},
		AppContext:     "endpoint",
		ResponseParser: masherytypes.ParseMasheryEndpoint,
	}

	if raw, err := c.GetObject(ctx, fetchSpec); err != nil {
		return nil, err
	} else {
		if rv, ok := raw.(masherytypes.MasheryEndpoint); ok {
			return &rv, nil
		} else {
			return nil, errors.New("invalid return type")
		}
	}
}

func DeleteEndpoint(ctx context.Context, serviceId, endpointId string, c *transport.V3Transport) error {
	return c.DeleteObject(ctx, transport.FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s", serviceId, endpointId),
		AppContext: "endpoint",
	})
}

// CountEndpointsOf Count the number of services that would match this criteria
func CountEndpointsOf(ctx context.Context, serviceId string, c *transport.V3Transport) (int64, error) {
	opCtx := transport.FetchSpec{
		Pagination: transport.NotRequired,
		Resource:   fmt.Sprintf("/services/%s/endpoints", serviceId),
		AppContext: "service endpoints",
	}

	return c.Count(ctx, opCtx)
}
