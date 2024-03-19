package v3client

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"testing"
)

func TestGetPackagePlanMethodFilter(t *testing.T) {
	methodIdent := masherytypes.PackagePlanServiceEndpointMethodIdentifier{}
	methodIdent.PackageId = "package-id"
	methodIdent.PlanId = "plan-id"
	methodIdent.ServiceId = "service-id"
	methodIdent.EndpointId = "endpoint-id"
	methodIdent.MethodId = "method-id"

	expRV := masherytypes.PackagePlanServiceEndpointMethodFilter{
		ResponseFilter: masherytypes.ResponseFilter{
			AddressableV3Object: masherytypes.AddressableV3Object{Id: "filter-id", Name: "filter-name"},
		},
		PackagePlanServiceEndpointMethod: methodIdent,
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packages/package-id/plans/plan-id/services/service-id/endpoints/endpoint-id/methods/method-id/responseFilter").
			WithMethod("get").
			RequestingFields(MasheryResponseFilterFields).
			WillReturnJsonOf(expRV)
	}

	autoTestGet(t,
		methodIdent,
		expRV,
		mockVisitor,
		func(client Client) ClientBoolExchangeFunc[masherytypes.PackagePlanServiceEndpointMethodIdentifier, masherytypes.PackagePlanServiceEndpointMethodFilter] {
			return client.GetPackagePlanMethodFilter
		},
	)
}

func TestCreatePackagePlanMethodFilter(t *testing.T) {
	serviceEndpointIdent := masherytypes.PackagePlanServiceEndpointMethodFilterIdentifier{}
	serviceEndpointIdent.PackageId = "package-id"
	serviceEndpointIdent.PlanId = "plan-id"
	serviceEndpointIdent.ServiceId = "service-id"
	serviceEndpointIdent.EndpointId = "endpoint-id"
	serviceEndpointIdent.MethodId = "method-id"
	serviceEndpointIdent.FilterId = "filter-id"

	expBody := masherytypes.IdReferenced{IdRef: "filter-id"}
	expRvCreate := masherytypes.PackagePlanServiceEndpointMethodFilter{
		ResponseFilter: masherytypes.ResponseFilter{
			AddressableV3Object: masherytypes.AddressableV3Object{Id: "filter-id", Name: "filter-name"},
		},
		PackagePlanServiceEndpointMethod: serviceEndpointIdent.AsPackagePlanServiceEndpointMethodIdentifier(),
	}

	onTheWire := cloneWithModification(expRvCreate, func(t1 *masherytypes.PackagePlanServiceEndpointMethodFilter) {
		t1.PackagePlanServiceEndpointMethod = masherytypes.PackagePlanServiceEndpointMethodIdentifier{}
	})

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packages/package-id/plans/plan-id/services/service-id/endpoints/endpoint-id/methods/method-id/responseFilter").
			WithMethod("post").
			RequestingFields(MasheryResponseFilterFields).
			Matching(PayloadMatcher(expBody)).
			WillReturnJsonOf(onTheWire)
	}

	autoTestRootAsymmetricCreate(t,
		serviceEndpointIdent,
		expRvCreate,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.PackagePlanServiceEndpointMethodFilterIdentifier, masherytypes.PackagePlanServiceEndpointMethodFilter] {
			return client.CreatePackagePlanMethodFilter
		},
	)
}

func TestDeletePackagePlanMethodFilter(t *testing.T) {
	methodIdent := masherytypes.PackagePlanServiceEndpointMethodIdentifier{}
	methodIdent.PackageId = "package-id"
	methodIdent.PlanId = "plan-id"
	methodIdent.ServiceId = "service-id"
	methodIdent.EndpointId = "endpoint-id"
	methodIdent.MethodId = "method-id"

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packages/package-id/plans/plan-id/services/service-id/endpoints/endpoint-id/methods/method-id/responseFilter").
			WithMethod("delete").
			RequestingNoFields().
			WillReturnUnspecified()
	}

	autoTestDelete(t,
		methodIdent,
		mockVisitor,
		func(cl Client) BiConsumerCanErr[context.Context, masherytypes.PackagePlanServiceEndpointMethodIdentifier] {
			return cl.DeletePackagePlanMethodFilter
		},
	)
}
