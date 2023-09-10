package v3client

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"testing"
)

func TestGetPlan(t *testing.T) {
	packagePlanId := masherytypes.PackagePlanIdentifier{}
	packagePlanId.PackageId = "package-id"
	packagePlanId.PlanId = "plan-id"

	expRvPlan := masherytypes.Plan{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   "plan-id",
			Name: "plan-name",
		},
		AdminKeyProvisioningEnabled: false,
		ParentPackageId:             packagePlanId.PackageIdentifier,
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packages/package-id/plans/plan-id").
			WithMethod("get").
			RequestingFields(MasheryPlanFields).
			WillReturnJsonOf(expRvPlan)
	}

	autoTestGet(t,
		packagePlanId,
		expRvPlan,
		mockVisitor,
		func(cl Client) ClientBoolExchangeFunc[masherytypes.PackagePlanIdentifier, masherytypes.Plan] {
			return cl.GetPlan
		},
	)
}

func TestCreatePlan(t *testing.T) {
	packageId := masherytypes.PackageIdentifier{PackageId: "package-id"}

	post := masherytypes.Plan{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   "",
			Name: "plan-name",
		},
		AdminKeyProvisioningEnabled: false,
	}

	apiResponseJson := cloneWithModification(post, func(t1 *masherytypes.Plan) { t1.Id = "plan-id" })

	expRvPlanReturned := cloneWithModification(apiResponseJson, func(t1 *masherytypes.Plan) { t1.ParentPackageId = packageId })

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packages/package-id/plans").
			WithMethod("post").
			RequestingFields(MasheryPlanFields).
			Matching(PayloadMatcher(post)).
			WillReturnJsonOf(apiResponseJson)
	}

	autoTestCreate(t,
		packageId,
		post,
		expRvPlanReturned,
		mockVisitor,
		func(cl Client) ClientDualExchangeFunc[masherytypes.PackageIdentifier, masherytypes.Plan, masherytypes.Plan] {
			return cl.CreatePlan
		},
	)
}

func TestUpdatePlan(t *testing.T) {
	expRvUpdatePayload := masherytypes.Plan{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   "plan-id",
			Name: "plan-name",
		},
		AdminKeyProvisioningEnabled: true,
		ParentPackageId:             masherytypes.PackageIdentifier{PackageId: "package-id"},
	}

	onTheWire := cloneWithModification(expRvUpdatePayload, func(t1 *masherytypes.Plan) { t1.ParentPackageId = masherytypes.PackageIdentifier{} })

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packages/package-id/plans/plan-id").
			WithMethod("put").
			RequestingFields(MasheryPlanFields).
			Matching(PayloadMatcher(onTheWire)).
			WillReturnJsonOf(onTheWire)
	}

	autoTestUpdate(t,
		expRvUpdatePayload,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.Plan, masherytypes.Plan] {
			return client.UpdatePlan
		},
	)
}

func TestDeletePlan(t *testing.T) {
	postIdent := masherytypes.PackagePlanIdentifier{
		PackageIdentifier: masherytypes.PackageIdentifier{PackageId: "package-id"},
		PlanId:            "plan-id",
	}

	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packages/package-id/plans/plan-id").
			WithMethod("delete").
			RequestingNoFields().
			WillReturnJsonOf(postIdent)
	}

	autoTestDelete(t,
		postIdent,
		mockVisitor,
		func(cl Client) BiConsumerCanErr[context.Context, masherytypes.PackagePlanIdentifier] {
			return cl.DeletePlan
		},
	)
}

func TestCountPlans(t *testing.T) {
	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packages/package-id/plans").
			WithMethod("get").
			RequestingNoFields().
			WillReturnJsonOf([]masherytypes.Plan{{}}).
			WillReturnTotalCount(98)
	}

	autoTestCount(t,
		masherytypes.PackageIdentifier{PackageId: "package-id"},
		98,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.PackageIdentifier, int64] {
			return client.CountPlans
		},
	)
}

func TestListPlans(t *testing.T) {
	mockedResponse := []masherytypes.Plan{
		{AddressableV3Object: masherytypes.AddressableV3Object{Id: "plan-id"}},
	}
	var mockVisitor BuildVisitor = func(matcher *RequestMatcher) {
		matcher.
			ForRequestPath("/packages/package-id/plans").
			WithMethod("get").
			RequestingFields(MasheryPlanFields).
			WillReturnJsonOf(mockedResponse)
	}

	identCtx := masherytypes.PackageIdentifier{PackageId: "package-id"}
	expPop := cloneAllWithModification(mockedResponse,
		func(t1 *masherytypes.Plan) {
			t1.ParentPackageId = identCtx
		})

	autoTestFetchAll(t,
		identCtx,
		expPop,
		mockVisitor,
		func(client Client) ClientExchangeFunc[masherytypes.PackageIdentifier, []masherytypes.Plan] {
			return client.ListPlans
		},
	)
}
