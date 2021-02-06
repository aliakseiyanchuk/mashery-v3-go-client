package v3client

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

// List methods associated with this endpoint, having only implicit fields returned.
func (c *HttpTransport) ListEndpointMethods(ctx context.Context, serviceId, endpointId string) ([]MasheryMethod, error) {
	spec := FetchSpec{
		Pagination:     PerPage,
		Resource:       fmt.Sprintf("/services/%s/endpoints/%s/methods", serviceId, endpointId),
		Query:          nil,
		AppContext:     "endpoint methods",
		ResponseParser: ParseMasheryMethodArray,
	}

	if d, err := c.fetchAll(ctx, spec); err != nil {
		return []MasheryMethod{}, err
	} else {
		// Convert individual fetches into the array of elements
		var rv []MasheryMethod
		for _, raw := range d {
			ms, ok := raw.([]MasheryMethod)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}

// List endpoints methods with their extended information.
func (c *HttpTransport) ListEndpointMethodsWithFullInfo(ctx context.Context, serviceId, endpointId string) ([]MasheryMethod, error) {
	spec := FetchSpec{
		Pagination: PerPage,
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods", serviceId, endpointId),
		Query: url.Values{
			"fields": {MasheryMethodsFieldsStr},
		},
		AppContext:     "endpoint methods",
		ResponseParser: ParseMasheryMethodArray,
	}

	if d, err := c.fetchAll(ctx, spec); err != nil {
		return []MasheryMethod{}, err
	} else {
		// Convert individual fetches into the array of elements
		var rv []MasheryMethod
		for _, raw := range d {
			ms, ok := raw.([]MasheryMethod)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}

// Create a new service.
func (c *HttpTransport) CreateEndpointMethod(ctx context.Context, serviceId, endpointId string, methoUpsert MasheryMethod) (*MasheryMethod, error) {
	rawResp, err := c.createObject(ctx, methoUpsert, FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods", serviceId, endpointId),
		AppContext: "endpoint method",
		Query: url.Values{
			"fields": {MasheryMethodsFieldsStr},
		},
		ResponseParser: ParseMasheryMethod,
	})

	if err == nil {
		rv, _ := rawResp.(MasheryMethod)
		return &rv, nil
	} else {
		return nil, err
	}
}

// Update mashery endpoint method using the specified upsertable.
func (c *HttpTransport) UpdateEndpointMethod(ctx context.Context, serviceId, endpointId string, methUpsert MasheryMethod) (*MasheryMethod, error) {
	if methUpsert.Id == "" {
		return nil, errors.New("illegal argument: endpoint Id must be set and not nil")
	}

	opContext := FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods/%s", serviceId, endpointId, methUpsert.Id),
		AppContext: "endpoint method",
		Query: url.Values{
			"fields": {MasheryMethodsFieldsStr},
		},
		ResponseParser: ParseMasheryMethod,
	}

	if d, err := c.updateObject(ctx, methUpsert, opContext); err == nil {
		rv, _ := d.(MasheryMethod)
		return &rv, nil
	} else {
		return nil, err
	}
}

func (c *HttpTransport) GetEndpointMethod(ctx context.Context, serviceId, endpointId, methodId string) (*MasheryMethod, error) {
	fetchSpec := FetchSpec{
		Pagination: NotRequired,
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods/%s", serviceId, endpointId, methodId),
		Query: url.Values{
			"fields": {MasheryMethodsFieldsStr},
		},
		AppContext:     "endpoint method",
		ResponseParser: ParseMasheryMethod,
	}

	if raw, err := c.getObject(ctx, fetchSpec); err != nil {
		return nil, err
	} else {
		if rv, ok := raw.(MasheryMethod); ok {
			return &rv, nil
		} else {
			return nil, errors.New("invalid return type")
		}
	}
}

func (c *HttpTransport) DeleteEndpointMethod(ctx context.Context, serviceId, endpointId, methodId string) error {
	return c.deleteObject(ctx, FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods/%s", serviceId, endpointId, methodId),
		AppContext: "endpoint method",
	})
}

// Count the number of services that would match this criteria
func (c *HttpTransport) CountEndpointsMethodsOf(ctx context.Context, serviceId, endpointId string) (int64, error) {
	opCtx := FetchSpec{
		Pagination: NotRequired,
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods", serviceId, endpointId),
		AppContext: "endpoints methods",
	}

	return c.count(ctx, opCtx)
}
