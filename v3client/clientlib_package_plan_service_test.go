package v3client

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"testing"
)

func TestListPlanServices(t *testing.T) {
	serviceIdent := masherytypes.PackagePlanIdentifier{}
	serviceIdent.PackageId = "package-id"
	serviceIdent.PlanId = "plan-id"

	expRV := []masherytypes.Service{
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "service-id", Name: "service-name"}},
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packages/package-id/plans/plan-id/services").
			WithMethod("get").
			RequestingNoFields().
			WillReturnJsonOf(expRV)
	}

	autoTestFetchAll(t,
		serviceIdent,
		expRV,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.PackagePlanIdentifier, []masherytypes.Service] {
			return client.ListPlanServices
		},
	)
}

func TestCountPlanService(t *testing.T) {
	ident := masherytypes.PackagePlanIdentifier{
		PackageIdentifier: masherytypes.PackageIdentifier{PackageId: "package-id"},
		PlanId:            "plan-id",
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packages/package-id/plans/plan-id/services").
			WithMethod("get").
			RequestingNoFields().
			WillReturnJsonOf([]masherytypes.AddressableV3Object{{}}).
			WillReturnTotalCount(7)
	}

	autoTestCount(t,
		ident,
		7,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.PackagePlanIdentifier, int64] {
			return client.CountPlanService
		},
	)
}

func TestCreatePlanService(t *testing.T) {
	serviceIdent := masherytypes.PackagePlanServiceIdentifier{}
	serviceIdent.ServiceId = "service-id"
	serviceIdent.PackageId = "package-id"
	serviceIdent.PlanId = "plan-id"

	expBody := masherytypes.IdReferenced{IdRef: "service-id"}
	expRvCreate := masherytypes.AddressableV3Object{Id: "service-id", Name: "service-name"}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packages/package-id/plans/plan-id/services").
			WithMethod("post").
			RequestingNoFields().
			Matching(PayloadMatcher(expBody)).
			WillReturnJsonOf(expRvCreate)
	}

	autoTestRootAsymmetricCreate(t,
		serviceIdent,
		expRvCreate,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.PackagePlanServiceIdentifier, masherytypes.AddressableV3Object] {
			return client.CreatePlanService
		},
	)
}

func TestCheckPlanServiceExists(t *testing.T) {
	serviceIdent := masherytypes.PackagePlanServiceIdentifier{}
	serviceIdent.ServiceId = "service-id"
	serviceIdent.PackageId = "package-id"
	serviceIdent.PlanId = "plan-id"

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packages/package-id/plans/plan-id/services/service-id").
			WithMethod("get").
			RequestingNoFields().
			WillReturnJsonOf([]masherytypes.AddressableV3Object{{}}).
			WillReturnTotalCount(3)
	}

	autoTestExists(t,
		serviceIdent,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.PackagePlanServiceIdentifier, bool] {
			return client.CheckPlanServiceExists
		},
	)
}

func TestDeletePlanService(t *testing.T) {
	postIdent := masherytypes.PackagePlanServiceIdentifier{}
	postIdent.PackageId = "package-id"
	postIdent.PlanId = "plan-id"
	postIdent.ServiceId = "service-id"

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packages/package-id/plans/plan-id/services/service-id").
			WithMethod("delete").
			RequestingNoFields().
			WillReturnUnspecified()
	}

	autoTestDelete(t,
		postIdent,
		mockVisitor,
		func(cl Client) BiConsumerCanErr[context.Context, masherytypes.PackagePlanServiceIdentifier] {
			return cl.DeletePlanService
		},
	)
}
