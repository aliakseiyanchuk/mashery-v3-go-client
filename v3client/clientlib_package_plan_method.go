package v3client

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"net/url"
)

const PackagePlanMethodAppCtx = "package plan method"

func ListPackagePlanMethods(ctx context.Context, id masherytypes.PackagePlanServiceEndpointIdentifier, c *transport.HttpTransport) ([]masherytypes.PackagePlanServiceEndpointMethod, error) {
	builder := transport.ObjectListFetchSpecBuilder[masherytypes.PackagePlanServiceEndpointMethod]{}
	builder.
		WithValueFactory(func() []masherytypes.PackagePlanServiceEndpointMethod {
			return []masherytypes.PackagePlanServiceEndpointMethod{}
		}).
		WithResource("/packages/%s/plans/%s/services/%s/endpoints/%s/methods", id.PackageId, id.PlanId, id.ServiceId, id.EndpointId).
		WithPagination(transport.PerItem).
		WithAppContext("package plan methods")

	return transport.FetchAll[masherytypes.PackagePlanServiceEndpointMethod](ctx, builder.Build(), c)
}

// GetPackagePlanMethod Retrieve the information about a package plan method.
func GetPackagePlanMethod(ctx context.Context, id masherytypes.PackagePlanServiceEndpointMethodIdentifier, c *transport.HttpTransport) (masherytypes.PackagePlanServiceEndpointMethod, bool, error) {

	builder := transport.ObjectFetchSpecBuilder[masherytypes.PackagePlanServiceEndpointMethod]{}
	builder.
		WithValueFactory(func() masherytypes.PackagePlanServiceEndpointMethod {
			return masherytypes.PackagePlanServiceEndpointMethod{}
		}).
		WithResource("/packages/%s/plans/%s/services/%s/endpoints/%s/methods/%s", id.PackageId, id.PlanId, id.ServiceId, id.EndpointId, id.MethodId).
		WithPagination(transport.PerPage).
		WithQuery(url.Values{
			"fields": {MasheryMethodsFieldsStr},
		}).
		WithAppContext(PackagePlanMethodAppCtx).
		WithReturn404AsNil(true)

	return transport.GetObject(ctx, builder.Build(), c)
}

// CreatePackagePlanServiceEndpointMethod Create a new service cache
func CreatePackagePlanServiceEndpointMethod(ctx context.Context, ident masherytypes.PackagePlanServiceEndpointMethodIdentifier, c *transport.HttpTransport) (masherytypes.PackagePlanServiceEndpointMethod, error) {
	upsert := masherytypes.ServiceEndpointMethod{
		BaseMethod: masherytypes.BaseMethod{
			AddressableV3Object: masherytypes.AddressableV3Object{
				Id: ident.MethodId,
			},
		},
	}

	builder := transport.ObjectExchangeSpecBuilder[masherytypes.ServiceEndpointMethod, masherytypes.PackagePlanServiceEndpointMethod]{}
	builder.
		WithBody(upsert).
		WithValueFactory(func() masherytypes.PackagePlanServiceEndpointMethod {
			return masherytypes.PackagePlanServiceEndpointMethod{}
		}).
		WithResource("/packages/%s/plans/%s/services/%s/endpoints/%s/methods", ident.PackageId, ident.PlanId, ident.ServiceId, ident.EndpointId).
		WithQuery(url.Values{
			"fields": {MasheryMethodsFieldsStr},
		}).
		WithAppContext(PackagePlanMethodAppCtx).
		WithReturn404AsNil(false)

	if rv, err := transport.ExchangeObject(ctx, builder.Build(), "post", c); err != nil {
		return masherytypes.PackagePlanServiceEndpointMethod{}, err
	} else {
		rv.PackagePlanServiceEndpoint = masherytypes.PackagePlanServiceEndpointIdentifier{
			PackagePlanIdentifier:     ident.PackagePlanIdentifier,
			ServiceEndpointIdentifier: ident.ServiceEndpointMethodIdentifier.ServiceEndpointIdentifier,
		}
		return rv, err
	}
}

// DeletePackagePlanMethod Create a new service.
func DeletePackagePlanMethod(ctx context.Context, id masherytypes.PackagePlanServiceEndpointMethodIdentifier, c *transport.HttpTransport) error {

	builder := transport.ObjectFetchSpecBuilder[masherytypes.PackagePlanServiceEndpointMethod]{}
	builder.
		WithValueFactory(func() masherytypes.PackagePlanServiceEndpointMethod {
			return masherytypes.PackagePlanServiceEndpointMethod{}
		}).
		WithResource("/packages/%s/plans/%s/services/%s/endpoints/%s/methods/%s", id.PackageId, id.PlanId, id.ServiceId, id.EndpointId, id.MethodId).
		WithAppContext(PackagePlanMethodAppCtx).
		WithIgnoreResponse(true)

	return transport.DeleteObject(ctx, builder.Build(), c)
}
