package v3client

import (
	"context"
	"fmt"
	"net/url"
)

const PackagePlanMethodFilterAppCtx = "package plan method filter"

func packagePlanEndpointMethodFilter(id MasheryPlanServiceEndpointMethod) string {
	return fmt.Sprintf("/packages/%s/plans/%s/services/%s/endpoints/%s/methods/%s/responseFilter", id.PackageId, id.PlanId, id.ServiceId, id.EndpointId, id.MethodId)
}

// Retrieve the information about a pacakge plan method.
func (c *HttpTransport) GetPackagePlanMethodFilter(ctx context.Context, id MasheryPlanServiceEndpointMethod) (*MasheryResponseFilter, error) {
	rv, err := c.getObject(ctx, FetchSpec{
		Pagination: PerItem,
		Resource:   packagePlanEndpointMethodFilter(id),
		Query: url.Values{
			"fields": {MasheryResponseFilterFieldsStr},
		},
		AppContext:     PackagePlanMethodFilterAppCtx,
		ResponseParser: ParseMasheryResponseFilter,
	})

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(MasheryResponseFilter)
		return &retServ, nil
	}
}

// Create a new service cache
func (c *HttpTransport) CreatePackagePlanMethodFilter(ctx context.Context, id MasheryPlanServiceEndpointMethod, ref MasheryServiceMethodFilter) (*MasheryResponseFilter, error) {
	upsert := IdReferenced{IdRef: ref.MethodId}

	rawResp, err := c.createObject(ctx, upsert, FetchSpec{
		Pagination: NotRequired,
		Resource:   packagePlanEndpointMethodFilter(id),
		Query: url.Values{
			"fields": {MasheryResponseFilterFieldsStr},
		},
		AppContext:     PackagePlanMethodFilterAppCtx,
		ResponseParser: ParseMasheryResponseFilter,
	})

	if err == nil {
		rv, _ := rawResp.(MasheryResponseFilter)
		return &rv, nil
	} else {
		return nil, err
	}
}

// Create a new service.
func (c *HttpTransport) DeletePackagePlanMethodFilter(ctx context.Context, id MasheryPlanServiceEndpointMethod) error {
	return c.deleteObject(ctx, FetchSpec{
		Pagination: NotRequired,
		Resource:   packagePlanEndpointMethodFilter(id),
		AppContext: PackagePlanMethodFilterAppCtx,
	})
}
