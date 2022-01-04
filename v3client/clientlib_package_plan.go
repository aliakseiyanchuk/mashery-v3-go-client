package v3client

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"net/url"
)

func CreatePlanService(ctx context.Context, planService masherytypes.MasheryPlanService, c *transport.V3Transport) (*masherytypes.AddressableV3Object, error) {
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

func DeletePlanService(ctx context.Context, planService masherytypes.MasheryPlanService, c *transport.V3Transport) error {
	return c.DeleteObject(ctx, transport.FetchSpec{
		Resource:   fmt.Sprintf("/packages/%s/plans/%s/services/%s", planService.PackageId, planService.PlanId, planService.ServiceId),
		AppContext: "plan service",
	})
}

func CreatePlanEndpoint(ctx context.Context, planEndp masherytypes.MasheryPlanServiceEndpoint, c *transport.V3Transport) (*masherytypes.AddressableV3Object, error) {
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

func DeletePlanEndpoint(ctx context.Context, planEndp masherytypes.MasheryPlanServiceEndpoint, c *transport.V3Transport) error {
	return c.DeleteObject(ctx, transport.FetchSpec{
		Resource:   fmt.Sprintf("/packages/%s/plans/%s/services/%s/endpoints/%s", planEndp.PackageId, planEndp.PlanId, planEndp.ServiceId, planEndp.EndpointId),
		AppContext: "plan endpoint",
	})
}

func ListPlanEndpoints(ctx context.Context, planService masherytypes.MasheryPlanService, c *transport.V3Transport) ([]masherytypes.AddressableV3Object, error) {
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

func CountPlanEndpoints(ctx context.Context, planService masherytypes.MasheryPlanService, c *transport.V3Transport) (int64, error) {
	opCtx := transport.FetchSpec{
		Pagination: transport.NotRequired,
		Resource:   fmt.Sprintf("/packages/%s/plans/%s/services/%s/endpoints", planService.PackageId, planService.PlanId, planService.ServiceId),
		AppContext: "count plan service endpoints",
	}

	return c.Count(ctx, opCtx)
}

func CountPlanService(ctx context.Context, packageId, planId string, c *transport.V3Transport) (int64, error) {
	opCtx := transport.FetchSpec{
		Pagination: transport.NotRequired,
		Resource:   fmt.Sprintf("/packages/%s/plans/%s/services", packageId, planId),
		AppContext: "count plans services",
	}

	return c.Count(ctx, opCtx)
}

func GetPlan(ctx context.Context, packageId string, planId string, c *transport.V3Transport) (*masherytypes.MasheryPlan, error) {
	rv, err := c.GetObject(ctx, transport.FetchSpec{
		Resource: fmt.Sprintf("/packages/%s/plans/%s", packageId, planId),
		Query: url.Values{
			"fields": {MasheryPlanFieldsStr},
		},
		AppContext:     "plan",
		ResponseParser: masherytypes.ParseMasheryPlan,
	})

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(masherytypes.MasheryPlan)
		retServ.ParentPackageId = packageId

		return &retServ, nil
	}
}

// CreatePlan Create a new service.
func CreatePlan(ctx context.Context, packageId string, plan masherytypes.MasheryPlan, c *transport.V3Transport) (*masherytypes.MasheryPlan, error) {
	rawResp, err := c.CreateObject(ctx, plan, transport.FetchSpec{
		Resource:   fmt.Sprintf("/packages/%s/plans", packageId),
		AppContext: "plan",
		Query: url.Values{
			"fields": {MasheryPlanFieldsStr},
		},
		ResponseParser: masherytypes.ParseMasheryPlan,
	})

	if err == nil {
		rv, _ := rawResp.(masherytypes.MasheryPlan)
		rv.ParentPackageId = packageId
		return &rv, nil
	} else {
		return nil, err
	}
}

// UpdatePlan Create a new service.
func UpdatePlan(ctx context.Context, plan masherytypes.MasheryPlan, c *transport.V3Transport) (*masherytypes.MasheryPlan, error) {
	if plan.Id == "" || plan.ParentPackageId == "" {
		return nil, errors.New("illegal argument: package Id and plan id must be set and not nil")
	}

	opContext := transport.FetchSpec{
		Resource:   fmt.Sprintf("/packages/%s/plans/%s", plan.ParentPackageId, plan.Id),
		AppContext: "plan",
		Query: url.Values{
			"fields": {MasheryPlanFieldsStr},
		},
		ResponseParser: masherytypes.ParseMasheryPlan,
	}

	if d, err := c.UpdateObject(ctx, plan, opContext); err == nil {
		rv, _ := d.(masherytypes.MasheryPlan)
		rv.ParentPackageId = plan.ParentPackageId
		return &rv, nil
	} else {
		return nil, err
	}
}

func DeletePlan(ctx context.Context, packageId, planId string, c *transport.V3Transport) error {
	opContext := transport.FetchSpec{
		Resource:   fmt.Sprintf("/packages/%s/plans/%s", packageId, planId),
		AppContext: "plan",
	}

	return c.DeleteObject(ctx, opContext)
}

func CountPlans(ctx context.Context, packageId string, c *transport.V3Transport) (int64, error) {
	opCtx := transport.FetchSpec{
		Pagination: transport.NotRequired,
		Resource:   fmt.Sprintf("/packages/%s/plans", packageId),
		AppContext: "count plans",
	}

	return c.Count(ctx, opCtx)
}

func ListPlans(ctx context.Context, packageId string, c *transport.V3Transport) ([]masherytypes.MasheryPlan, error) {
	opCtx := transport.FetchSpec{
		Pagination:     transport.PerPage,
		Resource:       fmt.Sprintf("/packages/%s/plans", packageId),
		Query:          nil,
		AppContext:     "list plans",
		ResponseParser: masherytypes.ParseMasheryPlanArray,
	}

	if d, err := c.FetchAll(ctx, opCtx); err != nil {
		return []masherytypes.MasheryPlan{}, nil
	} else {
		// Convert individual fetches into the array of elements
		var rv []masherytypes.MasheryPlan
		for _, raw := range d {
			ms, ok := raw.([]masherytypes.MasheryPlan)
			if ok {
				rv = append(rv, ms...)
			}
		}

		// Push the parent package
		for _, p := range rv {
			p.ParentPackageId = packageId
		}

		return rv, nil
	}
}

func ListPlanServices(ctx context.Context, packageId string, planId string, c *transport.V3Transport) ([]masherytypes.MasheryService, error) {
	opCtx := transport.FetchSpec{
		Pagination:     transport.PerPage,
		Resource:       fmt.Sprintf("/packages/%s/plans/%s/services", packageId, planId),
		Query:          nil,
		AppContext:     "list plan service",
		ResponseParser: masherytypes.ParseMasheryServiceArray,
	}

	if d, err := c.FetchAll(ctx, opCtx); err != nil {
		return []masherytypes.MasheryService{}, nil
	} else {
		// Convert individual fetches into the array of elements
		var rv []masherytypes.MasheryService
		for _, raw := range d {
			ms, ok := raw.([]masherytypes.MasheryService)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}
