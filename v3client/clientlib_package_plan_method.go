package v3client

import (
	"context"
	"fmt"
	"net/url"
)

const PackagePlanMethodAppCtx = "package plan method"

func packagePlanEndpointMethodsRoot(id MasheryPlanServiceEndpoint) string {
	return fmt.Sprintf("/packages/%s/plans/%s/services/%s/endpoints/%s/methods", id.PackageId, id.PlanId, id.ServiceId, id.EndpointId)
}

func packagePlanEndpointMethod(id MasheryPlanServiceEndpointMethod) string {
	return fmt.Sprintf("/packages/%s/plans/%s/services/%s/endpoints/%s/methods/%s", id.PackageId, id.PlanId, id.ServiceId, id.EndpointId, id.MethodId)
}

func (c *HttpTransport) ListPackagePlanMethods(ctx context.Context, id MasheryPlanServiceEndpoint) ([]MasheryMethod, error) {
	opCtx := FetchSpec{
		Pagination:     PerItem,
		Resource:       packagePlanEndpointMethodsRoot(id),
		Query:          nil,
		AppContext:     "package plan methods",
		ResponseParser: ParseMasheryMethodArray,
	}

	if d, err := c.fetchAll(ctx, opCtx); err != nil {
		return []MasheryMethod{}, nil
	} else {
		// Convert individual fetches into the array of elements
		var rv []MasheryMethod
		for _, raw := range d {
			ms, ok := raw.([]MasheryMethod)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}

// Retrieve the information about a pacakge plan method.
func (c *HttpTransport) GetPackagePlanMethod(ctx context.Context, id MasheryPlanServiceEndpointMethod) (*MasheryMethod, error) {
	rv, err := c.getObject(ctx, FetchSpec{
		Pagination: PerItem,
		Resource:   packagePlanEndpointMethod(id),
		Query: url.Values{
			"fields": {MasheryMethodsFieldsStr},
		},
		AppContext:     PackagePlanMethodAppCtx,
		ResponseParser: ParseMasheryMethod,
	})

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(MasheryMethod)
		return &retServ, nil
	}
}

// Create a new service cache
func (c *HttpTransport) CreatePackagePlanMethod(ctx context.Context, id MasheryPlanServiceEndpoint, upsert MasheryMethod) (*MasheryMethod, error) {
	rawResp, err := c.createObject(ctx, upsert, FetchSpec{
		Pagination: NotRequired,
		Resource:   packagePlanEndpointMethodsRoot(id),
		Query: url.Values{
			"fields": {MasheryMethodsFieldsStr},
		},
		AppContext:     PackagePlanMethodAppCtx,
		ResponseParser: ParseMasheryMethodArray,
	})

	if err == nil {
		rv, _ := rawResp.(MasheryMethod)
		return &rv, nil
	} else {
		return nil, err
	}
}

// Create a new service.
func (c *HttpTransport) DeletePackagePlanMethod(ctx context.Context, id MasheryPlanServiceEndpointMethod) error {
	return c.deleteObject(ctx, FetchSpec{
		Pagination: NotRequired,
		Resource:   packagePlanEndpointMethod(id),
		AppContext: PackagePlanMethodAppCtx,
	})
}
