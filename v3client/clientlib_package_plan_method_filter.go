package v3client

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"net/url"
)

var packagePlanServiceEndpointMethodFilterCRUDDecorator *GenericCRUDDecorator[masherytypes.PackagePlanServiceEndpointMethodIdentifier,
	masherytypes.PackagePlanServiceEndpointMethodIdentifier,
	masherytypes.PackagePlanServiceEndpointMethodFilter]

var packagePlanServiceEndpointMethodFilterCRUD *GenericCRUD[masherytypes.PackagePlanServiceEndpointMethodIdentifier,
	masherytypes.PackagePlanServiceEndpointMethodIdentifier,
	masherytypes.PackagePlanServiceEndpointMethodFilter,
]

func init() {
	packagePlanServiceEndpointMethodFilterCRUDDecorator = &GenericCRUDDecorator[masherytypes.PackagePlanServiceEndpointMethodIdentifier,
		masherytypes.PackagePlanServiceEndpointMethodIdentifier,
		masherytypes.PackagePlanServiceEndpointMethodFilter]{

		ValueSupplier: func() masherytypes.PackagePlanServiceEndpointMethodFilter {
			return masherytypes.PackagePlanServiceEndpointMethodFilter{}
		},
		ValueArraySupplier: func() []masherytypes.PackagePlanServiceEndpointMethodFilter {
			return []masherytypes.PackagePlanServiceEndpointMethodFilter{}
		},

		AcceptIdentFrom: func(t1 masherytypes.PackagePlanServiceEndpointMethodFilter, t2 *masherytypes.PackagePlanServiceEndpointMethodFilter) {
			t2.PackagePlanServiceEndpointMethod = t1.PackagePlanServiceEndpointMethod
		},
		AcceptObjectIdent: func(t1 masherytypes.PackagePlanServiceEndpointMethodIdentifier, t2 *masherytypes.PackagePlanServiceEndpointMethodFilter) {
			t2.PackagePlanServiceEndpointMethod = t1
		},
		AcceptParentIdent: func(t1 masherytypes.PackagePlanServiceEndpointMethodIdentifier, t2 *masherytypes.PackagePlanServiceEndpointMethodFilter) {
			t2.PackagePlanServiceEndpointMethod = t1
		},

		ResourceFor: func(id masherytypes.PackagePlanServiceEndpointMethodIdentifier) (string, error) {
			return fmt.Sprintf(
					"/packages/%s/plans/%s/services/%s/endpoints/%s/methods/%s/responseFilter",
					id.PackageId,
					id.PlanId,
					id.ServiceId,
					id.EndpointId,
					id.MethodId),
				nil
		},

		ResourceForUpsert: func(id masherytypes.PackagePlanServiceEndpointMethodFilter) (string, error) {
			return fmt.Sprintf(
					"/packages/%s/plans/%s/services/%s/endpoints/%s/methods/%s/responseFilter",
					id.PackagePlanServiceEndpointMethod.PackageId,
					id.PackagePlanServiceEndpointMethod.PlanId,
					id.PackagePlanServiceEndpointMethod.ServiceId,
					id.PackagePlanServiceEndpointMethod.EndpointId,
					id.Id),
				nil
		},

		ResourceForParent: func(id masherytypes.PackagePlanServiceEndpointMethodIdentifier) (string, error) {
			return fmt.Sprintf("/packages/%s/plans/%s/services/%s/endpoints/%s/methods/%s/responseFilter",
					id.PackageId,
					id.PlanId,
					id.ServiceId,
					id.EndpointId,
					id.MethodId),
				nil
		},

		DefaultFields: MasheryResponseFilterFields,
		Pagination:    transport.NotRequired,
	}
	packagePlanServiceEndpointMethodFilterCRUD = NewCRUD[masherytypes.PackagePlanServiceEndpointMethodIdentifier,
		masherytypes.PackagePlanServiceEndpointMethodIdentifier,
		masherytypes.PackagePlanServiceEndpointMethodFilter,
	]("package plan method filter", packagePlanServiceEndpointMethodFilterCRUDDecorator)
}

const PackagePlanMethodFilterAppCtx = "package plan method filter"

// CreatePackagePlanMethodFilter Create a new service cache
func CreatePackagePlanMethodFilter(ctx context.Context,
	ident masherytypes.PackagePlanServiceEndpointMethodFilterIdentifier,
	c *transport.HttpTransport) (masherytypes.PackagePlanServiceEndpointMethodFilter, error) {

	upsert := masherytypes.IdReferenced{IdRef: ident.FilterId}

	builder := transport.ObjectExchangeSpecBuilder[masherytypes.IdReferenced, masherytypes.PackagePlanServiceEndpointMethodFilter]{}
	builder.
		WithBody(upsert).
		WithValueFactory(func() masherytypes.PackagePlanServiceEndpointMethodFilter {
			return masherytypes.PackagePlanServiceEndpointMethodFilter{}
		}).
		WithResource("/packages/%s/plans/%s/services/%s/endpoints/%s/methods/%s/responseFilter", ident.PackageId, ident.PlanId, ident.ServiceId, ident.EndpointId, ident.MethodId).
		WithQuery(url.Values{
			"fields": {MasheryResponseFilterFieldsStr},
		}).
		WithAppContext(PackagePlanMethodFilterAppCtx)

	if rv, err := transport.ExchangeObject(ctx, builder.Build(), "post", c); err != nil {
		return masherytypes.PackagePlanServiceEndpointMethodFilter{}, err
	} else {
		rv.PackagePlanServiceEndpointMethod = ident.AsPackagePlanServiceEndpointMethodIdentifier()

		return rv, nil
	}
}
