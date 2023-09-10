package v3client

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"testing"
)

func TestGetRole(t *testing.T) {
	expRvRoles := masherytypes.Role{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   "abc",
			Name: "name",
		},
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/roles/abc").
			WithMethod("get").
			RequestingNoFields().
			WillReturnJsonOf(expRvRoles)
	}

	autoTestGet(
		t,
		"abc",
		expRvRoles,
		mockVisitor,
		func(cl Client) ClientBoolExchangeFunc[string, masherytypes.Role] {
			return cl.GetRole
		},
	)
}

func TestListRoles(t *testing.T) {
	expRvApp := []masherytypes.Role{
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "R1"}},
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "R2"}},
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "R3"}},
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/roles").
			WithMethod("get").
			RequestingNoFields().
			WillReturnJsonOf(expRvApp)
	}

	autoTestRootFetchAll(
		t,
		expRvApp,
		mockVisitor,
		func(cl Client) ClientArraySupplierFunc[masherytypes.Role] {
			return cl.ListRoles
		},
	)
}

func TestListRolesFiltered(t *testing.T) {
	expRvApp := []masherytypes.Role{
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "R4"}},
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "R5"}},
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "R6"}},
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/roles").
			WithMethod("get").
			FilteredOn("a", "b").
			RequestingNoFields().
			WillReturnJsonOf(expRvApp)
	}

	autoTestRootFetchFiltered(
		t,
		map[string]string{"a": "b"},
		expRvApp,
		mockVisitor,
		func(client Client) ClientFilteredArraySupplierFunc[masherytypes.Role] {
			return client.ListRolesFiltered
		},
	)
}
