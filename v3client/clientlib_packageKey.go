package v3client

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func (c *Client) GetPackageKey(ctx context.Context, id string) (*MasheryPackageKey, error) {
	qs := url.Values{
		"fields": {strings.Join(packageKeyFields, ",")},
	}

	rv, err := c.getObject(ctx, FetchSpec{
		Resource:       fmt.Sprintf("/packageKeys/%s", id),
		Query:          qs,
		AppContext:     "package key",
		ResponseParser: ParseMasheryPackageKey,
	})

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(MasheryPackageKey)
		return &retServ, nil
	}
}

// Create a new service.
func (c *Client) CreatePackageKey(ctx context.Context, service MasheryPackageKey) (*MasheryPackageKey, error) {
	rawResp, err := c.createObject(ctx, service, FetchSpec{
		Resource:       "/packageKeys",
		AppContext:     "package key",
		ResponseParser: ParseMasheryPackageKey,
	})

	if err == nil {
		rv, _ := rawResp.(MasheryPackageKey)
		return &rv, nil
	} else {
		return nil, err
	}
}

// Create a new service.
func (c *Client) UpdatePackageKey(ctx context.Context, packageKey MasheryPackageKey) (*MasheryPackageKey, error) {
	if packageKey.Id == "" {
		return nil, errors.New("illegal argument: package key Id must be set and not nil")
	}

	opContext := FetchSpec{
		Resource:       fmt.Sprintf("/packageKeys/%s", packageKey.Id),
		AppContext:     "package key",
		ResponseParser: ParseMasheryPackageKey,
	}

	if d, err := c.updateObject(ctx, packageKey, opContext); err == nil {
		rv, _ := d.(MasheryPackageKey)
		return &rv, nil
	} else {
		return nil, err
	}
}

func (c *Client) DeletePackageKey(ctx context.Context, keyId string) error {
	opSpec := FetchSpec{
		Resource:       fmt.Sprintf("/packageKeys/%s", keyId),
		AppContext:     "packageKeys",
		ResponseParser: nil,
	}

	return c.deleteObject(ctx, opSpec)
}

func (c *Client) ListPackageKeys(ctx context.Context) ([]MasheryPackageKey, error) {
	opCtx := FetchSpec{
		Pagination:     PerPage,
		Resource:       "/packageKeys",
		Query:          nil,
		AppContext:     "all package keys",
		ResponseParser: ParseMasheryPackageKeyArray,
	}

	if d, err := c.fetchAll(ctx, opCtx); err != nil {
		return []MasheryPackageKey{}, err
	} else {
		// Convert individual fetches into the array of elements
		var rv []MasheryPackageKey
		for _, raw := range d {
			ms, ok := raw.([]MasheryPackageKey)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}
