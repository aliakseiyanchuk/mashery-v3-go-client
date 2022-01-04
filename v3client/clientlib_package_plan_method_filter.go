package v3client

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"net/url"
)

const PackagePlanMethodFilterAppCtx = "package plan method filter"

func packagePlanEndpointMethodFilter(id masherytypes.MasheryPlanServiceEndpointMethod) string {
	return fmt.Sprintf("/packages/%s/plans/%s/services/%s/endpoints/%s/methods/%s/responseFilter", id.PackageId, id.PlanId, id.ServiceId, id.EndpointId, id.MethodId)
}

// GetPackagePlanMethodFilter Retrieve the information about a package plan method.
func GetPackagePlanMethodFilter(ctx context.Context, id masherytypes.MasheryPlanServiceEndpointMethod, c *transport.V3Transport) (*masherytypes.MasheryResponseFilter, error) {
	rv, err := c.GetObject(ctx, transport.FetchSpec{
		Pagination: transport.PerItem,
		Resource:   packagePlanEndpointMethodFilter(id),
		Query: url.Values{
			"fields": {MasheryResponseFilterFieldsStr},
		},
		AppContext:     PackagePlanMethodFilterAppCtx,
		ResponseParser: masherytypes.ParseMasheryResponseFilter,
	})

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(masherytypes.MasheryResponseFilter)
		return &retServ, nil
	}
}

// CreatePackagePlanMethodFilter Create a new service cache
func CreatePackagePlanMethodFilter(ctx context.Context, id masherytypes.MasheryPlanServiceEndpointMethod, ref masherytypes.MasheryServiceMethodFilter, c *transport.V3Transport) (*masherytypes.MasheryResponseFilter, error) {
	upsert := masherytypes.IdReferenced{IdRef: ref.MethodId}

	rawResp, err := c.CreateObject(ctx, upsert, transport.FetchSpec{
		Pagination: transport.NotRequired,
		Resource:   packagePlanEndpointMethodFilter(id),
		Query: url.Values{
			"fields": {MasheryResponseFilterFieldsStr},
		},
		AppContext:     PackagePlanMethodFilterAppCtx,
		ResponseParser: masherytypes.ParseMasheryResponseFilter,
	})

	if err == nil {
		rv, _ := rawResp.(masherytypes.MasheryResponseFilter)
		return &rv, nil
	} else {
		return nil, err
	}
}

// DeletePackagePlanMethodFilter Create a new service.
func DeletePackagePlanMethodFilter(ctx context.Context, id masherytypes.MasheryPlanServiceEndpointMethod, c *transport.V3Transport) error {
	return c.DeleteObject(ctx, transport.FetchSpec{
		Pagination: transport.NotRequired,
		Resource:   packagePlanEndpointMethodFilter(id),
		AppContext: PackagePlanMethodFilterAppCtx,
	})
}
