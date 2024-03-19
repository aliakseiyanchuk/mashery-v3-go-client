package v3client

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"testing"
)

func TestListOrganizations(t *testing.T) {
	expRvApp := []masherytypes.Organization{
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "R1"}},
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "R2"}},
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "R3"}},
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/organizations").
			WithMethod("get").
			RequestingNoFields().
			WillReturnJsonOf(expRvApp)
	}

	autoTestRootFetchAll(
		t,
		expRvApp,
		mockVisitor,
		func(cl Client) ClientArraySupplierFunc[masherytypes.Organization] {
			return cl.ListOrganizations
		},
	)
}

func TestListOrganizationsFiltered(t *testing.T) {
	expRvApp := []masherytypes.Organization{
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "R4"}},
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "R5"}},
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "R6"}},
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/organizations").
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
		func(client Client) ClientFilteredArraySupplierFunc[masherytypes.Organization] {
			return client.ListOrganizationsFiltered
		},
	)
}
