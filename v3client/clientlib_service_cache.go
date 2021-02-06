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
func GetServiceCache(ctx context.Context, id string, c *HttpTransport) (*MasheryServiceCache, error) {
	rv, err := c.getObject(ctx, masheryServiceCacheSpec(id))

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(MasheryServiceCache)
		return &retServ, nil
	}
}

// Create a new service cache
func CreateServiceCache(ctx context.Context, id string, service MasheryServiceCache, c *HttpTransport) (*MasheryServiceCache, error) {
	rawResp, err := c.createObject(ctx, service, masheryServiceCacheSpec(id))

	if err == nil {
		rv, _ := rawResp.(MasheryServiceCache)
		return &rv, nil
	} else {
		return nil, err
	}
}

// Update cache of this service
func UpdateServiceCache(ctx context.Context, id string, service MasheryServiceCache, c *HttpTransport) (*MasheryServiceCache, error) {
	if d, err := c.updateObject(ctx, service, masheryServiceCacheSpec(id)); err == nil {
		rv, _ := d.(MasheryServiceCache)
		return &rv, nil
	} else {
		return nil, err
	}
}

// Create a new service.
func DeleteServiceCache(ctx context.Context, id string, c *HttpTransport) error {
	return c.deleteObject(ctx, masheryServiceCacheSpec(id))
}
