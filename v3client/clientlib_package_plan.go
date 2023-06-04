package v3client

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"net/url"
)

func CreatePlanService(ctx context.Context, planService masherytypes.PackagePlanServiceIdentifier, c *transport.V3Transport) (*masherytypes.AddressableV3Object, error) {
	ref := masherytypes.IdReferenced{IdRef: planService.ServiceId}

	rv, err := c.CreateObject(ctx, ref, transport.FetchSpec{
		Pagination:     transport.NotRequired,
		Resource:       fmt.Sprintf("/packages/%s/plans/%s/services", planService.PackageId, planService.PlanId),
		Query:          nil,
		AppContext:     "plan service",
		ResponseParser: masherytypes.ParseMasheryAddressableObject,
	})

	if err != nil {
		return nil, err
	} else {
		rvc := rv.(masherytypes.AddressableV3Object)
		return &rvc, nil
	}
}

func CheckPlanServiceExists(ctx context.Context, planService masherytypes.PackagePlanServiceIdentifier, c *transport.V3Transport) (bool, error) {
	rv, err := c.GetObject(ctx, transport.FetchSpec{
		Pagination:     transport.NotRequired,
		Resource:       fmt.Sprintf("/packages/%s/plans/%s/services/%s", planService.PackageId, planService.PlanId, planService.ServiceId),
		Query:          nil,
		AppContext:     "check plan service exists",
		ResponseParser: masherytypes.ParseMasheryAddressableObject,
	})

	if err != nil {
		return false, err
	} else {
		return rv != nil, nil
	}
}

func DeletePlanService(ctx context.Context, planService masherytypes.PackagePlanServiceIdentifier, c *transport.V3Transport) error {
	return c.DeleteObject(ctx, transport.FetchSpec{
		Resource:   fmt.Sprintf("/packages/%s/plans/%s/services/%s", planService.PackageId, planService.PlanId, planService.ServiceId),
		AppContext: "plan service",
	})
}

func CreatePlanEndpoint(ctx context.Context, planEndp masherytypes.PackagePlanServiceEndpointIdentifier, c *transport.V3Transport) (*masherytypes.AddressableV3Object, error) {
	ref := masherytypes.IdReferenced{IdRef: planEndp.EndpointId}
	rv, err := c.CreateObject(ctx, ref, transport.FetchSpec{
		Resource:       fmt.Sprintf("/packages/%s/plans/%s/services/%s/endpoints", planEndp.PackageId, planEndp.PlanId, planEndp.ServiceId),
		AppContext:     "create plan endpoint",
		ResponseParser: masherytypes.ParseMasheryAddressableObject,
	})

	if err != nil {
		return nil, err
	} else {
		rvc := rv.(masherytypes.AddressableV3Object)
		return &rvc, nil
	}
}

func CheckPlanEndpointExists(ctx context.Context, planEndp masherytypes.PackagePlanServiceEndpointIdentifier, c *transport.V3Transport) (bool, error) {
	rv, err := c.GetObject(ctx, transport.FetchSpec{
		Resource:       fmt.Sprintf("/packages/%s/plans/%s/services/%s/endpoints/%s", planEndp.PackageId, planEndp.PlanId, planEndp.ServiceId, planEndp.EndpointId),
		AppContext:     "check plan endpoint exists",
		ResponseParser: masherytypes.ParseMasheryAddressableObject,
	})

	if err != nil {
		return false, err
	} else {
		return rv != nil, nil
	}
}

func DeletePlanEndpoint(ctx context.Context, planEndp masherytypes.PackagePlanServiceEndpointIdentifier, c *transport.V3Transport) error {
	return c.DeleteObject(ctx, transport.FetchSpec{
		Resource:   fmt.Sprintf("/packages/%s/plans/%s/services/%s/endpoints/%s", planEndp.PackageId, planEndp.PlanId, planEndp.ServiceId, planEndp.EndpointId),
		AppContext: "plan endpoint",
	})
}

func ListPlanEndpoints(ctx context.Context, planService masherytypes.PackagePlanServiceIdentifier, c *transport.V3Transport) ([]masherytypes.AddressableV3Object, error) {
	rv, err := c.FetchAll(ctx, transport.FetchSpec{
		Pagination:     transport.PerItem,
		Resource:       fmt.Sprintf("/packages/%s/plans/%s/services/%s/endpoints", planService.PackageId, planService.PlanId, planService.ServiceId),
		Query:          nil,
		AppContext:     "list plan service endpoints",
		ResponseParser: masherytypes.ParseMasheryAddressableObjectsArray,
	})

	if err != nil {
		return []masherytypes.AddressableV3Object{}, err
	} else {
		rvCast := make([]masherytypes.AddressableV3Object, len(rv))
		for i, v := range rv {
			rvCast[i] = v.(masherytypes.AddressableV3Object)
		}

		return rvCast, nil
	}
}

func CountPlanEndpoints(ctx context.Context, planService masherytypes.PackagePlanServiceIdentifier, c *transport.V3Transport) (int64, error) {
	opCtx := transport.FetchSpec{
		Pagination: transport.NotRequired,
		Resource:   fmt.Sprintf("/packages/%s/plans/%s/services/%s/endpoints", planService.PackageId, planService.PlanId, planService.ServiceId),
		AppContext: "count plan service endpoints",
	}

	return c.Count(ctx, opCtx)
}

