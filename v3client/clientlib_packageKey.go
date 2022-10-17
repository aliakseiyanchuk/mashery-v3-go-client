package v3client

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/errwrap"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"net/url"
)

func GetPackageKey(ctx context.Context, id masherytypes.PackageKeyIdentifier, c *transport.V3Transport) (*masherytypes.PackageKey, error) {

	rv, err := c.GetObject(ctx, transport.FetchSpec{
		Resource: fmt.Sprintf("/packageKeys/%s", id.PackageKeyId),
		Query: url.Values{
			"filter": MasheryPackageKeyFullFields,
		},
		AppContext:     "package key",
		ResponseParser: masherytypes.ParseMasheryPackageKey,
	})

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(masherytypes.PackageKey)
		return &retServ, nil
	}
}

// CreatePackageKey Create a new service.
func CreatePackageKey(ctx context.Context, appId masherytypes.ApplicationIdentifier, packageKey masherytypes.PackageKey, c *transport.V3Transport) (*masherytypes.PackageKey, error) {
	if !packageKey.LinksPackageAndPlan() {
		return nil, &errwrap.WrappedError{
			Context: "create package key",
			Cause:   errors.New("package key must supply associated package and plan"),
		}
	}
	rawResp, err := c.CreateObject(ctx, packageKey, transport.FetchSpec{
		Resource:       fmt.Sprintf("/applications/%s/packageKeys", appId),
		AppContext:     "package key",
		ResponseParser: masherytypes.ParseMasheryPackageKey,
	})

	if err == nil {
		rv, _ := rawResp.(masherytypes.PackageKey)
		return &rv, nil
	} else {
		return nil, err
	}
}

// UpdatePackageKey Create a new service.
func UpdatePackageKey(ctx context.Context, packageKey masherytypes.PackageKey, c *transport.V3Transport) (*masherytypes.PackageKey, error) {
	if packageKey.Id == "" {
		return nil, errors.New("illegal argument: package key Id must be set and not nil")
	}

	opContext := transport.FetchSpec{
		Resource:       fmt.Sprintf("/packageKeys/%s", packageKey.Id),
		AppContext:     "package key",
		ResponseParser: masherytypes.ParseMasheryPackageKey,
	}

	if d, err := c.UpdateObject(ctx, packageKey, opContext); err == nil {
		rv, _ := d.(masherytypes.PackageKey)
		return &rv, nil
	} else {
		return nil, err
	}
}

func DeletePackageKey(ctx context.Context, keyId masherytypes.PackageKeyIdentifier, c *transport.V3Transport) error {
	opSpec := transport.FetchSpec{
		Resource:       fmt.Sprintf("/packageKeys/%s", keyId.PackageKeyId),
		AppContext:     "package key",
		ResponseParser: nil,
	}

	return c.DeleteObject(ctx, opSpec)
}

func ListPackageKeysFiltered(ctx context.Context, params map[string]string, fields []string, c *transport.V3Transport) ([]masherytypes.PackageKey, error) {
	return listPackageKeysWithQuery(ctx, c.V3FilteringParams(params, fields), c)
}

func ListPackageKeys(ctx context.Context, c *transport.V3Transport) ([]masherytypes.PackageKey, error) {
	return listPackageKeysWithQuery(ctx, nil, c)
}

func listPackageKeysWithQuery(ctx context.Context, qs url.Values, c *transport.V3Transport) ([]masherytypes.PackageKey, error) {
	opCtx := transport.FetchSpec{
		Pagination:     transport.PerPage,
		Resource:       "/packageKeys",
		Query:          qs,
		AppContext:     "all package keys",
		ResponseParser: masherytypes.ParseMasheryPackageKeyArray,
	}

	if d, err := c.FetchAll(ctx, opCtx); err != nil {
		return []masherytypes.PackageKey{}, err
	} else {
		// Convert individual fetches into the array of elements
		var rv []masherytypes.PackageKey
		for _, raw := range d {
			ms, ok := raw.([]masherytypes.PackageKey)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}
