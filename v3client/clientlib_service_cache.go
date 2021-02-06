package v3client

import (
	"context"
	"fmt"
)

func masheryServiceCacheSpec(id string) FetchSpec {
	return FetchSpec{
		Resource:       fmt.Sprintf("/services/%s/cache", id),
		Query:          nil,
		AppContext:     "service cache",
		ResponseParser: ParseMasheryServiceCache,
	}
}

// Retrieve the service cache
func (c *HttpTransport) GetServiceCache(ctx context.Context, id string) (*MasheryServiceCache, error) {
	rv, err := c.getObject(ctx, masheryServiceCacheSpec(id))

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(MasheryServiceCache)
		return &retServ, nil
	}
}

// Create a new service cache
func (c *HttpTransport) CreateServiceCache(ctx context.Context, id string, service MasheryServiceCache) (*MasheryServiceCache, error) {
	rawResp, err := c.createObject(ctx, service, masheryServiceCacheSpec(id))

	if err == nil {
		rv, _ := rawResp.(MasheryServiceCache)
		return &rv, nil
	} else {
		return nil, err
	}
}

// Update cache of this service
func (c *HttpTransport) UpdateServiceCache(ctx context.Context, id string, service MasheryServiceCache) (*MasheryServiceCache, error) {
	if d, err := c.updateObject(ctx, service, masheryServiceCacheSpec(id)); err == nil {
		rv, _ := d.(MasheryServiceCache)
		return &rv, nil
	} else {
		return nil, err
	}
}

// Create a new service.
func (c *HttpTransport) DeleteServiceCache(ctx context.Context, id string) error {
	return c.deleteObject(ctx, masheryServiceCacheSpec(id))
}
