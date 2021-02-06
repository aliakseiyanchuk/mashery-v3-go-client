package v3client

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

func CreatePlanService(ctx context.Context, planService MasheryPlanService, c *HttpTransport) (*AddressableV3Object, error) {
	ref := IdReferenced{IdRef: planService.ServiceId}
	rv, err := c.createObject(ctx, ref, FetchSpec{
		Pagination:     NotRequired,
		Resource:       fmt.Sprintf("/packages/%s/plans/%s/services", planService.PackageId, planService.PlanId),
		Query:          nil,
		AppContext:     "plan service",
		ResponseParser: ParseMasheryAddressableObject,
	})

	if err != nil {
		return nil, err
	} else {
		rvc := rv.(AddressableV3Object)
		return &rvc, nil
	}
}

func DeletePlanService(ctx context.Context, planService MasheryPlanService, c *HttpTransport) error {
	return c.deleteObject(ctx, FetchSpec{
		Resource:   fmt.Sprintf("/packages/%s/plans/%s/services/%s", planService.PackageId, planService.PlanId, planService.ServiceId),
		AppContext: "plan service",
	})
}

func CreatePlanEndpoint(ctx context.Context, planEndp MasheryPlanServiceEndpoint, c *HttpTransport) (*AddressableV3Object, error) {
	ref := IdReferenced{IdRef: planEndp.EndpointId}
	rv, err := c.createObject(ctx, ref, FetchSpec{
		Resource:       fmt.Sprintf("/packages/%s/plans/%s/services/%s/endpoints", planEndp.PackageId, planEndp.PlanId, planEndp.ServiceId),
		AppContext:     "create plan endpoint",
		ResponseParser: ParseMasheryAddressableObject,
	})

	if err != nil {
		return nil, err
	} else {
		rvc := rv.(AddressableV3Object)
		return &rvc, nil
	}
}

func DeletePlanEndpoint(ctx context.Context, planEndp MasheryPlanServiceEndpoint, c *HttpTransport) error {
	return c.deleteObject(ctx, FetchSpec{
		Resource:   fmt.Sprintf("/packages/%s/plans/%s/services/%s/endpoints/%s", planEndp.PackageId, planEndp.PlanId, planEndp.ServiceId, planEndp.EndpointId),
		AppContext: "plan endpoint",
	})
}

func ListPlanEndpoints(ctx context.Context, planService MasheryPlanService, c *HttpTransport) ([]AddressableV3Object, error) {
	rv, err := c.fetchAll(ctx, FetchSpec{
		Pagination:     PerItem,
		Resource:       fmt.Sprintf("/packages/%s/plans/%s/services/%s/endpoints", planService.PackageId, planService.PlanId, planService.ServiceId),
		Query:          nil,
		AppContext:     "list plan service endpoints",
		ResponseParser: ParseMasheryAddressableObjectsArray,
	})

	if err != nil {
		return []AddressableV3Object{}, err
	} else {
		rvCast := make([]AddressableV3Object, len(rv))
		for i, v := range rv {
			rvCast[i] = v.(AddressableV3Object)
		}

		return rvCast, nil
	}
}

func CountPlanEndpoints(ctx context.Context, planService MasheryPlanService, c *HttpTransport) (int64, error) {
	opCtx := FetchSpec{
		Pagination: NotRequired,
		Resource:   fmt.Sprintf("/packages/%s/plans/%s/services/%s/endpoints", planService.PackageId, planService.PlanId, planService.ServiceId),
		AppContext: "count plan service endpoints",
	}

	return c.count(ctx, opCtx)
}

func CountPlanService(ctx context.Context, packageId, planId string, c *HttpTransport) (int64, error) {
	opCtx := FetchSpec{
		Pagination: NotRequired,
		Resource:   fmt.Sprintf("/packages/%s/plans/%s/services", packageId, planId),
		AppContext: "count plans services",
	}

	return c.count(ctx, opCtx)
}

