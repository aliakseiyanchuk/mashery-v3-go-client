package v3client

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"testing"
)

func TestSetServiceRoles(t *testing.T) {
	serviceIdent := masherytypes.ServiceIdentifier{}
	serviceIdent.ServiceId = "service-id"

	expBody := []masherytypes.RolePermission{
		{
			Role: masherytypes.Role{
				AddressableV3Object: masherytypes.AddressableV3Object{Id: "role-id"},
			},
			Action: "read",
		},
	}

	onTheWire := setServiceRolesWrapper{
		Roles: expBody,
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/roles").
			WithMethod("put").
			RequestingNoFilters().
			Matching(PayloadMatcher(onTheWire)).
			WillReturnJsonOf(expBody)
	}

	autoTestBiConsume(t,
		serviceIdent,
		expBody,
		mockVisitor,
		func(client Client) ClientBiConsumerFunc[masherytypes.ServiceIdentifier, []masherytypes.RolePermission] {
			return client.SetServiceRoles
		},
	)
}

func TestGetServiceRoles(t *testing.T) {
	serviceId := masherytypes.ServiceIdentifier{ServiceId: "service-id"}

	onTheWire := []masherytypes.RolePermission{
		{
			Role: masherytypes.Role{
				AddressableV3Object: masherytypes.AddressableV3Object{Id: "role-id"},
			},
			Action: "read",
		},
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/roles").
			WithMethod("get").
			RequestingNoFields().
			WillReturnJsonOf(onTheWire)
	}

	autoTestGet(t,
		serviceId,
		onTheWire,
		mockVisitor,
		func(cl Client) ClientBoolExchangeFunc[masherytypes.ServiceIdentifier, []masherytypes.RolePermission] {
			return cl.GetServiceRoles
		},
	)
}

func TestDeleteServiceRoles(t *testing.T) {
	serviceId := masherytypes.ServiceIdentifier{ServiceId: "service-id"}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/services/service-id/roles").
			WithMethod("delete").
			RequestingNoFields().
			WillReturnUnspecified()
	}

	autoTestDelete(t,
		serviceId,
		mockVisitor,
		func(cl Client) BiConsumerCanErr[context.Context, masherytypes.ServiceIdentifier] {
			return cl.DeleteServiceRoles
		},
	)
}
