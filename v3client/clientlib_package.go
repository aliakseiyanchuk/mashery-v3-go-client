package v3client

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"net/url"
)

func GetPackage(ctx context.Context, id masherytypes.PackageIdentifier, c *transport.V3Transport) (*masherytypes.Package, error) {
	rv, err := c.GetObject(ctx, transport.FetchSpec{
		Resource: fmt.Sprintf("/packages/%s", id.PackageId),
		Query: url.Values{
			"fields": {MasheryPackageFieldsStr},
		},
		AppContext:     "service",
		ResponseParser: masherytypes.ParseMasheryPackage,
	})

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(masherytypes.Package)
		return &retServ, nil
	}
}

// CreatePackage Create a new service.
func CreatePackage(ctx context.Context, pack masherytypes.Package, c *transport.V3Transport) (*masherytypes.Package, error) {
	rawResp, err := c.CreateObject(ctx, pack, transport.FetchSpec{
		Resource:   "/packages",
		AppContext: "package",
		Query: url.Values{
			"fields": {MasheryPackageFieldsStr},
		},
		ResponseParser: masherytypes.ParseMasheryPackage,
	})

	if err == nil {
		rv, _ := rawResp.(masherytypes.Package)
		return &rv, nil
	} else {
		return nil, err
	}
}

// UpdatePackage Create a new service.
func UpdatePackage(ctx context.Context, pack masherytypes.Package, c *transport.V3Transport) (*masherytypes.Package, error) {
	if pack.Id == "" {
		return nil, errors.New("illegal argument: package Id must be set and not nil")
	}

	opContext := transport.FetchSpec{
		Resource:       fmt.Sprintf("/packages/%s", pack.Id),
		AppContext:     "package",
		ResponseParser: masherytypes.ParseService,
	}

	if d, err := c.UpdateObject(ctx, pack, opContext); err == nil {
		rv, _ := d.(masherytypes.Package)
		return &rv, nil
	} else {
		return nil, err
	}
}

func DeletePackage(ctx context.Context, packId masherytypes.PackageIdentifier, c *transport.V3Transport) error {
	opContext := transport.FetchSpec{
		Resource:   fmt.Sprintf("/packages/%s", packId.PackageId),
		AppContext: "package",
	}

	return c.DeleteObject(ctx, opContext)
}

func ListPackages(ctx context.Context, c *transport.V3Transport) ([]masherytypes.Package, error) {
	opCtx := transport.FetchSpec{
		Pagination:     transport.PerItem,
		Resource:       "/packages",
		Query:          nil,
		AppContext:     "all service",
		ResponseParser: masherytypes.ParseMasheryPackageArray,
	}

	if d, err := c.FetchAll(ctx, opCtx); err != nil {
		return []masherytypes.Package{}, err
	} else {
		// Convert individual fetches into the array of elements
		var rv []masherytypes.Package
		for _, raw := range d {
			ms, ok := raw.([]masherytypes.Package)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}
