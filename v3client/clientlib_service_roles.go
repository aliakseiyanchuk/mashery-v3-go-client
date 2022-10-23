package v3client

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
)

func masheryServiceRolesSpec(id masherytypes.ServiceIdentifier) transport.FetchSpec {
	return transport.FetchSpec{
		Resource:       fmt.Sprintf("/services/%s/roles", id.ServiceId),
		Query:          nil,
		AppContext:     "get service roles",
		ResponseParser: masherytypes.ParseRolePermissionArray,
		Return404AsNil: true,
	}
}

func masheryServiceRolesPutSpec(id masherytypes.ServiceIdentifier) transport.FetchSpec {
	return transport.FetchSpec{
		Resource:       fmt.Sprintf("/services/%s", id),
		Query:          nil,
		AppContext:     "put service roles",
		ResponseParser: masherytypes.NilParser,
	}
}

// GetServiceRoles retrieve the roles that are attached to this service.
func GetServiceRoles(ctx context.Context, id masherytypes.ServiceIdentifier, c *transport.V3Transport) ([]masherytypes.MasheryRolePermission, error) {
	d, err := c.FetchAll(ctx, masheryServiceRolesSpec(id))

	if err != nil {
		return nil, err
	} else {
		var rv []masherytypes.MasheryRolePermission
		for _, raw := range d {
			ms, ok := raw.([]masherytypes.MasheryRolePermission)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}

type setServiceRolesWrapper struct {
	Roles []masherytypes.MasheryRolePermission `json:"roles"`
}

// SetServiceRoles set service roles for the given service. Empty array effectively deletes all associated roles.
func SetServiceRoles(ctx context.Context, id masherytypes.ServiceIdentifier, roles []masherytypes.MasheryRolePermission, c *transport.V3Transport) error {
	wrappedUpsert := setServiceRolesWrapper{Roles: roles}

	_, err := c.UpdateObject(ctx, wrappedUpsert, masheryServiceRolesPutSpec(id))

	if err == nil {
		return nil
	} else {
		return err
	}
}

// DeleteServiceRoles delete service roles
func DeleteServiceRoles(ctx context.Context, id masherytypes.ServiceIdentifier, c *transport.V3Transport) error {
	return c.DeleteObject(ctx, masheryServiceRolesSpec(id))
}
