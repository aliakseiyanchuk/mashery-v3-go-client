package v3client

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
)

var packagePlanCRUDDecorator *GenericCRUDDecorator[masherytypes.PackageIdentifier, masherytypes.PackagePlanIdentifier, masherytypes.Plan]
var packagePlanCRDU *GenericCRUD[masherytypes.PackageIdentifier, masherytypes.PackagePlanIdentifier, masherytypes.Plan]

func init() {
	packagePlanCRUDDecorator = &GenericCRUDDecorator[masherytypes.PackageIdentifier, masherytypes.PackagePlanIdentifier, masherytypes.Plan]{
		ValueSupplier:      func() masherytypes.Plan { return masherytypes.Plan{} },
		ValueArraySupplier: func() []masherytypes.Plan { return []masherytypes.Plan{} },

		AcceptParentIdent: func(t1 masherytypes.PackageIdentifier, t2 *masherytypes.Plan) {
			t2.ParentPackageId = t1
		},
		AcceptObjectIdent: func(t1 masherytypes.PackagePlanIdentifier, t2 *masherytypes.Plan) {
			t2.ParentPackageId = t1.PackageIdentifier
		},
		AcceptIdentFrom: func(t1 masherytypes.Plan, t2 *masherytypes.Plan) {
			t2.ParentPackageId = t1.ParentPackageId
		},

		ResourceFor: func(ident masherytypes.PackagePlanIdentifier) (string, error) {
			return fmt.Sprintf("/packages/%s/plans/%s", ident.PackageId, ident.PlanId), nil
		},
		ResourceForUpsert: func(t masherytypes.Plan) (string, error) {
			if len(t.Id) > 0 {
				return fmt.Sprintf("/packages/%s/plans/%s", t.ParentPackageId.PackageId, t.Id), nil
			}
			return "", errors.New("insufficient identification")
		},
		ResourceForParent: func(ident masherytypes.PackageIdentifier) (string, error) {
			return fmt.Sprintf("/packages/%s/plans", ident.PackageId), nil
		},
		DefaultFields: MasheryPlanFields,
		Pagination:    transport.PerPage,
	}
	packagePlanCRDU = NewCRUD[masherytypes.PackageIdentifier, masherytypes.PackagePlanIdentifier, masherytypes.Plan]("package plan", packagePlanCRUDDecorator)
}

func addressableV3ObjectFactory() masherytypes.AddressableV3Object {
	return masherytypes.AddressableV3Object{}
}
func addressableV3ObjectArrayFactory() []masherytypes.AddressableV3Object {
	return []masherytypes.AddressableV3Object{}
}

func CreatePlanService(ctx context.Context, planService masherytypes.PackagePlanServiceIdentifier, c *transport.HttpTransport) (masherytypes.AddressableV3Object, error) {
	ref := masherytypes.IdReferenced{IdRef: planService.ServiceId}

	builder := transport.ObjectExchangeSpecBuilder[masherytypes.IdReferenced, masherytypes.AddressableV3Object]{}
	builder.
		WithBody(ref).
		WithValueFactory(addressableV3ObjectFactory).
		WithResource("/packages/%s/plans/%s/services", planService.PackageId, planService.PlanId).
		WithAppContext("plan service")

	return transport.ExchangeObject(ctx, builder.Build(), "post", c)
}

func CheckPlanServiceExists(ctx context.Context, planService masherytypes.PackagePlanServiceIdentifier, c *transport.HttpTransport) (bool, error) {
	builder := transport.ObjectListFetchSpecBuilder[masherytypes.AddressableV3Object]{}
	builder.
		WithValueFactory(addressableV3ObjectArrayFactory).
		WithResource("/packages/%s/plans/%s/services/%s", planService.PackageId, planService.PlanId, planService.ServiceId).
		WithAppContext("check plan service exists").
		WithReturn404AsNil(true)

	return transport.Exists(ctx, builder.Build().AsObjectFetchSpec(), c)
}

func DeletePlanService(ctx context.Context, planService masherytypes.PackagePlanServiceIdentifier, c *transport.HttpTransport) error {
	builder := transport.ObjectFetchSpecBuilder[masherytypes.AddressableV3Object]{}
	builder.
		WithValueFactory(addressableV3ObjectFactory).
		WithResource("/packages/%s/plans/%s/services/%s", planService.PackageId, planService.PlanId, planService.ServiceId).
		WithAppContext("delete plan service").
		WithIgnoreResponse(true)

	return transport.DeleteObject(ctx, builder.Build(), c)
}

