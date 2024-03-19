package v3client

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"testing"
)

func TestGetPublicDomains(t *testing.T) {
	expRvApp := []masherytypes.DomainAddress{
		{Address: "a"},
		{Address: "b"},
		{Address: "c"},
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/domains/public/hostnames").
			WithMethod("get").
			RequestingNoFields().
			WillReturnJsonOf(expRvApp)
	}

	autoTestRootFetchAll(
		t,
		expRvApp,
		mockVisitor,
		func(cl Client) ClientArraySupplierFunc[masherytypes.DomainAddress] {
			return cl.GetPublicDomains
		},
	)
}

func TestGetSystemDomains(t *testing.T) {
	expRvApp := []masherytypes.DomainAddress{
		{Address: "a"},
		{Address: "b"},
		{Address: "c"},
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/domains/system/hostnames").
			WithMethod("get").
			RequestingNoFields().
			WillReturnJsonOf(expRvApp)
	}

	autoTestRootFetchAll(
		t,
		expRvApp,
		mockVisitor,
		func(cl Client) ClientArraySupplierFunc[masherytypes.DomainAddress] {
			return cl.GetSystemDomains
		},
	)
}
