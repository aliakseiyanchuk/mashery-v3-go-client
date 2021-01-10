package v3client

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

func (c *Client) GetPlan(ctx context.Context, packageId string, planId string) (*MasheryPlan, error) {
	qs := url.Values{
		"fields": {
			"id", "created", "updated", "name", "description", "eav", "selfServiceKeyProvisioningEnabled",
			"adminKeyProvisioningEnabled", "notes", "maxNumKeysAllowed", "numKeysBeforeReview", "qpsLimitCeiling", "qpsLimitExempt",
			"qpsLimitKeyOverrideAllowed", "rateLimitCeiling", "rateLimitExempt", "rateLimitKeyOverrideAllowed", "rateLimitPeriod",
			"responseFilterOverrideAllowed", "status", "emailTemplateSetId", "services.dd.pp", "services.id", "services.endpoints",
		},
	}

	rv, err := c.getObject(ctx, FetchSpec{
		Resource:       fmt.Sprintf("/packages/%s/plans/%s", packageId, planId),
		Query:          qs,
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
func (c *Client) CreatePlan(ctx context.Context, packageId string, plan MasheryPlan) (*MasheryPlan, error) {
	rawResp, err := c.createObject(ctx, plan, FetchSpec{
		Resource:       fmt.Sprintf("/packages/%s/plans", packageId),
		AppContext:     "plan",
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
func (c *Client) UpdatePlan(ctx context.Context, plan MasheryPlan) (*MasheryPlan, error) {
	if plan.Id == "" || plan.ParentPackageId == "" {
		return nil, errors.New("illegal argument: package Id and plan id must be set and not nil")
	}

	opContext := FetchSpec{
		Resource:       fmt.Sprintf("/packages/%s/plans/%s", plan.ParentPackageId, plan.Id),
		AppContext:     "plan",
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

func (c *Client) ListPlans(ctx context.Context, packageId string) ([]MasheryPlan, error) {
	opCtx := FetchSpec{
		Pagination:     PerPage,
		Resource:       fmt.Sprintf("/packages/%s/plans", packageId),
		Query:          nil,
		AppContext:     "all service",
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

func (c *Client) ListPlanServices(ctx context.Context, packageId string, planId string) ([]MasheryService, error) {
	opCtx := FetchSpec{
		Pagination:     PerPage,
		Resource:       fmt.Sprintf("/packages/%s/plans/%s/services", packageId, planId),
		Query:          nil,
		AppContext:     "all service",
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
