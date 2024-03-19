package v3client

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"testing"
)

func TestListPlanEndpoints(t *testing.T) {
	serviceIdent := masherytypes.PackagePlanServiceIdentifier{}
	serviceIdent.PackageId = "package-id"
	serviceIdent.PlanId = "plan-id"
	serviceIdent.ServiceId = "service-id"

	expRV := []masherytypes.AddressableV3Object{
		{Id: "endpoint-d", Name: "endpoint-name"},
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packages/package-id/plans/plan-id/services/service-id/endpoints").
			WithMethod("get").
			RequestingNoFields().
			WillReturnJsonOf(expRV)
	}

	autoTestFetchAll(t,
		serviceIdent,
		expRV,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.PackagePlanServiceIdentifier, []masherytypes.AddressableV3Object] {
			return client.ListPlanEndpoints
		},
	)
}

func TestCreatePlanEndpoint(t *testing.T) {
	serviceIdent := masherytypes.PackagePlanServiceEndpointIdentifier{}
	serviceIdent.PackageId = "package-id"
	serviceIdent.PlanId = "plan-id"
	serviceIdent.ServiceId = "service-id"
	serviceIdent.EndpointId = "endpoint-id"

	expBody := masherytypes.IdReferenced{IdRef: "endpoint-id"}
	expRvCreate := masherytypes.AddressableV3Object{Id: "endpoint-id", Name: "endpoint-name"}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packages/package-id/plans/plan-id/services/service-id/endpoints").
			WithMethod("post").
			RequestingNoFields().
			Matching(PayloadMatcher(expBody)).
			WillReturnJsonOf(expRvCreate)
	}

	autoTestRootAsymmetricCreate(t,
		serviceIdent,
		expRvCreate,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.PackagePlanServiceEndpointIdentifier, masherytypes.AddressableV3Object] {
			return client.CreatePlanEndpoint
		},
	)
}

func TestCheckPlanEndpointExists(t *testing.T) {
	serviceIdent := masherytypes.PackagePlanServiceEndpointIdentifier{}
	serviceIdent.PackageId = "package-id"
	serviceIdent.PlanId = "plan-id"
	serviceIdent.ServiceId = "service-id"
	serviceIdent.EndpointId = "endpoint-id"

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packages/package-id/plans/plan-id/services/service-id/endpoints/endpoint-id").
			WithMethod("get").
			RequestingNoFields().
			WillReturnJsonOf(masherytypes.AddressableV3Object{})
	}

	autoTestExists(t,
		serviceIdent,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.PackagePlanServiceEndpointIdentifier, bool] {
			return client.CheckPlanEndpointExists
		},
	)
}

func TestDeletePlanEndpoint(t *testing.T) {
	serviceIdent := masherytypes.PackagePlanServiceEndpointIdentifier{}
	serviceIdent.PackageId = "package-id"
	serviceIdent.PlanId = "plan-id"
	serviceIdent.ServiceId = "service-id"
	serviceIdent.EndpointId = "endpoint-id"

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packages/package-id/plans/plan-id/services/service-id/endpoints/endpoint-id").
			WithMethod("delete").
			RequestingNoFields().
			WillReturnUnspecified()
	}

	autoTestDelete(t,
		serviceIdent,
		mockVisitor,
		func(cl Client) BiConsumerCanErr[context.Context, masherytypes.PackagePlanServiceEndpointIdentifier] {
			return cl.DeletePlanEndpoint
		},
	)
}