func CountPlanService(ctx context.Context, ident masherytypes.PackagePlanIdentifier, c *transport.V3Transport) (int64, error) {
	opCtx := transport.FetchSpec{
		Pagination: transport.NotRequired,
		Resource:   fmt.Sprintf("/packages/%s/plans/%s/services", ident.PackageId, ident.PlanId),
		AppContext: "count plans services",
	}

	return c.Count(ctx, opCtx)
}

func GetPlan(ctx context.Context, ident masherytypes.PackagePlanIdentifier, c *transport.V3Transport) (*masherytypes.Plan, error) {
	rv, err := c.GetObject(ctx, transport.FetchSpec{
		Resource: fmt.Sprintf("/packages/%s/plans/%s", ident.PackageId, ident.PlanId),
		Query: url.Values{
			"fields": {MasheryPlanFieldsStr},
		},
		AppContext:     "plan",
		ResponseParser: masherytypes.ParseMasheryPlan,
	})

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(masherytypes.Plan)
		retServ.ParentPackageId = ident.PackageIdentifier

		return &retServ, nil
	}
}

// CreatePlan Create a new service.
func CreatePlan(ctx context.Context, packageId masherytypes.PackageIdentifier, plan masherytypes.Plan, c *transport.V3Transport) (*masherytypes.Plan, error) {
	rawResp, err := c.CreateObject(ctx, plan, transport.FetchSpec{
		Resource:   fmt.Sprintf("/packages/%s/plans", packageId.PackageId),
		AppContext: "plan",
		Query: url.Values{
			"fields": {MasheryPlanFieldsStr},
		},
		ResponseParser: masherytypes.ParseMasheryPlan,
	})

	if err == nil {
		rv, _ := rawResp.(masherytypes.Plan)
		rv.ParentPackageId = packageId
		return &rv, nil
	} else {
		return nil, err
	}
}

// UpdatePlan Create a new service.
func UpdatePlan(ctx context.Context, plan masherytypes.Plan, c *transport.V3Transport) (*masherytypes.Plan, error) {
	opContext := transport.FetchSpec{
		Resource:   fmt.Sprintf("/packages/%s/plans/%s", plan.ParentPackageId.PackageId, plan.Id),
		AppContext: "plan",
		Query: url.Values{
			"fields": {MasheryPlanFieldsStr},
		},
		ResponseParser: masherytypes.ParseMasheryPlan,
	}

	if d, err := c.UpdateObject(ctx, plan, opContext); err == nil {
		rv, _ := d.(masherytypes.Plan)
		rv.ParentPackageId = plan.ParentPackageId
		return &rv, nil
	} else {
		return nil, err
	}
}

func DeletePlan(ctx context.Context, ident masherytypes.PackagePlanIdentifier, c *transport.V3Transport) error {
	opContext := transport.FetchSpec{
		Resource:   fmt.Sprintf("/packages/%s/plans/%s", ident.PackageId, ident.PlanId),
		AppContext: "plan",
	}

	return c.DeleteObject(ctx, opContext)
}

func CountPlans(ctx context.Context, packageId masherytypes.PackageIdentifier, c *transport.V3Transport) (int64, error) {
	opCtx := transport.FetchSpec{
		Pagination: transport.NotRequired,
		Resource:   fmt.Sprintf("/packages/%s/plans", packageId.PackageId),
		AppContext: "count plans",
	}

	return c.Count(ctx, opCtx)
}

func ListPlans(ctx context.Context, packageId masherytypes.PackageIdentifier, c *transport.V3Transport) ([]masherytypes.Plan, error) {
	opCtx := transport.FetchSpec{
		Pagination:     transport.PerPage,
		Resource:       fmt.Sprintf("/packages/%s/plans", packageId.PackageId),
		Query:          nil,
		AppContext:     "list plans",
		ResponseParser: masherytypes.ParseMasheryPlanArray,
	}

	if d, err := c.FetchAll(ctx, opCtx); err != nil {
		return []masherytypes.Plan{}, nil
	} else {
		// Convert individual fetches into the array of elements
		var rv []masherytypes.Plan
		for _, raw := range d {
			ms, ok := raw.([]masherytypes.Plan)
			if ok {
				rv = append(rv, ms...)
			}
		}

		for _, p := range rv {
			p.ParentPackageId = packageId
		}

		return rv, nil
	}
}

func ListPlanServices(ctx context.Context, ident masherytypes.PackagePlanIdentifier, c *transport.V3Transport) ([]masherytypes.Service, error) {
	opCtx := transport.FetchSpec{
		Pagination:     transport.PerPage,
		Resource:       fmt.Sprintf("/packages/%s/plans/%s/services", ident.PackageId, ident.PlanId),
		Query:          nil,
		AppContext:     "list plan service",
		ResponseParser: masherytypes.ParseMasheryServiceArray,
	}

	if d, err := c.FetchAll(ctx, opCtx); err != nil {
		return []masherytypes.Service{}, nil
	} else {
		// Convert individual fetches into the array of elements
		var rv []masherytypes.Service
		for _, raw := range d {
			ms, ok := raw.([]masherytypes.Service)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}
