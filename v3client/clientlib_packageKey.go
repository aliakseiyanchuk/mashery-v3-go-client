package v3client

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

func (c *HttpTransport) GetPackageKey(ctx context.Context, id string) (*MasheryPackageKey, error) {

	rv, err := c.getObject(ctx, FetchSpec{
		Resource: fmt.Sprintf("/packageKeys/%s", id),
		Query: url.Values{
			"filter": MasheryPackageKeyFullFields,
		},
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
func (c *HttpTransport) CreatePackageKey(ctx context.Context, appId string, packageKey MasheryPackageKey) (*MasheryPackageKey, error) {
	if !packageKey.LinksPackageAndPlan() {
		return nil, &WrappedError{
			Context: "create package key",
			Cause:   errors.New("package key must supply associated package and plan"),
		}
	}
	rawResp, err := c.createObject(ctx, packageKey, FetchSpec{
		Resource:       fmt.Sprintf("/applications/%s/packageKeys", appId),
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
func (c *HttpTransport) UpdatePackageKey(ctx context.Context, packageKey MasheryPackageKey) (*MasheryPackageKey, error) {
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

func (c *HttpTransport) DeletePackageKey(ctx context.Context, keyId string) error {
	opSpec := FetchSpec{
		Resource:       fmt.Sprintf("/packageKeys/%s", keyId),
		AppContext:     "package key",
		ResponseParser: nil,
	}

	return c.deleteObject(ctx, opSpec)
}

func (c *HttpTransport) ListPackageKeysFiltered(ctx context.Context, params map[string]string, fields []string) ([]MasheryPackageKey, error) {
	return c.listPackageKeysWithQuery(ctx, c.v3FilteringParams(params, fields))
}

func (c *HttpTransport) ListPackageKeys(ctx context.Context) ([]MasheryPackageKey, error) {
	return c.listPackageKeysWithQuery(ctx, nil)
}

func (c *HttpTransport) listPackageKeysWithQuery(ctx context.Context, qs url.Values) ([]MasheryPackageKey, error) {
	opCtx := FetchSpec{
		Pagination:     PerPage,
		Resource:       "/packageKeys",
		Query:          qs,
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
