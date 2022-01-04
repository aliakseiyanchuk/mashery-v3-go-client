package v3client

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"net/url"
)

func GetEmailTemplateSet(ctx context.Context, id string, c *transport.V3Transport) (*masherytypes.MasheryEmailTemplateSet, error) {
	rv, err := c.GetObject(ctx, transport.FetchSpec{
		Resource: fmt.Sprintf("/emailTemplateSets/%s", id),
		Query: url.Values{
			"fields": {MasheryEmailTemplateSetFieldsStr},
		},
		AppContext:     "email template set",
		ResponseParser: masherytypes.ParseMasheryEmailTemplateSet,
	})

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(masherytypes.MasheryEmailTemplateSet)
		return &retServ, nil
	}
}

func ListEmailTemplateSets(ctx context.Context, c *transport.V3Transport) ([]masherytypes.MasheryEmailTemplateSet, error) {
	return listEmailTemplateSet(ctx, nil, c)
}

func ListEmailTemplateSetsFiltered(ctx context.Context, params map[string]string, fields []string, c *transport.V3Transport) ([]masherytypes.MasheryEmailTemplateSet, error) {
	return listEmailTemplateSet(ctx, c.V3FilteringParams(params, fields), c)
}

func listEmailTemplateSet(ctx context.Context, qs url.Values, c *transport.V3Transport) ([]masherytypes.MasheryEmailTemplateSet, error) {
	opCtx := transport.FetchSpec{
		Pagination:     transport.PerPage,
		Resource:       "/emailTemplateSets",
		Query:          qs,
		AppContext:     "all email template sets",
		ResponseParser: masherytypes.ParseMasheryEmailTemplateSetArray,
	}

	if d, err := c.FetchAll(ctx, opCtx); err != nil {
		return []masherytypes.MasheryEmailTemplateSet{}, nil
	} else {
		// Convert individual fetches into the array of elements
		var rv []masherytypes.MasheryEmailTemplateSet
		for _, raw := range d {
			ms, ok := raw.([]masherytypes.MasheryEmailTemplateSet)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}
