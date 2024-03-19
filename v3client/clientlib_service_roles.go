package v3client

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
)

// GetServiceRoles retrieve the roles that are attached to this service.
func GetServiceRoles(ctx context.Context, id masherytypes.ServiceIdentifier, c *transport.HttpTransport) ([]masherytypes.RolePermission, bool, error) {
	objectListSpecBuilder := transport.ObjectListFetchSpecBuilder[masherytypes.RolePermission]{}
	objectListSpecBuilder.
		WithValueFactory(func() []masherytypes.RolePermission {
			return []masherytypes.RolePermission{}
		}).
		WithReturn404AsNil(true).
		WithResource(fmt.Sprintf("/services/%s/roles", id.ServiceId)).
		WithAppContext("service role")

	return transport.FetchAllWithExists(ctx, objectListSpecBuilder.Build(), c)
}

type setServiceRolesWrapper struct {
	Roles []masherytypes.RolePermission `json:"roles"`
}

// SetServiceRoles set service roles for the given service. Empty array effectively deletes all associated roles.
func SetServiceRoles(ctx context.Context, id masherytypes.ServiceIdentifier, roles []masherytypes.RolePermission, c *transport.HttpTransport) error {
	wrappedUpsert := setServiceRolesWrapper{Roles: roles}

	objectUpsertSpecBuilder := transport.ObjectUpsertSpecBuilder[setServiceRolesWrapper]{}
	objectUpsertSpecBuilder.
		WithUpsert(wrappedUpsert).
		WithValueFactory(func() setServiceRolesWrapper {
			return setServiceRolesWrapper{}
		}).
		WithIgnoreResponse(true).
		WithResource("/services/%s/roles", id.ServiceId).
		WithAppContext("put service role")

	_, err := transport.UpdateObject(ctx, objectUpsertSpecBuilder.Build(), c)
	return err
}

// DeleteServiceRoles delete service roles
func DeleteServiceRoles(ctx context.Context, id masherytypes.ServiceIdentifier, c *transport.HttpTransport) error {
	objectUpsertSpecBuilder := transport.ObjectFetchSpecBuilder[masherytypes.RolePermission]{}
	objectUpsertSpecBuilder.
		WithValueFactory(func() masherytypes.RolePermission {
			return masherytypes.RolePermission{}
		}).
		WithIgnoreResponse(true).
		WithResource("/services/%s/roles", id.ServiceId).
		WithAppContext("delete service role")

	return transport.DeleteObject(ctx, objectUpsertSpecBuilder.Build(), c)
}
