package v3client

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"net/url"
)

// ListEndpointMethodFilters List filters associated with this endpoint method, having only implicit fields returned.
func ListEndpointMethodFilters(ctx context.Context, serviceId, endpointId, methodId string, c *transport.V3Transport) ([]masherytypes.MasheryResponseFilter, error) {
	spec := transport.FetchSpec{
		Pagination:     transport.PerPage,
		Resource:       fmt.Sprintf("/services/%s/endpoints/%s/methods/%s/responseFilters", serviceId, endpointId, methodId),
		Query:          nil,
		AppContext:     "endpoint methods filters",
		ResponseParser: masherytypes.ParseMasheryResponseFilterArray,
	}

	if d, err := c.FetchAll(ctx, spec); err != nil {
		return []masherytypes.MasheryResponseFilter{}, err
	} else {
		var rv []masherytypes.MasheryResponseFilter
		for _, raw := range d {
			ms, ok := raw.([]masherytypes.MasheryResponseFilter)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}

// ListEndpointMethodFiltersWithFullInfo List endpoints methods filters with their extended information.
func ListEndpointMethodFiltersWithFullInfo(ctx context.Context, serviceId, endpointId, methodId string, c *transport.V3Transport) ([]masherytypes.MasheryResponseFilter, error) {
	spec := transport.FetchSpec{
		Pagination: transport.PerPage,
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods/%s/responseFilters", serviceId, endpointId, methodId),
		Query: url.Values{
			"fields": {MasheryResponseFilterFieldsStr},
		},
		AppContext:     "endpoint methods filters",
		ResponseParser: masherytypes.ParseMasheryResponseFilterArray,
	}

	if d, err := c.FetchAll(ctx, spec); err != nil {
		return []masherytypes.MasheryResponseFilter{}, err
	} else {
		// Convert individual fetches into the array of elements
		var rv []masherytypes.MasheryResponseFilter
		for _, raw := range d {
			ms, ok := raw.([]masherytypes.MasheryResponseFilter)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}

// CreateEndpointMethodFilter Create a new service.
func CreateEndpointMethodFilter(ctx context.Context, serviceId, endpointId, methodId string, filterUpsert masherytypes.MasheryResponseFilter, c *transport.V3Transport) (*masherytypes.MasheryResponseFilter, error) {
	rawResp, err := c.CreateObject(ctx, filterUpsert, transport.FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods/%s/responseFilters", serviceId, endpointId, methodId),
		AppContext: "endpoint method filters",
		Query: url.Values{
			"fields": {MasheryResponseFilterFieldsStr},
		},
		ResponseParser: masherytypes.ParseMasheryResponseFilter,
	})

	if err == nil {
		rv, _ := rawResp.(masherytypes.MasheryResponseFilter)
		return &rv, nil
	} else {
		return nil, err
	}
}

// UpdateEndpointMethodFilter Update mashery endpoint method using the specified upsertable.
func UpdateEndpointMethodFilter(ctx context.Context, serviceId, endpointId, methodId string, methUpsert masherytypes.MasheryResponseFilter, c *transport.V3Transport) (*masherytypes.MasheryResponseFilter, error) {
	if methUpsert.Id == "" {
		return nil, errors.New("illegal argument: response filter must be set and not nil")
	}

	opContext := transport.FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods/%s/responseFilters/%s", serviceId, endpointId, methodId, methUpsert.Id),
		AppContext: "endpoint method filters",
		Query: url.Values{
			"fields": {MasheryResponseFilterFieldsStr},
		},
		ResponseParser: masherytypes.ParseMasheryResponseFilter,
	}

	if d, err := c.UpdateObject(ctx, methUpsert, opContext); err == nil {
		rv, _ := d.(masherytypes.MasheryResponseFilter)
		return &rv, nil
	} else {
		return nil, err
	}
}

func GetEndpointMethodFilter(ctx context.Context, serviceId, endpointId, methodId, filterId string, c *transport.V3Transport) (*masherytypes.MasheryResponseFilter, error) {
	fetchSpec := transport.FetchSpec{
		Pagination: transport.NotRequired,
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods/%s/responseFilters/%s", serviceId, endpointId, methodId, filterId),
		Query: url.Values{
			"fields": {MasheryResponseFilterFieldsStr},
		},
		AppContext:     "endpoint method filters",
		ResponseParser: masherytypes.ParseMasheryResponseFilter,
	}

	if raw, err := c.GetObject(ctx, fetchSpec); err != nil {
		return nil, err
	} else {
		if rv, ok := raw.(masherytypes.MasheryResponseFilter); ok {
			return &rv, nil
		} else {
			return nil, errors.New("invalid return type")
		}
	}
}

func DeleteEndpointMethodFilter(ctx context.Context, serviceId, endpointId, methodId, filterId string, c *transport.V3Transport) error {
	return c.DeleteObject(ctx, transport.FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods/%s/responseFilters/%s", serviceId, endpointId, methodId, filterId),
		AppContext: "endpoint method filters",
	})
}

// CountEndpointsMethodsFiltersOf Count the number of services that would match this criteria
func CountEndpointsMethodsFiltersOf(ctx context.Context, serviceId, endpointId, methodId string, c *transport.V3Transport) (int64, error) {
	opCtx := transport.FetchSpec{
		Pagination: transport.NotRequired,
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods/%s/responseFilters", serviceId, endpointId, methodId),
		AppContext: "endpoint method filters",
	}

	return c.Count(ctx, opCtx)
}
