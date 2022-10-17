package v3client

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"net/url"
)

const PackagePlanMethodAppCtx = "package plan method"

func packagePlanEndpointMethodsRoot(id masherytypes.PackagePlanServiceEndpointIdentifier) string {
	return fmt.Sprintf("/packages/%s/plans/%s/services/%s/endpoints/%s/methods", id.PackageId, id.PlanId, id.ServiceId, id.EndpointId)
}

func packagePlanEndpointMethod(id masherytypes.PackagePlanServiceEndpointMethodIdentifier) string {
	return fmt.Sprintf("/packages/%s/plans/%s/services/%s/endpoints/%s/methods/%s", id.PackageId, id.PlanId, id.ServiceId, id.EndpointId, id.MethodId)
}

func ListPackagePlanMethods(ctx context.Context, id masherytypes.PackagePlanServiceEndpointIdentifier, c *transport.V3Transport) ([]masherytypes.PackagePlanServiceEndpointMethod, error) {
	opCtx := transport.FetchSpec{
		Pagination:     transport.PerItem,
		Resource:       packagePlanEndpointMethodsRoot(id),
		Query:          nil,
		AppContext:     "package plan methods",
		ResponseParser: masherytypes.ParsePackagePlanServiceEndpointMethodFilterArray,
	}

	if d, err := c.FetchAll(ctx, opCtx); err != nil {
		return []masherytypes.PackagePlanServiceEndpointMethod{}, nil
	} else {
		// Convert individual fetches into the array of elements
		var rv []masherytypes.PackagePlanServiceEndpointMethod
		for _, raw := range d {
			ms, ok := raw.([]masherytypes.PackagePlanServiceEndpointMethod)
			if ok {
				rv = append(rv, ms...)
			}
		}

		for _, v := range rv {
			v.PackagePlanServiceEndpoint = masherytypes.PackagePlanServiceEndpointIdentifier{
				PackagePlanIdentifier:     id.PackagePlanIdentifier,
				ServiceEndpointIdentifier: id.ServiceEndpointIdentifier,
			}
		}

		return rv, nil
	}
}

// GetPackagePlanMethod Retrieve the information about a package plan method.
func GetPackagePlanMethod(ctx context.Context, id masherytypes.PackagePlanServiceEndpointMethodIdentifier, c *transport.V3Transport) (*masherytypes.PackagePlanServiceEndpointMethod, error) {
	rv, err := c.GetObject(ctx, transport.FetchSpec{
		Pagination: transport.PerItem,
		Resource:   packagePlanEndpointMethod(id),
		Query: url.Values{
			"fields": {MasheryMethodsFieldsStr},
		},
		AppContext:     PackagePlanMethodAppCtx,
		ResponseParser: masherytypes.ParsePacakgePlanServiceEndpointMethod,
	})

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(masherytypes.PackagePlanServiceEndpointMethod)
		retServ.PackagePlanServiceEndpoint = masherytypes.PackagePlanServiceEndpointIdentifier{
			PackagePlanIdentifier:     id.PackagePlanIdentifier,
			ServiceEndpointIdentifier: id.ServiceEndpointMethodIdentifier.ServiceEndpointIdentifier,
		}
		return &retServ, nil
	}
}

// CreatePackagePlanServiceEndpointMethod Create a new service cache
func CreatePackagePlanServiceEndpointMethod(ctx context.Context, ident masherytypes.PackagePlanServiceEndpointMethodIdentifier, c *transport.V3Transport) (*masherytypes.PackagePlanServiceEndpointMethod, error) {
	upsert := masherytypes.ServiceEndpointMethod{
		BaseMethod: masherytypes.BaseMethod{
			AddressableV3Object: masherytypes.AddressableV3Object{
				Id: ident.MethodId,
			},
		},
	}

	rawResp, err := c.CreateObject(ctx, upsert, transport.FetchSpec{
		Pagination: transport.NotRequired,
		Resource: packagePlanEndpointMethodsRoot(masherytypes.PackagePlanServiceEndpointIdentifier{
			PackagePlanIdentifier:     ident.PackagePlanIdentifier,
			ServiceEndpointIdentifier: ident.ServiceEndpointMethodIdentifier.ServiceEndpointIdentifier,
		}),
		Query: url.Values{
			"fields": {MasheryMethodsFieldsStr},
		},
		AppContext:     PackagePlanMethodAppCtx,
		ResponseParser: masherytypes.ParsePacakgePlanServiceEndpointMethod,
	})

	if err == nil {
		rv, _ := rawResp.(masherytypes.PackagePlanServiceEndpointMethod)
		rv.PackagePlanServiceEndpoint = masherytypes.PackagePlanServiceEndpointIdentifier{
			PackagePlanIdentifier:     ident.PackagePlanIdentifier,
			ServiceEndpointIdentifier: ident.ServiceEndpointMethodIdentifier.ServiceEndpointIdentifier,
		}
		return &rv, nil
	} else {
		return nil, err
	}
}

// DeletePackagePlanMethod Create a new service.
func DeletePackagePlanMethod(ctx context.Context, id masherytypes.PackagePlanServiceEndpointMethodIdentifier, c *transport.V3Transport) error {
	return c.DeleteObject(ctx, transport.FetchSpec{
		Pagination: transport.NotRequired,
		Resource:   packagePlanEndpointMethod(id),
		AppContext: PackagePlanMethodAppCtx,
	})
}
