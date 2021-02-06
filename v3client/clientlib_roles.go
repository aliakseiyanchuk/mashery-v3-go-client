package v3client

import (
	"context"
	"fmt"
	"net/url"
)

func GetRole(ctx context.Context, id string, c *HttpTransport) (*MasheryRole, error) {
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

func ListRoles(ctx context.Context, c *HttpTransport) ([]MasheryRole, error) {
	return listRoles(ctx, nil, c)
}

func ListRolesFiltered(ctx context.Context, params map[string]string, fields []string, c *HttpTransport) ([]MasheryRole, error) {
	return listRoles(ctx, c.v3FilteringParams(params, fields), c)
}

func listRoles(ctx context.Context, qs url.Values, c *HttpTransport) ([]MasheryRole, error) {
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