func GetPlan(ctx context.Context, packageId string, planId string, c *HttpTransport) (*MasheryPlan, error) {
	rv, err := c.getObject(ctx, FetchSpec{
		Resource: fmt.Sprintf("/packages/%s/plans/%s", packageId, planId),
		Query: url.Values{
			"fields": {MasheryPlanFieldsStr},
		},
		AppContext:     "plan",
		ResponseParser: ParseMasheryPlan,
	})

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(MasheryPlan)
		retServ.ParentPackageId = packageId

		return &retServ, nil
	}
}

// Create a new service.
func CreatePlan(ctx context.Context, packageId string, plan MasheryPlan, c *HttpTransport) (*MasheryPlan, error) {
	rawResp, err := c.createObject(ctx, plan, FetchSpec{
		Resource:   fmt.Sprintf("/packages/%s/plans", packageId),
		AppContext: "plan",
		Query: url.Values{
			"fields": {MasheryPlanFieldsStr},
		},
		ResponseParser: ParseMasheryPlan,
	})

	if err == nil {
		rv, _ := rawResp.(MasheryPlan)
		rv.ParentPackageId = packageId
		return &rv, nil
	} else {
		return nil, err
	}
}

// Create a new service.
func UpdatePlan(ctx context.Context, plan MasheryPlan, c *HttpTransport) (*MasheryPlan, error) {
	if plan.Id == "" || plan.ParentPackageId == "" {
		return nil, errors.New("illegal argument: package Id and plan id must be set and not nil")
	}

	opContext := FetchSpec{
		Resource:   fmt.Sprintf("/packages/%s/plans/%s", plan.ParentPackageId, plan.Id),
		AppContext: "plan",
		Query: url.Values{
			"fields": {MasheryPlanFieldsStr},
		},
		ResponseParser: ParseMasheryPlan,
	}

	if d, err := c.updateObject(ctx, plan, opContext); err == nil {
		rv, _ := d.(MasheryPlan)
		rv.ParentPackageId = plan.ParentPackageId
		return &rv, nil
	} else {
		return nil, err
	}
}

func DeletePlan(ctx context.Context, packageId, planId string, c *HttpTransport) error {
	opContext := FetchSpec{
		Resource:   fmt.Sprintf("/packages/%s/plans/%s", packageId, planId),
		AppContext: "plan",
	}

	return c.deleteObject(ctx, opContext)
}

func CountPlans(ctx context.Context, packageId string, c *HttpTransport) (int64, error) {
	opCtx := FetchSpec{
		Pagination: NotRequired,
		Resource:   fmt.Sprintf("/packages/%s/plans", packageId),
		AppContext: "count plans",
	}

	return c.count(ctx, opCtx)
}

func ListPlans(ctx context.Context, packageId string, c *HttpTransport) ([]MasheryPlan, error) {
	opCtx := FetchSpec{
		Pagination:     PerPage,
		Resource:       fmt.Sprintf("/packages/%s/plans", packageId),
		Query:          nil,
		AppContext:     "list plans",
		ResponseParser: ParseMasheryPlanArray,
	}

	if d, err := c.fetchAll(ctx, opCtx); err != nil {
		return []MasheryPlan{}, nil
	} else {
		// Convert individual fetches into the array of elements
		var rv []MasheryPlan
		for _, raw := range d {
			ms, ok := raw.([]MasheryPlan)
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

func ListPlanServices(ctx context.Context, packageId string, planId string, c *HttpTransport) ([]MasheryService, error) {
	opCtx := FetchSpec{
		Pagination:     PerPage,
		Resource:       fmt.Sprintf("/packages/%s/plans/%s/services", packageId, planId),
		Query:          nil,
		AppContext:     "list plan service",
		ResponseParser: ParseMasheryServiceArray,
	}

	if d, err := c.fetchAll(ctx, opCtx); err != nil {
		return []MasheryService{}, nil
	} else {
		// Convert individual fetches into the array of elements
		var rv []MasheryService
		for _, raw := range d {
			ms, ok := raw.([]MasheryService)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}
