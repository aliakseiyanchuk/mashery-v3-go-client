package v3client

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"testing"
)

func TestGetPackageKey(t *testing.T) {
	keyId := masherytypes.PackageKeyIdentifier{PackageKeyId: "package-key-id"}

	expRvPackage := masherytypes.PackageKey{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   "member-id",
			Name: "member-name",
		},
		Status: "active",
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packageKeys/package-key-id").
			WithMethod("get").
			RequestingFields(MasheryPackageKeyFields).
			WillReturnJsonOf(expRvPackage)
	}

	autoTestGet(t,
		keyId,
		expRvPackage,
		mockVisitor,
		func(cl Client) ClientBoolExchangeFunc[masherytypes.PackageKeyIdentifier, masherytypes.PackageKey] {
			return cl.GetPackageKey
		},
	)
}

func TestCreatePackageKey(t *testing.T) {

	expAppCtx := masherytypes.ApplicationIdentifier{ApplicationId: "app-id"}
	expRvCreatePayload := masherytypes.PackageKey{
		AddressableV3Object: masherytypes.AddressableV3Object{Id: "", Name: "key-name"},
		Status:              "active",
		Plan: &masherytypes.Plan{
			AddressableV3Object: masherytypes.AddressableV3Object{Id: "plan-id"},
		},
		Package: &masherytypes.Package{
			AddressableV3Object: masherytypes.AddressableV3Object{Id: "package-id"},
		},
	}

	apiResponseJson := cloneWithModification(expRvCreatePayload, func(t1 *masherytypes.PackageKey) { t1.Id = "package-key-id" })

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/applications/app-id/packageKeys").
			WithMethod("post").
			RequestingFields(MasheryPackageKeyFields).
			Matching(PayloadMatcher(expRvCreatePayload)).
			WillReturnJsonOf(apiResponseJson)
	}

	autoTestCreate(t,
		expAppCtx,
		expRvCreatePayload,
		apiResponseJson,
		mockVisitor,
		func(client Client) ClientDualExchangeFunc[masherytypes.ApplicationIdentifier, masherytypes.PackageKey, masherytypes.PackageKey] {
			return client.CreatePackageKey
		},
	)
}

func TestUpdatePackageKey(t *testing.T) {
	payload := masherytypes.PackageKey{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   "package-key-id",
			Name: "package-key-name",
		},
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packageKeys/package-key-id").
			WithMethod("put").
			RequestingFields(MasheryPackageKeyFields).
			Matching(PayloadMatcher(payload)).
			WillReturnJsonOf(payload)
	}

	autoTestUpdate(t,
		payload,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.PackageKey, masherytypes.PackageKey] {
			return client.UpdatePackageKey
		},
	)
}

func TestDeletePackageKey(t *testing.T) {
	postIdent := masherytypes.PackageKeyIdentifier{
		PackageKeyId: "package-key-id",
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packageKeys/package-key-id").
			WithMethod("delete").
			RequestingNoFields().
			WillReturnJsonOf(postIdent)
	}

	autoTestDelete(t,
		postIdent,
		mockVisitor,
		func(cl Client) BiConsumerCanErr[context.Context, masherytypes.PackageKeyIdentifier] {
			return cl.DeletePackageKey
		},
	)
}

func TestListPackagesKey(t *testing.T) {
	mockedResponse := []masherytypes.Package{
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "package-id"}},
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packages").
			WithMethod("get").
			RequestingNoFilters().
			RequestingFields(MasheryPackageFields).
			WillReturnJsonOf(mockedResponse)
	}

	autoTestRootFetchAll(t,
		mockedResponse,
		mockVisitor,
		func(client Client) ClientArraySupplierFunc[masherytypes.Package] {
			return client.ListPackages
		},
	)
}

func TestListPackagesKeyFiltered(t *testing.T) {
	mockedResponse := []masherytypes.PackageKey{
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "package-id"}},
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packageKeys").
			WithMethod("get").
			FilteredOn("a", "b").
			RequestingFields(MasheryPackageKeyFields).
			WillReturnJsonOf(mockedResponse)
	}

	autoTestRootFetchFiltered(t,
		map[string]string{"a": "b"},
		mockedResponse,
		mockVisitor,
		func(client Client) ClientFilteredArraySupplierFunc[masherytypes.PackageKey] {
			return client.ListPackageKeysFiltered
		},
	)
}
