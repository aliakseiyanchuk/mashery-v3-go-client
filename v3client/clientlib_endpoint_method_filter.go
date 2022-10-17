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
func ListEndpointMethodFilters(ctx context.Context, ident masherytypes.ServiceEndpointMethodIdentifier, c *transport.V3Transport) ([]masherytypes.ServiceEndpointMethodFilter, error) {
	spec := transport.FetchSpec{
		Pagination:     transport.PerPage,
		Resource:       fmt.Sprintf("/services/%s/endpoints/%s/methods/%s/responseFilters", ident.ServiceId, ident.EndpointId, ident.MethodId),
		Query:          nil,
		AppContext:     "endpoint methods filters",
		ResponseParser: masherytypes.ParseServiceEndpointMethodFilterArray,
	}

	if d, err := c.FetchAll(ctx, spec); err != nil {
		return []masherytypes.ServiceEndpointMethodFilter{}, err
	} else {
		var rv []masherytypes.ServiceEndpointMethodFilter
		for _, raw := range d {
			ms, ok := raw.([]masherytypes.ServiceEndpointMethodFilter)
			if ok {
				rv = append(rv, ms...)
			}
		}

		for _, v := range rv {
			v.ServiceEndpointMethod = ident
		}

		return rv, nil
	}
}

// ListEndpointMethodFiltersWithFullInfo List endpoints methods filters with their extended information.
func ListEndpointMethodFiltersWithFullInfo(ctx context.Context, ident masherytypes.ServiceEndpointMethodIdentifier, c *transport.V3Transport) ([]masherytypes.ServiceEndpointMethodFilter, error) {
	spec := transport.FetchSpec{
		Pagination: transport.PerPage,
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods/%s/responseFilters", ident.ServiceId, ident.EndpointId, ident.MethodId),
		Query: url.Values{
			"fields": {MasheryResponseFilterFieldsStr},
		},
		AppContext:     "endpoint methods filters",
		ResponseParser: masherytypes.ParseServiceEndpointMethodFilterArray,
	}

	if d, err := c.FetchAll(ctx, spec); err != nil {
		return []masherytypes.ServiceEndpointMethodFilter{}, err
	} else {
		// Convert individual fetches into the array of elements
		var rv []masherytypes.ServiceEndpointMethodFilter
		for _, raw := range d {
			ms, ok := raw.([]masherytypes.ServiceEndpointMethodFilter)
			if ok {
				rv = append(rv, ms...)
			}
		}

		for _, v := range rv {
			v.ServiceEndpointMethod = ident
		}

		return rv, nil
	}
}

// CreateEndpointMethodFilter Create a new service.
func CreateEndpointMethodFilter(ctx context.Context, ident masherytypes.ServiceEndpointMethodIdentifier,
	filterUpsert masherytypes.ServiceEndpointMethodFilter,
	c *transport.V3Transport) (*masherytypes.ServiceEndpointMethodFilter, error) {

	rawResp, err := c.CreateObject(ctx, filterUpsert, transport.FetchSpec{
		Resource: fmt.Sprintf("/services/%s/endpoints/%s/methods/%s/responseFilters",
			ident.ServiceId,
			ident.EndpointId,
			ident.MethodId),
		AppContext: "endpoint method filters",
		Query: url.Values{
			"fields": {MasheryResponseFilterFieldsStr},
		},
		ResponseParser: masherytypes.ParseServiceEndpointMethodFilter,
	})

	if err == nil {
		rv, _ := rawResp.(masherytypes.ServiceEndpointMethodFilter)
		rv.ServiceEndpointMethod = ident
		return &rv, nil
	} else {
		return nil, err
	}
}

// UpdateEndpointMethodFilter Update mashery endpoint method using the specified upsertable.
func UpdateEndpointMethodFilter(ctx context.Context, methUpsert masherytypes.ServiceEndpointMethodFilter,
	c *transport.V3Transport) (*masherytypes.ServiceEndpointMethodFilter, error) {
	if methUpsert.Id == "" {
		return nil, errors.New("illegal argument: response filter must be set and not nil")
	}

	opContext := transport.FetchSpec{
		Resource: fmt.Sprintf("/services/%s/endpoints/%s/methods/%s/responseFilters/%s",
			methUpsert.ServiceEndpointMethod.ServiceId,
			methUpsert.ServiceEndpointMethod.EndpointId,
			methUpsert.ServiceEndpointMethod.MethodId,
			methUpsert.Id),
		AppContext: "endpoint method filters",
		Query: url.Values{
			"fields": {MasheryResponseFilterFieldsStr},
		},
		ResponseParser: masherytypes.ParsePackagePlanServiceEndpointMethodFilter,
	}

	if d, err := c.UpdateObject(ctx, methUpsert, opContext); err == nil {
		rv, _ := d.(masherytypes.ServiceEndpointMethodFilter)
		rv.ServiceEndpointMethod = methUpsert.ServiceEndpointMethod
		return &rv, nil
	} else {
		return nil, err
	}
}

func GetEndpointMethodFilter(ctx context.Context, ident masherytypes.ServiceEndpointMethodFilterIdentifier,
	c *transport.V3Transport) (*masherytypes.ServiceEndpointMethodFilter, error) {
	fetchSpec := transport.FetchSpec{
		Pagination: transport.NotRequired,
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods/%s/responseFilters/%s", ident.ServiceId, ident.EndpointId, ident.MethodId, ident.FilterId),
		Query: url.Values{
			"fields": {MasheryResponseFilterFieldsStr},
		},
		AppContext:     "endpoint method filters",
		ResponseParser: masherytypes.ParseServiceEndpointMethodFilter,
	}

	if raw, err := c.GetObject(ctx, fetchSpec); err != nil {
		return nil, err
	} else {
		if rv, ok := raw.(masherytypes.ServiceEndpointMethodFilter); ok {
			rv.ServiceEndpointMethod = ident.ServiceEndpointMethodIdentifier
			return &rv, nil
		} else {
			return nil, errors.New("invalid return type")
		}
	}
}

func DeleteEndpointMethodFilter(ctx context.Context, ident masherytypes.ServiceEndpointMethodFilterIdentifier, c *transport.V3Transport) error {
	return c.DeleteObject(ctx, transport.FetchSpec{
		Resource: fmt.Sprintf("/services/%s/endpoints/%s/methods/%s/responseFilters/%s",
			ident.ServiceId, ident.EndpointId, ident.MethodId, ident.FilterId),
		AppContext: "endpoint method filters",
	})
}

// CountEndpointsMethodsFiltersOf Count the number of services that would match this criteria
func CountEndpointsMethodsFiltersOf(ctx context.Context, ident masherytypes.ServiceEndpointMethodIdentifier, c *transport.V3Transport) (int64, error) {
	opCtx := transport.FetchSpec{
		Pagination: transport.NotRequired,
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods/%s/responseFilters", ident.ServiceId, ident.EndpointId, ident.MethodId),
		AppContext: "endpoint method filters",
	}

	return c.Count(ctx, opCtx)
}