func CreatePlanEndpoint(ctx context.Context, planEndp masherytypes.PackagePlanServiceEndpointIdentifier, c *transport.HttpTransport) (masherytypes.AddressableV3Object, error) {
	ref := masherytypes.IdReferenced{IdRef: planEndp.EndpointId}

	builder := transport.ObjectExchangeSpecBuilder[masherytypes.IdReferenced, masherytypes.AddressableV3Object]{}
	builder.
		WithBody(ref).
		WithValueFactory(addressableV3ObjectFactory).
		WithResource("/packages/%s/plans/%s/services/%s/endpoints", planEndp.PackageId, planEndp.PlanId, planEndp.ServiceId).
		WithAppContext("create plan endpoint")

	return transport.ExchangeObject(ctx, builder.Build(), "post", c)
}

func CheckPlanEndpointExists(ctx context.Context, planEndp masherytypes.PackagePlanServiceEndpointIdentifier, c *transport.HttpTransport) (bool, error) {
	builder := transport.ObjectFetchSpecBuilder[masherytypes.AddressableV3Object]{}
	builder.
		WithValueFactory(addressableV3ObjectFactory).
		WithResource("/packages/%s/plans/%s/services/%s/endpoints/%s", planEndp.PackageId, planEndp.PlanId, planEndp.ServiceId, planEndp.EndpointId).
		WithAppContext("check plan endpoint exists").
		WithReturn404AsNil(true)

	return transport.Exists(ctx, builder.Build(), c)
}

func DeletePlanEndpoint(ctx context.Context, planEndp masherytypes.PackagePlanServiceEndpointIdentifier, c *transport.HttpTransport) error {
	builder := transport.ObjectFetchSpecBuilder[masherytypes.AddressableV3Object]{}
	builder.
		WithValueFactory(addressableV3ObjectFactory).
		WithResource("/packages/%s/plans/%s/services/%s/endpoints/%s", planEndp.PackageId, planEndp.PlanId, planEndp.ServiceId, planEndp.EndpointId).
		WithAppContext("delete plan endpoint").
		WithReturn404AsNil(false).
		WithIgnoreResponse(true)

	return transport.DeleteObject(ctx, builder.Build(), c)
}

func ListPlanEndpoints(ctx context.Context, planService masherytypes.PackagePlanServiceIdentifier, c *transport.HttpTransport) ([]masherytypes.AddressableV3Object, error) {
	builder := transport.ObjectListFetchSpecBuilder[masherytypes.AddressableV3Object]{}
	builder.
		WithValueFactory(addressableV3ObjectArrayFactory).
		WithPagination(transport.PerItem).
		WithResource("/packages/%s/plans/%s/services/%s/endpoints", planService.PackageId, planService.PlanId, planService.ServiceId).
		WithAppContext("list plan service endpoints")

	return transport.FetchAll(ctx, builder.Build(), c)
}

func CountPlanEndpoints(ctx context.Context, planService masherytypes.PackagePlanServiceIdentifier, c *transport.HttpTransport) (int64, error) {
	builder := transport.ObjectListFetchSpecBuilder[masherytypes.PackagePlanServiceIdentifier]{}
	builder.
		WithValueFactory(func() []masherytypes.PackagePlanServiceIdentifier {
			return []masherytypes.PackagePlanServiceIdentifier{}
		}).
		WithResource("/packages/%s/plans/%s/services/%s/endpoints", planService.PackageId, planService.PlanId, planService.ServiceId).
		WithAppContext("count plan service endpoints")

	return transport.Count(ctx, builder.Build(), c)
}

func CountPlanService(ctx context.Context, ident masherytypes.PackagePlanIdentifier, c *transport.HttpTransport) (int64, error) {
	builder := transport.ObjectListFetchSpecBuilder[masherytypes.PackagePlanIdentifier]{}
	builder.
		WithValueFactory(func() []masherytypes.PackagePlanIdentifier {
			return []masherytypes.PackagePlanIdentifier{}
		}).
		WithResource("/packages/%s/plans/%s/services", ident.PackageId, ident.PlanId).
		WithAppContext("count plans services")

	return transport.Count(ctx, builder.Build(), c)
}

func ListPlanServices(ctx context.Context, ident masherytypes.PackagePlanIdentifier, c *transport.HttpTransport) ([]masherytypes.Service, error) {
	builder := transport.ObjectListFetchSpecBuilder[masherytypes.Service]{}
	builder.
		WithValueFactory(func() []masherytypes.Service {
			return []masherytypes.Service{}
		}).
		WithResource("/packages/%s/plans/%s/services", ident.PackageId, ident.PlanId).
		WithPagination(transport.PerPage).
		WithAppContext("list plan service")

	return transport.FetchAll(ctx, builder.Build(), c)
}
