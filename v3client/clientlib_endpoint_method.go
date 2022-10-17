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
func ListEndpointMethods(ctx context.Context, ident masherytypes.ServiceEndpointIdentifier, c *transport.V3Transport) ([]masherytypes.ServiceEndpointMethod, error) {
	spec := transport.FetchSpec{
		Pagination:     transport.PerPage,
		Resource:       fmt.Sprintf("/services/%s/endpoints/%s/methods", ident.ServiceIdentifier, ident.EndpointId),
		Query:          nil,
		AppContext:     "endpoint methods",
		ResponseParser: masherytypes.ParseServiceEndpointMethodArray,
	}

	if d, err := c.FetchAll(ctx, spec); err != nil {
		return []masherytypes.ServiceEndpointMethod{}, err
	} else {
		// Convert individual fetches into the array of elements
		var rv []masherytypes.ServiceEndpointMethod
		for _, raw := range d {
			ms, ok := raw.([]masherytypes.ServiceEndpointMethod)
			if ok {
				rv = append(rv, ms...)
			}

			for _, v := range rv {
				v.ParentEndpointId = ident
			}
		}

		return rv, nil
	}
}

// ListEndpointMethodsWithFullInfo List endpoints methods with their extended information.
func ListEndpointMethodsWithFullInfo(ctx context.Context, ident masherytypes.ServiceEndpointIdentifier, c *transport.V3Transport) ([]masherytypes.ServiceEndpointMethod, error) {
	spec := transport.FetchSpec{
		Pagination: transport.PerPage,
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods", ident.ServiceId, ident.EndpointId),
		Query: url.Values{
			"fields": {MasheryMethodsFieldsStr},
		},
		AppContext:     "endpoint methods",
		ResponseParser: masherytypes.ParseServiceEndpointMethodArray,
	}

	if d, err := c.FetchAll(ctx, spec); err != nil {
		return []masherytypes.ServiceEndpointMethod{}, err
	} else {
		// Convert individual fetches into the array of elements
		var rv []masherytypes.ServiceEndpointMethod
		for _, raw := range d {
			ms, ok := raw.([]masherytypes.ServiceEndpointMethod)
			if ok {
				rv = append(rv, ms...)
			}

			for _, v := range rv {
				v.ParentEndpointId = ident
			}
		}

		return rv, nil
	}
}

// CreateEndpointMethod Create a new service.
func CreateEndpointMethod(ctx context.Context, ident masherytypes.ServiceEndpointIdentifier, methodUpsert masherytypes.ServiceEndpointMethod, c *transport.V3Transport) (*masherytypes.ServiceEndpointMethod, error) {
	rawResp, err := c.CreateObject(ctx, methodUpsert, transport.FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods", ident.ServiceId, ident.EndpointId),
		AppContext: "endpoint method",
		Query: url.Values{
			"fields": {MasheryMethodsFieldsStr},
		},
		ResponseParser: masherytypes.ParseServiceEndpointMethod,
	})

	if err == nil {
		rv, _ := rawResp.(masherytypes.ServiceEndpointMethod)
		rv.ParentEndpointId = ident
		return &rv, nil
	} else {
		return nil, err
	}
}

// UpdateEndpointMethod Update mashery endpoint method using the specified upsertable.
func UpdateEndpointMethod(ctx context.Context, methUpsert masherytypes.ServiceEndpointMethod, c *transport.V3Transport) (*masherytypes.ServiceEndpointMethod, error) {
	if methUpsert.Id == "" {
		return nil, errors.New("illegal argument: endpoint Id must be set and not nil")
	}

	opContext := transport.FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods/%s", methUpsert.ParentEndpointId.ServiceId, methUpsert.ParentEndpointId.EndpointId, methUpsert.Id),
		AppContext: "endpoint method",
		Query: url.Values{
			"fields": {MasheryMethodsFieldsStr},
		},
		ResponseParser: masherytypes.ParseServiceEndpointMethod,
	}

	if d, err := c.UpdateObject(ctx, methUpsert, opContext); err == nil {
		rv, _ := d.(masherytypes.ServiceEndpointMethod)
		rv.ParentEndpointId = methUpsert.ParentEndpointId
		return &rv, nil
	} else {
		return nil, err
	}
}

func GetEndpointMethod(ctx context.Context, ident masherytypes.ServiceEndpointMethodIdentifier, c *transport.V3Transport) (*masherytypes.ServiceEndpointMethod, error) {
	fetchSpec := transport.FetchSpec{
		Pagination: transport.NotRequired,
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods/%s", ident.ServiceId, ident.EndpointId, ident.MethodId),
		Query: url.Values{
			"fields": {MasheryMethodsFieldsStr},
		},
		AppContext:     "endpoint method",
		ResponseParser: masherytypes.ParseServiceEndpointMethod,
	}

	if raw, err := c.GetObject(ctx, fetchSpec); err != nil {
		return nil, err
	} else {
		if rv, ok := raw.(masherytypes.ServiceEndpointMethod); ok {
			rv.ParentEndpointId = ident.ServiceEndpointIdentifier
			return &rv, nil
		} else {
			return nil, errors.New("invalid return type")
		}
	}
}

func DeleteEndpointMethod(ctx context.Context, ident masherytypes.ServiceEndpointMethodIdentifier, c *transport.V3Transport) error {
	return c.DeleteObject(ctx, transport.FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods/%s", ident.ServiceId, ident.EndpointId, ident.MethodId),
		AppContext: "endpoint method",
	})
}

// CountEndpointsMethodsOf Count the number of services that would match this criteria
func CountEndpointsMethodsOf(ctx context.Context, ident masherytypes.ServiceEndpointIdentifier, c *transport.V3Transport) (int64, error) {
	opCtx := transport.FetchSpec{
		Pagination: transport.NotRequired,
		Resource:   fmt.Sprintf("/services/%s/endpoints/%s/methods", ident.ServiceId, ident.EndpointId),
		AppContext: "endpoints methods",
	}

	return c.Count(ctx, opCtx)
}
