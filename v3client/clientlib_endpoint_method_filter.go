package v3client

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

// List filters associated with this endpoint method, having only implicit fields returned.
func (c *HttpTransport) ListEndpointMethodFilters(ctx context.Context, serviceId, endpointId, methodId string) ([]MasheryResponseFilter, error) {
	spec := FetchSpec{
		Pagination:     PerPage,
		Resource:       fmt.Sprintf("/services/%s/endpoints/%s/methods/%s/responseFilters", serviceId, endpointId, methodId),
		Query:          nil,
		AppContext:     "endpoint methods filters",
		ResponseParser: ParseMasheryResponseFilterArray,
	}

	if d, err := c.fetchAll(ctx, spec); err != nil {
		return []MasheryResponseFilter{}, err
	} else {
		var rv []MasheryResponseFilter
		for _, raw := range d {
			ms, ok := raw.([]MasheryResponseFilter)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}

// List endpoints methods filters with their extended information.
func (c *HttpTransport) ListEndpointMethodFiltersWithFullInfo(ctx context.Context, serviceId, endpointId, methodId string) ([]MasheryResponseFilter, error) {
	spec := FetchSpec{
		Pagination: PerPage,
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods/%s/responseFilters", serviceId, endpointId, methodId),
		Query: url.Values{
			"fields": {MasheryResponseFilterFieldsStr},
		},
		AppContext:     "endpoint methods filters",
		ResponseParser: ParseMasheryResponseFilterArray,
	}

	if d, err := c.fetchAll(ctx, spec); err != nil {
		return []MasheryResponseFilter{}, err
	} else {
		// Convert individual fetches into the array of elements
		var rv []MasheryResponseFilter
		for _, raw := range d {
			ms, ok := raw.([]MasheryResponseFilter)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}

// Create a new service.
func (c *HttpTransport) CreateEndpointMethodFilter(ctx context.Context, serviceId, endpointId, methodId string, filterUpsert MasheryResponseFilter) (*MasheryResponseFilter, error) {
	rawResp, err := c.createObject(ctx, filterUpsert, FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods/%s/responseFilters", serviceId, endpointId, methodId),
		AppContext: "endpoint method filters",
		Query: url.Values{
			"fields": {MasheryResponseFilterFieldsStr},
		},
		ResponseParser: ParseMasheryResponseFilter,
	})

	if err == nil {
		rv, _ := rawResp.(MasheryResponseFilter)
		return &rv, nil
	} else {
		return nil, err
	}
}

// Update mashery endpoint method using the specified upsertable.
func (c *HttpTransport) UpdateEndpointMethodFilter(ctx context.Context, serviceId, endpointId, methodId string, methUpsert MasheryResponseFilter) (*MasheryResponseFilter, error) {
	if methUpsert.Id == "" {
		return nil, errors.New("illegal argument: response filter must be set and not nil")
	}

	opContext := FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods/%s/responseFilters/%s", serviceId, endpointId, methodId, methUpsert.Id),
		AppContext: "endpoint method filters",
		Query: url.Values{
			"fields": {MasheryResponseFilterFieldsStr},
		},
		ResponseParser: ParseMasheryResponseFilter,
	}

	if d, err := c.updateObject(ctx, methUpsert, opContext); err == nil {
		rv, _ := d.(MasheryResponseFilter)
		return &rv, nil
	} else {
		return nil, err
	}
}

func (c *HttpTransport) GetEndpointMethodFilter(ctx context.Context, serviceId, endpointId, methodId, filterId string) (*MasheryResponseFilter, error) {
	fetchSpec := FetchSpec{
		Pagination: NotRequired,
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods/%s/responseFilters/%s", serviceId, endpointId, methodId, filterId),
		Query: url.Values{
			"fields": {MasheryResponseFilterFieldsStr},
		},
		AppContext:     "endpoint method filters",
		ResponseParser: ParseMasheryResponseFilter,
	}

	if raw, err := c.getObject(ctx, fetchSpec); err != nil {
		return nil, err
	} else {
		if rv, ok := raw.(MasheryResponseFilter); ok {
			return &rv, nil
		} else {
			return nil, errors.New("invalid return type")
		}
	}
}

func (c *HttpTransport) DeleteEndpointMethodFilter(ctx context.Context, serviceId, endpointId, methodId, filterId string) error {
	return c.deleteObject(ctx, FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods/%s/responseFilters/%s", serviceId, endpointId, methodId, filterId),
		AppContext: "endpoint method filters",
	})
}

// Count the number of services that would match this criteria
func (c *HttpTransport) CountEndpointsMethodsFiltersOf(ctx context.Context, serviceId, endpointId, methodId string) (int64, error) {
	opCtx := FetchSpec{
		Pagination: NotRequired,
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods/%s/responseFilters", serviceId, endpointId),
		AppContext: "endpoint method filters",
	}

	return c.count(ctx, opCtx)
}
