package v3client

import (
	"context"
	"fmt"
)

func masheryServiceRolesSpec(id string) FetchSpec {
	return FetchSpec{
		Resource:       fmt.Sprintf("/services/%s/roles", id),
		Query:          nil,
		AppContext:     "get service roles",
		ResponseParser: ParseMasheryRolePermissionArray,
		Return404AsNil: true,
	}
}

func masheryServiceRolesPutSpec(id string) FetchSpec {
	return FetchSpec{
		Resource:       fmt.Sprintf("/services/%s", id),
		Query:          nil,
		AppContext:     "put service roles",
		ResponseParser: nilParser,
	}
}

// GetServiceRoles retrieve the roles that are attached to this service.
func GetServiceRoles(ctx context.Context, id string, c *HttpTransport) ([]MasheryRolePermission, error) {
	d, err := c.fetchAll(ctx, masheryServiceRolesSpec(id))

	if err != nil {
		return nil, err
	} else {
		var rv []MasheryRolePermission
		for _, raw := range d {
			ms, ok := raw.([]MasheryRolePermission)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}

type setServiceRolesWrapper struct {
	Roles []MasheryRolePermission `json:"roles"`
}

// SetServiceRoles set service roles for the given service. Empty array effectively deletes all associated roles.
func SetServiceRoles(ctx context.Context, id string, roles []MasheryRolePermission, c *HttpTransport) error {
	wrappedUpsert := setServiceRolesWrapper{Roles: roles}

	_, err := c.updateObject(ctx, wrappedUpsert, masheryServiceRolesPutSpec(id))

	if err == nil {
		return nil
	} else {
		return err
	}
}
