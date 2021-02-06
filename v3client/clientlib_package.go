package v3client

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

func GetPackage(ctx context.Context, id string, c *HttpTransport) (*MasheryPackage, error) {
	rv, err := c.getObject(ctx, FetchSpec{
		Resource: fmt.Sprintf("/packages/%s", id),
		Query: url.Values{
			"fields": {MasheryPackageFieldsStr},
		},
		AppContext:     "service",
		ResponseParser: ParseMasheryPackage,
	})

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(MasheryPackage)
		return &retServ, nil
	}
}

// Create a new service.
func CreatePackage(ctx context.Context, pack MasheryPackage, c *HttpTransport) (*MasheryPackage, error) {
	rawResp, err := c.createObject(ctx, pack, FetchSpec{
		Resource:   "/packages",
		AppContext: "package",
		Query: url.Values{
			"fields": {MasheryPackageFieldsStr},
		},
		ResponseParser: ParseMasheryPackage,
	})

	if err == nil {
		rv, _ := rawResp.(MasheryPackage)
		return &rv, nil
	} else {
		return nil, err
	}
}

// Create a new service.
func UpdatePackage(ctx context.Context, pack MasheryPackage, c *HttpTransport) (*MasheryPackage, error) {
	if pack.Id == "" {
		return nil, errors.New("illegal argument: package Id must be set and not nil")
	}

	opContext := FetchSpec{
		Resource:       fmt.Sprintf("/packages/%s", pack.Id),
		AppContext:     "package",
		ResponseParser: ParseMasheryService,
	}

	if d, err := c.updateObject(ctx, pack, opContext); err == nil {
		rv, _ := d.(MasheryPackage)
		return &rv, nil
	} else {
		return nil, err
	}
}

func DeletePackage(ctx context.Context, packId string, c *HttpTransport) error {
	opContext := FetchSpec{
		Resource:   fmt.Sprintf("/packages/%s", packId),
		AppContext: "package",
	}

	return c.deleteObject(ctx, opContext)
}

func ListPackages(ctx context.Context, c *HttpTransport) ([]MasheryPackage, error) {
	opCtx := FetchSpec{
		Pagination:     PerItem,
		Resource:       "/packages",
		Query:          nil,
		AppContext:     "all service",
		ResponseParser: ParseMasheryPackageArray,
	}

	if d, err := c.fetchAll(ctx, opCtx); err != nil {
		return []MasheryPackage{}, err
	} else {
		// Convert individual fetches into the array of elements
		var rv []MasheryPackage
		for _, raw := range d {
			ms, ok := raw.([]MasheryPackage)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}
