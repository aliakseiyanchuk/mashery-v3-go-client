package v3client

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"net/url"
)

var serviceAllFieldsQuery = url.Values{
	"fields": {MasheryServiceFullFieldsStr},
}

func GetService(ctx context.Context, id masherytypes.ServiceIdentifier, c *transport.V3Transport) (*masherytypes.Service, error) {
	rv, err := c.GetObject(ctx, transport.FetchSpec{
		Resource:       fmt.Sprintf("/services/%s", id.ServiceId),
		Query:          serviceAllFieldsQuery,
		AppContext:     "service",
		ResponseParser: masherytypes.ParseService,
	})

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(masherytypes.Service)
		return &retServ, nil
	}
}

// CreateService Create a new service.
func CreateService(ctx context.Context, service masherytypes.Service, c *transport.V3Transport) (*masherytypes.Service, error) {
	rawResp, err := c.CreateObject(ctx, service, transport.FetchSpec{
		Resource:   "/services",
		AppContext: "services",
		Query: url.Values{
			"fields": {MasheryServiceFullFieldsStr},
		},
		ResponseParser: masherytypes.ParseService,
	})

	if err == nil {
		rv, _ := rawResp.(masherytypes.Service)
		return &rv, nil
	} else {
		return nil, err
	}
}

// UpdateService updates a Mashery service
func UpdateService(ctx context.Context, service masherytypes.Service, c *transport.V3Transport) (*masherytypes.Service, error) {
	if service.Id == "" {
		return nil, errors.New("illegal argument: service Id must be set and not nil")
	}

	opContext := transport.FetchSpec{
		Resource:       fmt.Sprintf("/services/%s", service.Id),
		AppContext:     "service",
		ResponseParser: masherytypes.ParseService,
	}

	if d, err := c.UpdateObject(ctx, service, opContext); err == nil {
		rv, _ := d.(masherytypes.Service)
		return &rv, nil
	} else {
		return nil, err
	}
}

// DeleteService Delete a service.
func DeleteService(ctx context.Context, serviceId masherytypes.ServiceIdentifier, c *transport.V3Transport) error {
	opContext := transport.FetchSpec{
		Resource:   fmt.Sprintf("/services/%s", serviceId.ServiceId),
		AppContext: "service",
	}

	return c.DeleteObject(ctx, opContext)
}

// ListServicesFiltered List services that are filtered according to the condition that is V3-supported and containing the fields
// that the requester specifies
func ListServicesFiltered(ctx context.Context, params map[string]string, fields []string, c *transport.V3Transport) ([]masherytypes.Service, error) {
	return listServicesWithQuery(ctx, c.V3FilteringParams(params, fields), c)
}

func ListServices(ctx context.Context, c *transport.V3Transport) ([]masherytypes.Service, error) {
	return listServicesWithQuery(ctx, nil, c)
}

func listServicesWithQuery(ctx context.Context, qs url.Values, c *transport.V3Transport) ([]masherytypes.Service, error) {
	opCtx := transport.FetchSpec{
		Pagination:     transport.PerItem,
		Resource:       "/services",
		Query:          qs,
		AppContext:     "all service",
		ResponseParser: masherytypes.ParseMasheryServiceArray,
	}

	if d, err := c.FetchAll(ctx, opCtx); err != nil {
		return []masherytypes.Service{}, nil
	} else {
		// Convert individual fetches into the array of elements
		var rv []masherytypes.Service
		for _, raw := range d {
			ms, ok := raw.([]masherytypes.Service)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}

// CountServices Count the number of services that would match this criteria
func CountServices(ctx context.Context, params map[string]string, c *transport.V3Transport) (int64, error) {
	opCtx := transport.FetchSpec{
		Pagination: transport.NotRequired,
		Resource:   "/services",
		Query: url.Values{
			"filter": {transport.V3FilterExpression(params)},
		},
		AppContext:     "all service count",
		ResponseParser: masherytypes.ParseMasheryServiceArray,
	}

	return c.Count(ctx, opCtx)
}
