package v3client

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"net/url"
)

// ListEndpointMethods List methods associated with this endpoint, having only implicit fields returned.
func ListEndpointMethods(ctx context.Context, serviceId, endpointId string, c *transport.V3Transport) ([]masherytypes.MasheryMethod, error) {
	spec := transport.FetchSpec{
		Pagination:     transport.PerPage,
		Resource:       fmt.Sprintf("/services/%s/endpoints/%s/methods", serviceId, endpointId),
		Query:          nil,
		AppContext:     "endpoint methods",
		ResponseParser: masherytypes.ParseMasheryMethodArray,
	}

	if d, err := c.FetchAll(ctx, spec); err != nil {
		return []masherytypes.MasheryMethod{}, err
	} else {
		// Convert individual fetches into the array of elements
		var rv []masherytypes.MasheryMethod
		for _, raw := range d {
			ms, ok := raw.([]masherytypes.MasheryMethod)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}

// ListEndpointMethodsWithFullInfo List endpoints methods with their extended information.
func ListEndpointMethodsWithFullInfo(ctx context.Context, serviceId, endpointId string, c *transport.V3Transport) ([]masherytypes.MasheryMethod, error) {
	spec := transport.FetchSpec{
		Pagination: transport.PerPage,
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods", serviceId, endpointId),
		Query: url.Values{
			"fields": {MasheryMethodsFieldsStr},
		},
		AppContext:     "endpoint methods",
		ResponseParser: masherytypes.ParseMasheryMethodArray,
	}

	if d, err := c.FetchAll(ctx, spec); err != nil {
		return []masherytypes.MasheryMethod{}, err
	} else {
		// Convert individual fetches into the array of elements
		var rv []masherytypes.MasheryMethod
		for _, raw := range d {
			ms, ok := raw.([]masherytypes.MasheryMethod)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}

// CreateEndpointMethod Create a new service.
func CreateEndpointMethod(ctx context.Context, serviceId, endpointId string, methodUpsert masherytypes.MasheryMethod, c *transport.V3Transport) (*masherytypes.MasheryMethod, error) {
	rawResp, err := c.CreateObject(ctx, methodUpsert, transport.FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods", serviceId, endpointId),
		AppContext: "endpoint method",
		Query: url.Values{
			"fields": {MasheryMethodsFieldsStr},
		},
		ResponseParser: masherytypes.ParseMasheryMethod,
	})

	if err == nil {
		rv, _ := rawResp.(masherytypes.MasheryMethod)
		return &rv, nil
	} else {
		return nil, err
	}
}

// UpdateEndpointMethod Update mashery endpoint method using the specified upsertable.
func UpdateEndpointMethod(ctx context.Context, serviceId, endpointId string, methUpsert masherytypes.MasheryMethod, c *transport.V3Transport) (*masherytypes.MasheryMethod, error) {
	if methUpsert.Id == "" {
		return nil, errors.New("illegal argument: endpoint Id must be set and not nil")
	}

	opContext := transport.FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods/%s", serviceId, endpointId, methUpsert.Id),
		AppContext: "endpoint method",
		Query: url.Values{
			"fields": {MasheryMethodsFieldsStr},
		},
		ResponseParser: masherytypes.ParseMasheryMethod,
	}

	if d, err := c.UpdateObject(ctx, methUpsert, opContext); err == nil {
		rv, _ := d.(masherytypes.MasheryMethod)
		return &rv, nil
	} else {
		return nil, err
	}
}

func GetEndpointMethod(ctx context.Context, serviceId, endpointId, methodId string, c *transport.V3Transport) (*masherytypes.MasheryMethod, error) {
	fetchSpec := transport.FetchSpec{
		Pagination: transport.NotRequired,
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods/%s", serviceId, endpointId, methodId),
		Query: url.Values{
			"fields": {MasheryMethodsFieldsStr},
		},
		AppContext:     "endpoint method",
		ResponseParser: masherytypes.ParseMasheryMethod,
	}

	if raw, err := c.GetObject(ctx, fetchSpec); err != nil {
		return nil, err
	} else {
		if rv, ok := raw.(masherytypes.MasheryMethod); ok {
			return &rv, nil
		} else {
			return nil, errors.New("invalid return type")
		}
	}
}

func DeleteEndpointMethod(ctx context.Context, serviceId, endpointId, methodId string, c *transport.V3Transport) error {
	return c.DeleteObject(ctx, transport.FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods/%s", serviceId, endpointId, methodId),
		AppContext: "endpoint method",
	})
}

// CountEndpointsMethodsOf Count the number of services that would match this criteria
func CountEndpointsMethodsOf(ctx context.Context, serviceId, endpointId string, c *transport.V3Transport) (int64, error) {
	opCtx := transport.FetchSpec{
		Pagination: transport.NotRequired,
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods", serviceId, endpointId),
		AppContext: "endpoints methods",
	}

	return c.Count(ctx, opCtx)
}
