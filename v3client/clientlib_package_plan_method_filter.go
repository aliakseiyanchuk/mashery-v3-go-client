package v3client

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"net/url"
)

const PackagePlanMethodFilterAppCtx = "package plan method filter"

func packagePlanEndpointMethodFilter(id masherytypes.PackagePlanServiceEndpointMethodIdentifier) string {
	return fmt.Sprintf("/packages/%s/plans/%s/services/%s/endpoints/%s/methods/%s/responseFilter", id.PackageId, id.PlanId, id.ServiceId, id.EndpointId, id.MethodId)
}

// GetPackagePlanMethodFilter Retrieve the information about a package plan method.
func GetPackagePlanMethodFilter(ctx context.Context, id masherytypes.PackagePlanServiceEndpointMethodIdentifier, c *transport.V3Transport) (*masherytypes.PackagePlanServiceEndpointMethodFilter, error) {
	rv, err := c.GetObject(ctx, transport.FetchSpec{
		Pagination: transport.PerItem,
		Resource:   packagePlanEndpointMethodFilter(id),
		Query: url.Values{
			"fields": {MasheryResponseFilterFieldsStr},
		},
		AppContext:     PackagePlanMethodFilterAppCtx,
		ResponseParser: masherytypes.ParsePackagePlanServiceEndpointMethodFilter,
	})

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(masherytypes.PackagePlanServiceEndpointMethodFilter)
		retServ.PackagePlanServiceEndpointMethod = masherytypes.PackagePlanServiceEndpointMethodIdentifier{
			ServiceEndpointMethodIdentifier: id.ServiceEndpointMethodIdentifier,
			PackagePlanIdentifier:           id.PackagePlanIdentifier,
		}
		return &retServ, nil
	}
}

// CreatePackagePlanMethodFilter Create a new service cache
func CreatePackagePlanMethodFilter(ctx context.Context,
	ident masherytypes.PackagePlanServiceEndpointMethodFilterIdentifier,
	c *transport.V3Transport) (*masherytypes.PackagePlanServiceEndpointMethodFilter, error) {

	upsert := masherytypes.IdReferenced{IdRef: ident.MethodId}

	rawResp, err := c.CreateObject(ctx, upsert, transport.FetchSpec{
		Pagination: transport.NotRequired,
		Resource: packagePlanEndpointMethodFilter(masherytypes.PackagePlanServiceEndpointMethodIdentifier{
			ServiceEndpointMethodIdentifier: ident.ServiceEndpointMethodIdentifier,
			PackagePlanIdentifier:           ident.PackagePlanIdentifier,
		}),
		Query: url.Values{
			"fields": {MasheryResponseFilterFieldsStr},
		},
		AppContext:     PackagePlanMethodFilterAppCtx,
		ResponseParser: masherytypes.ParsePackagePlanServiceEndpointMethodFilter,
	})

	if err == nil {
		rv, _ := rawResp.(masherytypes.PackagePlanServiceEndpointMethodFilter)
		rv.PackagePlanServiceEndpointMethod = masherytypes.PackagePlanServiceEndpointMethodIdentifier{
			ServiceEndpointMethodIdentifier: ident.ServiceEndpointMethodIdentifier,
			PackagePlanIdentifier:           ident.PackagePlanIdentifier,
		}

		return &rv, nil
	} else {
		return nil, err
	}
}

// DeletePackagePlanMethodFilter Create a new service.
func DeletePackagePlanMethodFilter(ctx context.Context, id masherytypes.PackagePlanServiceEndpointMethodIdentifier, c *transport.V3Transport) error {
	return c.DeleteObject(ctx, transport.FetchSpec{
		Pagination: transport.NotRequired,
		Resource:   packagePlanEndpointMethodFilter(id),
		AppContext: PackagePlanMethodFilterAppCtx,
	})
}
