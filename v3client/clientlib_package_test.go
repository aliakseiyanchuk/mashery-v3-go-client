package v3client

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"testing"
)

func TestGetPackage(t *testing.T) {
	memberId := masherytypes.PackageIdentifier{PackageId: "package-id"}

	expRvPackage := masherytypes.Package{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   "member-id",
			Name: "member-name",
		},
		Description: "desc-1",
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packages/package-id").
			WithMethod("get").
			RequestingFields(MasheryPackageFields).
			WillReturnJsonOf(expRvPackage)
	}

	autoTestGet(t,
		memberId,
		expRvPackage,
		mockVisitor,
		func(cl Client) ClientBoolExchangeFunc[masherytypes.PackageIdentifier, masherytypes.Package] {
			return cl.GetPackage
		},
	)
}

func TestCreatePackage(t *testing.T) {

	expRvCreatePayload := masherytypes.Package{
		AddressableV3Object: masherytypes.AddressableV3Object{Id: "", Name: "member-name"},
		Description:         "desc-1",
	}

	apiResponseJson := cloneWithModification(expRvCreatePayload, func(t1 *masherytypes.Package) { t1.Id = "package-id" })

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packages").
			WithMethod("post").
			RequestingFields(MasheryPackageFields).
			Matching(PayloadMatcher(expRvCreatePayload)).
			WillReturnJsonOf(apiResponseJson)
	}

	autoTestRootCreate(t,
		expRvCreatePayload,
		apiResponseJson,
		mockVisitor,
		func(cl Client) ClientExchangeFunc[masherytypes.Package, masherytypes.Package] {
			return cl.CreatePackage
		},
	)
}

func TestUpdatePackage(t *testing.T) {
	payload := masherytypes.Package{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   "package-id",
			Name: "package-name",
		},
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packages/package-id").
			WithMethod("put").
			RequestingFields(MasheryPackageFields).
			Matching(PayloadMatcher(payload)).
			WillReturnJsonOf(payload)
	}

	autoTestUpdate(t,
		payload,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.Package, masherytypes.Package] {
			return client.UpdatePackage
		},
	)
}

func TestDeletePackage(t *testing.T) {
	postIdent := masherytypes.PackageIdentifier{
		PackageId: "package-id",
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packages/package-id").
			WithMethod("delete").
			RequestingNoFields().
			WillReturnJsonOf(postIdent)
	}

	autoTestDelete(t,
		postIdent,
		mockVisitor,
		func(cl Client) BiConsumerCanErr[context.Context, masherytypes.PackageIdentifier] {
			return cl.DeletePackage
		},
	)
}

func TestListPackages(t *testing.T) {
	mockedResponse := []masherytypes.Package{
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "package-id"}},
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packages").
			WithMethod("get").
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
