package v3client

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"net/url"
)

func GetRole(ctx context.Context, id string, c *transport.V3Transport) (*masherytypes.Role, error) {
	rv, err := c.GetObject(ctx, transport.FetchSpec{
		Resource:       fmt.Sprintf("/roles/%s", id),
		Query:          nil,
		AppContext:     "role",
		ResponseParser: masherytypes.ParseService,
	})

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(masherytypes.Role)
		return &retServ, nil
	}
}

func ListRoles(ctx context.Context, c *transport.V3Transport) ([]masherytypes.Role, error) {
	return listRoles(ctx, nil, c)
}

func ListRolesFiltered(ctx context.Context, params map[string]string, fields []string, c *transport.V3Transport) ([]masherytypes.Role, error) {
	return listRoles(ctx, c.V3FilteringParams(params, fields), c)
}

func listRoles(ctx context.Context, qs url.Values, c *transport.V3Transport) ([]masherytypes.Role, error) {
	opCtx := transport.FetchSpec{
		Pagination:     transport.PerPage,
		Resource:       "/roles",
		Query:          qs,
		AppContext:     "all roles",
		ResponseParser: masherytypes.ParseMasheryRoleArray,
	}

	if d, err := c.FetchAll(ctx, opCtx); err != nil {
		return []masherytypes.Role{}, nil
	} else {
		// Convert individual fetches into the array of elements
		var rv []masherytypes.Role
		for _, raw := range d {
			ms, ok := raw.([]masherytypes.Role)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}
