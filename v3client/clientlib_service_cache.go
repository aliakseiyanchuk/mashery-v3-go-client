package v3client

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
)

func masheryServiceCacheSpec(id string) transport.FetchSpec {
	return transport.FetchSpec{
		Resource:       fmt.Sprintf("/services/%s/cache", id),
		Query:          nil,
		AppContext:     "service cache",
		ResponseParser: masherytypes.ParseMasheryServiceCache,
	}
}

// GetServiceCache Retrieve the service cache
func GetServiceCache(ctx context.Context, id string, c *transport.V3Transport) (*masherytypes.MasheryServiceCache, error) {
	rv, err := c.GetObject(ctx, masheryServiceCacheSpec(id))

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(masherytypes.MasheryServiceCache)
		return &retServ, nil
	}
}

// CreateServiceCache Create a new service cache
func CreateServiceCache(ctx context.Context, id string, service masherytypes.MasheryServiceCache, c *transport.V3Transport) (*masherytypes.MasheryServiceCache, error) {
	rawResp, err := c.CreateObject(ctx, service, masheryServiceCacheSpec(id))

	if err == nil {
		rv, _ := rawResp.(masherytypes.MasheryServiceCache)
		return &rv, nil
	} else {
		return nil, err
	}
}

// UpdateServiceCache Update cache of this service
func UpdateServiceCache(ctx context.Context, id string, service masherytypes.MasheryServiceCache, c *transport.V3Transport) (*masherytypes.MasheryServiceCache, error) {
	if d, err := c.UpdateObject(ctx, service, masheryServiceCacheSpec(id)); err == nil {
		rv, _ := d.(masherytypes.MasheryServiceCache)
		return &rv, nil
	} else {
		return nil, err
	}
}

// DeleteServiceCache Create a new service.
func DeleteServiceCache(ctx context.Context, id string, c *transport.V3Transport) error {
	return c.DeleteObject(ctx, masheryServiceCacheSpec(id))
}
