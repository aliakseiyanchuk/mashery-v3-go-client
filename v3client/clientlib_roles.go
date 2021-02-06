package v3client

import (
	"context"
	"fmt"
	"net/url"
)

func (c *HttpTransport) GetRole(ctx context.Context, id string) (*MasheryRole, error) {
	rv, err := c.getObject(ctx, FetchSpec{
		Resource:       fmt.Sprintf("/roles/%s", id),
		Query:          nil,
		AppContext:     "role",
		ResponseParser: ParseMasheryService,
	})

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(MasheryRole)
		return &retServ, nil
	}
}

func (c *HttpTransport) ListRoles(ctx context.Context) ([]MasheryRole, error) {
	return c.listRoles(ctx, nil)
}

func (c *HttpTransport) ListRolesFiltered(ctx context.Context, params map[string]string, fields []string) ([]MasheryRole, error) {
	return c.listRoles(ctx, c.v3FilteringParams(params, fields))
}

func (c *HttpTransport) listRoles(ctx context.Context, qs url.Values) ([]MasheryRole, error) {
	opCtx := FetchSpec{
		Pagination:     PerPage,
		Resource:       "/roles",
		Query:          qs,
		AppContext:     "all roles",
		ResponseParser: ParseMasheryRoleArray,
	}

	if d, err := c.fetchAll(ctx, opCtx); err != nil {
		return []MasheryRole{}, nil
	} else {
		// Convert individual fetches into the array of elements
		var rv []MasheryRole
		for _, raw := range d {
			ms, ok := raw.([]MasheryRole)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}
