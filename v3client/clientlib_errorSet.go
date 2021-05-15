package v3client

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

func ListErrorSets(ctx context.Context, serviceId string, qs url.Values, c *HttpTransport) ([]MasheryErrorSet, error) {
	opCtx := FetchSpec{
		Pagination:     PerItem,
		Resource:       fmt.Sprintf("/services/%s/errorSets", serviceId),
		Query:          qs,
		AppContext:     "all error sets of a service",
		ResponseParser: ParseServiceErrorSetArray,
	}

	if d, err := c.fetchAll(ctx, opCtx); err != nil {
		return []MasheryErrorSet{}, nil
	} else {
		// Convert individual fetches into the array of elements
		var rv []MasheryErrorSet
		for _, raw := range d {
			ms, ok := raw.([]MasheryErrorSet)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}

func GetErrorSet(ctx context.Context, serviceId, setId string, c *HttpTransport) (*MasheryErrorSet, error) {
	opCtx := FetchSpec{
		Resource: fmt.Sprintf("/services/%s/errorSets/%s", serviceId, setId),
		Query: map[string][]string{
			"fields": MasheryErrorSetFields,
		},
		AppContext:     "specific error sets",
		ResponseParser: ParseErrorSet,
	}

	if raw, err := c.getObject(ctx, opCtx); err != nil {
		return nil, err
	} else {
		if rv, ok := raw.(MasheryErrorSet); ok {
			return &rv, nil
		} else {
			return nil, errors.New("invalid return type")
		}
	}
}

func CreateErrorSet(ctx context.Context, serviceId string, set MasheryErrorSet, c *HttpTransport) (*MasheryErrorSet, error) {
	rawResp, err := c.createObject(ctx, set, FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/errorSets", serviceId),
		AppContext: "create errorSet",

		Query: url.Values{
			"fields": MasheryErrorSetFields,
		},
		ResponseParser: ParseErrorSet,
	})

	if err == nil {
		rv, _ := rawResp.(MasheryErrorSet)
		return &rv, nil
	} else {
		return nil, err
	}
}

func UpdateErrorSet(ctx context.Context, serviceId string, setData MasheryErrorSet, c *HttpTransport) (*MasheryErrorSet, error) {
	if setData.Id == "" {
		return nil, errors.New("illegal argument: endpoint Id must be set and not nil")
	}

	opContext := FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/errorSets/%s", serviceId, setData.Id),
		AppContext: "update error set",
		Query: url.Values{
			"fields": {MasheryEndpointFieldsStr},
		},
		ResponseParser: ParseErrorSet,
	}

	if d, err := c.updateObject(ctx, setData, opContext); err == nil {
		rv, _ := d.(MasheryErrorSet)
		return &rv, nil
	} else {
		return nil, err
	}
}

func UpdateErrorSetMessage(ctx context.Context, serviceId string, setId string, msg MasheryErrorMessage, c *HttpTransport) (*MasheryErrorMessage, error) {
	if msg.Id == "" {
		return nil, errors.New("illegal argument: message Id must not be empty")
	}

	opContext := FetchSpec{
		Resource:       fmt.Sprintf("/services/%s/errorSets/%s/errorMessages/%s", serviceId, setId, msg.Id),
		AppContext:     "update error set message",
		ResponseParser: ParseErrorSetMessage,
	}

	if d, err := c.updateObject(ctx, msg, opContext); err == nil {
		rv, _ := d.(MasheryErrorMessage)
		return &rv, nil
	} else {
		return nil, err
	}
}

func DeleteErrorSet(ctx context.Context, serviceId, setId string, c *HttpTransport) error {
	opContext := FetchSpec{
		Resource:   fmt.Sprintf("/services/%s/errorSets/%s", serviceId, setId),
		AppContext: "service error set",
	}

	return c.deleteObject(ctx, opContext)
}
