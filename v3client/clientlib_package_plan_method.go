package v3client

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"net/url"
)

const PackagePlanMethodAppCtx = "package plan method"

func packagePlanEndpointMethodsRoot(id masherytypes.MasheryPlanServiceEndpoint) string {
	return fmt.Sprintf("/packages/%s/plans/%s/services/%s/endpoints/%s/methods", id.PackageId, id.PlanId, id.ServiceId, id.EndpointId)
}

func packagePlanEndpointMethod(id masherytypes.MasheryPlanServiceEndpointMethod) string {
	return fmt.Sprintf("/packages/%s/plans/%s/services/%s/endpoints/%s/methods/%s", id.PackageId, id.PlanId, id.ServiceId, id.EndpointId, id.MethodId)
}

func ListPackagePlanMethods(ctx context.Context, id masherytypes.MasheryPlanServiceEndpoint, c *transport.V3Transport) ([]masherytypes.MasheryMethod, error) {
	opCtx := transport.FetchSpec{
		Pagination:     transport.PerItem,
		Resource:       packagePlanEndpointMethodsRoot(id),
		Query:          nil,
		AppContext:     "package plan methods",
		ResponseParser: masherytypes.ParseMasheryMethodArray,
	}

	if d, err := c.FetchAll(ctx, opCtx); err != nil {
		return []masherytypes.MasheryMethod{}, nil
	} else {
		// Convert individual fetches into the array of elements
		var rv []masherytypes.MasheryMethod
		for _, raw := range d {
			ms, ok := raw.([]masherytypes.MasheryMethod)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}

// GetPackagePlanMethod Retrieve the information about a package plan method.
func GetPackagePlanMethod(ctx context.Context, id masherytypes.MasheryPlanServiceEndpointMethod, c *transport.V3Transport) (*masherytypes.MasheryMethod, error) {
	rv, err := c.GetObject(ctx, transport.FetchSpec{
		Pagination: transport.PerItem,
		Resource:   packagePlanEndpointMethod(id),
		Query: url.Values{
			"fields": {MasheryMethodsFieldsStr},
		},
		AppContext:     PackagePlanMethodAppCtx,
		ResponseParser: masherytypes.ParseMasheryMethod,
	})

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(masherytypes.MasheryMethod)
		return &retServ, nil
	}
}

// CreatePackagePlanMethod Create a new service cache
func CreatePackagePlanMethod(ctx context.Context, id masherytypes.MasheryPlanServiceEndpoint, upsert masherytypes.MasheryMethod, c *transport.V3Transport) (*masherytypes.MasheryMethod, error) {
	rawResp, err := c.CreateObject(ctx, upsert, transport.FetchSpec{
		Pagination: transport.NotRequired,
		Resource:   packagePlanEndpointMethodsRoot(id),
		Query: url.Values{
			"fields": {MasheryMethodsFieldsStr},
		},
		AppContext:     PackagePlanMethodAppCtx,
		ResponseParser: masherytypes.ParseMasheryMethodArray,
	})

	if err == nil {
		rv, _ := rawResp.(masherytypes.MasheryMethod)
		return &rv, nil
	} else {
		return nil, err
	}
}

// DeletePackagePlanMethod Create a new service.
func DeletePackagePlanMethod(ctx context.Context, id masherytypes.MasheryPlanServiceEndpointMethod, c *transport.V3Transport) error {
	return c.DeleteObject(ctx, transport.FetchSpec{
		Pagination: transport.NotRequired,
		Resource:   packagePlanEndpointMethod(id),
		AppContext: PackagePlanMethodAppCtx,
	})
}
