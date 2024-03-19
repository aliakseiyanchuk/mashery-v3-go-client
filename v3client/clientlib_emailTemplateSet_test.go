package v3client

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"testing"
)

func TestGetEmailTemplateSet(t *testing.T) {
	expRvTemplateSet := masherytypes.EmailTemplateSet{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   "abc",
			Name: "name",
		},
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/emailTemplateSets/abc").
			WithMethod("get").
			RequestingFields(MasheryEmailTemplateSetFields).
			WillReturnJsonOf(expRvTemplateSet)
	}

	autoTestGet(
		t,
		"abc",
		expRvTemplateSet,
		mockVisitor,
		func(cl Client) ClientBoolExchangeFunc[string, masherytypes.EmailTemplateSet] {
			return cl.GetEmailTemplateSet
		},
	)
}

func TestListEmailTemplateSets(t *testing.T) {
	expRvApp := []masherytypes.EmailTemplateSet{
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "A"}},
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "B"}},
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "C"}},
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/emailTemplateSets").
			WithMethod("get").
			RequestingFields(MasheryEmailTemplateSetFields).
			WillReturnJsonOf(expRvApp)
	}

	autoTestRootFetchAll(
		t,
		expRvApp,
		mockVisitor,
		func(cl Client) ClientArraySupplierFunc[masherytypes.EmailTemplateSet] {
			return cl.ListEmailTemplateSets
		},
	)
}

func TestListEmailTemplateSetsFiltered(t *testing.T) {
	expRvApp := []masherytypes.EmailTemplateSet{
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "A"}},
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "B"}},
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "C"}},
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/emailTemplateSets").
			WithMethod("get").
			FilteredOn("a", "b").
			RequestingFields(MasheryEmailTemplateSetFields).
			WillReturnJsonOf(expRvApp)
	}

	autoTestRootFetchFiltered(
		t,
		map[string]string{"a": "b"},
		expRvApp,
		mockVisitor,
		func(client Client) ClientFilteredArraySupplierFunc[masherytypes.EmailTemplateSet] {
			return client.ListEmailTemplateSetsFiltered
		},
	)
}
