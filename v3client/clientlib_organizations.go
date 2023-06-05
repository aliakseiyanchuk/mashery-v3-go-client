package v3client

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
)

func ListOrganizations(ctx context.Context, c *transport.V3Transport) ([]masherytypes.Organization, error) {
	return ListOrganizationsFiltered(ctx, nil, c)
}

func ListOrganizationsFiltered(ctx context.Context, qs map[string]string, c *transport.V3Transport) ([]masherytypes.Organization, error) {
	opCtx := transport.FetchSpec{
		Pagination:     transport.PerPage,
		Resource:       "/organizations",
		Query:          c.V3FilteringParams(qs, nil),
		AppContext:     "all organizations",
		ResponseParser: masherytypes.ParseMasheryOriganizationsArray,
	}

	if d, err := c.FetchAll(ctx, opCtx); err != nil {
		return []masherytypes.Organization{}, err
	} else {
		// Convert individual fetches into the array of elements
		var rv []masherytypes.Organization
		for _, raw := range d {
			ms, ok := raw.([]masherytypes.Organization)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}
