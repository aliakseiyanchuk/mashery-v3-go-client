package v3client

import (
	"context"
	"fmt"
	"net/url"
)

func (c *HttpTransport) GetEmailTemplateSet(ctx context.Context, id string) (*MasheryEmailTemplateSet, error) {
	rv, err := c.getObject(ctx, FetchSpec{
		Resource: fmt.Sprintf("/emailTemplateSets/%s", id),
		Query: url.Values{
			"fields": {MasheryEmailTemplateSetFieldsStr},
		},
		AppContext:     "email template set",
		ResponseParser: ParseMasheryEmailTemplateSet,
	})

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(MasheryEmailTemplateSet)
		return &retServ, nil
	}
}

func (c *HttpTransport) ListEmailTemplateSets(ctx context.Context) ([]MasheryEmailTemplateSet, error) {
	return c.listEmailTemplateSet(ctx, nil)
}

func (c *HttpTransport) ListEmailTemplateSetsFiltered(ctx context.Context, params map[string]string, fields []string) ([]MasheryEmailTemplateSet, error) {
	return c.listEmailTemplateSet(ctx, c.v3FilteringParams(params, fields))
}

func (c *HttpTransport) listEmailTemplateSet(ctx context.Context, qs url.Values) ([]MasheryEmailTemplateSet, error) {
	opCtx := FetchSpec{
		Pagination:     PerPage,
		Resource:       "/emailTemplateSets",
		Query:          qs,
		AppContext:     "all email template sets",
		ResponseParser: ParseMasheryEmailTemplateSetArray,
	}

	if d, err := c.fetchAll(ctx, opCtx); err != nil {
		return []MasheryEmailTemplateSet{}, nil
	} else {
		// Convert individual fetches into the array of elements
		var rv []MasheryEmailTemplateSet
		for _, raw := range d {
			ms, ok := raw.([]MasheryEmailTemplateSet)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}
