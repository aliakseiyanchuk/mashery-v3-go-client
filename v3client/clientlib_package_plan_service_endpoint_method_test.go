package v3client

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"testing"
)

func TestListPackagePlanMethods(t *testing.T) {
	serviceEndpointIdent := masherytypes.PackagePlanServiceEndpointIdentifier{}
	serviceEndpointIdent.PackageId = "package-id"
	serviceEndpointIdent.PlanId = "plan-id"
	serviceEndpointIdent.ServiceId = "service-id"
	serviceEndpointIdent.EndpointId = "endpoint-id"

	expRV := []masherytypes.PackagePlanServiceEndpointMethod{
		{
			BaseMethod: masherytypes.BaseMethod{
				AddressableV3Object: masherytypes.AddressableV3Object{Id: "method-id", Name: "method-name"},
			},
			PackagePlanServiceEndpoint: serviceEndpointIdent,
		},
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packages/package-id/plans/plan-id/services/service-id/endpoints/endpoint-id/methods").
			WithMethod("get").
			RequestingNoFields().
			WillReturnJsonOf(expRV)
	}

	autoTestFetchAll(t,
		serviceEndpointIdent,
		expRV,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.PackagePlanServiceEndpointIdentifier, []masherytypes.PackagePlanServiceEndpointMethod] {
			return client.ListPackagePlanMethods
		},
	)
}

func TestGetPackagePlanMethod(t *testing.T) {
	serviceIdent := masherytypes.PackagePlanServiceEndpointMethodIdentifier{}
	serviceIdent.PackageId = "package-id"
	serviceIdent.PlanId = "plan-id"
	serviceIdent.ServiceId = "service-id"
	serviceIdent.EndpointId = "endpoint-id"
	serviceIdent.MethodId = "method-id"

	expRV := masherytypes.PackagePlanServiceEndpointMethod{
		BaseMethod: masherytypes.BaseMethod{
			AddressableV3Object: masherytypes.AddressableV3Object{Id: "method-id", Name: "method-name"},
		},
		PackagePlanServiceEndpoint: serviceIdent.GetPackagePlanServiceEndpointIdentifier(),
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packages/package-id/plans/plan-id/services/service-id/endpoints/endpoint-id/methods/method-id").
			WithMethod("get").
			RequestingFields(MasheryMethodsFields).
			WillReturnJsonOf(expRV)
	}

	autoTestGet(t,
		serviceIdent,
		expRV,
		mockVisitor,
		func(client Client) ClientBoolExchangeFunc[masherytypes.PackagePlanServiceEndpointMethodIdentifier, masherytypes.PackagePlanServiceEndpointMethod] {
			return client.GetPackagePlanMethod
		},
	)
}

func TestCreatePackagePlanMethod(t *testing.T) {
	serviceEndpointIdent := masherytypes.PackagePlanServiceEndpointMethodIdentifier{}
	serviceEndpointIdent.PackageId = "package-id"
	serviceEndpointIdent.PlanId = "plan-id"
	serviceEndpointIdent.ServiceId = "service-id"
	serviceEndpointIdent.EndpointId = "endpoint-id"
	serviceEndpointIdent.MethodId = "method-id"

	expBody := masherytypes.IdReferenced{IdRef: "method-id"}
	expRvCreate := masherytypes.PackagePlanServiceEndpointMethod{
		BaseMethod: masherytypes.BaseMethod{
			AddressableV3Object: masherytypes.AddressableV3Object{Id: "method-id", Name: "method-name"},
		},
		PackagePlanServiceEndpoint: serviceEndpointIdent.GetPackagePlanServiceEndpointIdentifier(),
	}

	onTheWire := cloneWithModification(expRvCreate, func(t1 *masherytypes.PackagePlanServiceEndpointMethod) {
		t1.PackagePlanServiceEndpoint = masherytypes.PackagePlanServiceEndpointIdentifier{}
	})

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packages/package-id/plans/plan-id/services/service-id/endpoints/endpoint-id/methods").
			WithMethod("post").
			RequestingFields(MasheryMethodsFields).
			Matching(PayloadMatcher(expBody)).
			WillReturnJsonOf(onTheWire)
	}

	autoTestRootAsymmetricCreate(t,
		serviceEndpointIdent,
		expRvCreate,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.PackagePlanServiceEndpointMethodIdentifier, masherytypes.PackagePlanServiceEndpointMethod] {
			return client.CreatePackagePlanMethod
		},
	)
}

func TestDeletePackagePlanMethod(t *testing.T) {
	serviceIdent := masherytypes.PackagePlanServiceEndpointMethodIdentifier{}
	serviceIdent.PackageId = "package-id"
	serviceIdent.PlanId = "plan-id"
	serviceIdent.ServiceId = "service-id"
	serviceIdent.EndpointId = "endpoint-id"
	serviceIdent.MethodId = "method-id"

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packages/package-id/plans/plan-id/services/service-id/endpoints/endpoint-id/methods/method-id").
			WithMethod("delete").
			RequestingNoFields().
			WillReturnUnspecified()
	}

	autoTestDelete(t,
		serviceIdent,
		mockVisitor,
		func(cl Client) BiConsumerCanErr[context.Context, masherytypes.PackagePlanServiceEndpointMethodIdentifier] {
			return cl.DeletePackagePlanMethod
		},
	)
}
