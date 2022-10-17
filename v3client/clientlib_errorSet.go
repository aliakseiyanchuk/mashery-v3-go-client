package v3client

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"net/url"
)

func ListErrorSets(ctx context.Context, serviceId masherytypes.ServiceIdentifier, qs url.Values, c *transport.V3Transport) ([]masherytypes.ErrorSet, error) {
	opCtx := transport.FetchSpec{
		Pagination:     transport.PerItem,
		Resource:       fmt.Sprintf("/services/%s/errorSets", serviceId.ServiceId),
		Query:          qs,
		AppContext:     "all error sets of a service",
		ResponseParser: masherytypes.ParseServiceErrorSetArray,
	}

	if d, err := c.FetchAll(ctx, opCtx); err != nil {
		return []masherytypes.ErrorSet{}, nil
	} else {
		// Convert individual fetches into the array of elements
		var rv []masherytypes.ErrorSet
		for _, raw := range d {
			ms, ok := raw.([]masherytypes.ErrorSet)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}

func GetErrorSet(ctx context.Context, ident masherytypes.ErrorSetIdentifier, c *transport.V3Transport) (*masherytypes.ErrorSet, error) {
	opCtx := transport.FetchSpec{
		Resource: fmt.Sprintf("/services/%s/errorSets/%s", ident.ServiceId, ident.ErrorSetId),
		Query: map[string][]string{
			"fields": MasheryErrorSetFields,
		},
		AppContext:     "specific error sets",
		ResponseParser: masherytypes.ParseErrorSet,
	}

	if raw, err := c.GetObject(ctx, opCtx); err != nil {
		return nil, err
	} else {
		if rv, ok := raw.(masherytypes.ErrorSet); ok {
			return &rv, nil
		} else {
			return nil, errors.New("invalid return type")
		}
	}
}

func CreateErrorSet(ctx context.Context, serviceId masherytypes.ServiceIdentifier, set masherytypes.ErrorSet, c *transport.V3Transport) (*masherytypes.ErrorSet, error) {
	rawResp, err := c.CreateObject(ctx, set, transport.FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/errorSets", serviceId),
		AppContext: "create errorSet",

		Query: url.Values{
			"fields": MasheryErrorSetFields,
		},
		ResponseParser: masherytypes.ParseErrorSet,
	})

	if err == nil {
		rv, _ := rawResp.(masherytypes.ErrorSet)
		return &rv, nil
	} else {
		return nil, err
	}
}

func UpdateErrorSet(ctx context.Context, setData masherytypes.ErrorSet, c *transport.V3Transport) (*masherytypes.ErrorSet, error) {
	if setData.Id == "" {
		return nil, errors.New("illegal argument: endpoint Id must be set and not nil")
	}

	opContext := transport.FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/errorSets/%s", setData.ParentServiceId.ServiceId, setData.Id),
		AppContext: "update error set",
		Query: url.Values{
			"fields": {MasheryEndpointFieldsStr},
		},
		ResponseParser: masherytypes.ParseErrorSet,
	}

	if d, err := c.UpdateObject(ctx, setData, opContext); err == nil {
		rv, _ := d.(masherytypes.ErrorSet)
		return &rv, nil
	} else {
		return nil, err
	}
}

func UpdateErrorSetMessage(ctx context.Context, ident masherytypes.ErrorSetIdentifier, msg masherytypes.MasheryErrorMessage, c *transport.V3Transport) (*masherytypes.MasheryErrorMessage, error) {
	if msg.Id == "" {
		return nil, errors.New("illegal argument: message Id must not be empty")
	}

	opContext := transport.FetchSpec{
		Resource:       fmt.Sprintf("/services/%s/errorSets/%s/errorMessages/%s", ident.ServiceId, ident.ErrorSetId, msg.Id),
		AppContext:     "update error set message",
		ResponseParser: masherytypes.ParseErrorSetMessage,
	}

	if d, err := c.UpdateObject(ctx, msg, opContext); err == nil {
		rv, _ := d.(masherytypes.MasheryErrorMessage)
		return &rv, nil
	} else {
		return nil, err
	}
}

func DeleteErrorSet(ctx context.Context, ident masherytypes.ErrorSetIdentifier, c *transport.V3Transport) error {
	opContext := transport.FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/errorSets/%s", ident.ServiceId, ident.ErrorSetId),
		AppContext: "service error set",
	}

	return c.DeleteObject(ctx, opContext)
}
