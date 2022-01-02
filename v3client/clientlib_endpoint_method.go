package v3client

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

// ListEndpointMethods List methods associated with this endpoint, having only implicit fields returned.
func ListEndpointMethods(ctx context.Context, serviceId, endpointId string, c *HttpTransport) ([]MasheryMethod, error) {
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

// ListEndpointMethodsWithFullInfo List endpoints methods with their extended information.
func ListEndpointMethodsWithFullInfo(ctx context.Context, serviceId, endpointId string, c *HttpTransport) ([]MasheryMethod, error) {
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

// CreateEndpointMethod Create a new service.
func CreateEndpointMethod(ctx context.Context, serviceId, endpointId string, methodUpsert MasheryMethod, c *HttpTransport) (*MasheryMethod, error) {
	rawResp, err := c.createObject(ctx, methodUpsert, FetchSpec{
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

// UpdateEndpointMethod Update mashery endpoint method using the specified upsertable.
func UpdateEndpointMethod(ctx context.Context, serviceId, endpointId string, methUpsert MasheryMethod, c *HttpTransport) (*MasheryMethod, error) {
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

func GetEndpointMethod(ctx context.Context, serviceId, endpointId, methodId string, c *HttpTransport) (*MasheryMethod, error) {
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

func DeleteEndpointMethod(ctx context.Context, serviceId, endpointId, methodId string, c *HttpTransport) error {
	return c.deleteObject(ctx, FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods/%s", serviceId, endpointId, methodId),
		AppContext: "endpoint method",
	})
}

// CountEndpointsMethodsOf Count the number of services that would match this criteria
func CountEndpointsMethodsOf(ctx context.Context, serviceId, endpointId string, c *HttpTransport) (int64, error) {
	opCtx := FetchSpec{
		Pagination: NotRequired,
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods", serviceId, endpointId),
		AppContext: "endpoints methods",
	}

	return c.count(ctx, opCtx)
}
