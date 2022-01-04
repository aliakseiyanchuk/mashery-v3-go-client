package v3client

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
)

func ListPublicDomains(ctx context.Context, c *transport.V3Transport) ([]string, error) {
	spec := transport.FetchSpec{
		Pagination:     transport.PerPage,
		Resource:       "/domains/public/hostnames",
		Query:          nil,
		AppContext:     "unique public domains",
		ResponseParser: masherytypes.ParseMasheryDomainAddressArray,
	}

	if d, err := c.FetchAll(ctx, spec); err != nil {
		return []string{}, err
	} else {
		return mulitDomainAddressToStringArray(d), nil
	}
}

func mulitDomainAddressToStringArray(inp []interface{}) []string {
	var rv []masherytypes.DomainAddress
	for _, raw := range inp {
		if dAddr, ok := raw.([]masherytypes.DomainAddress); ok {
			rv = append(rv, dAddr...)
		}
	}

	strRv := make([]string, len(rv))
	for i, addr := range rv {
		strRv[i] = addr.Address
	}

	return strRv
}

func ListSystemDomains(ctx context.Context, c *transport.V3Transport) ([]string, error) {
	spec := transport.FetchSpec{
		Pagination:     transport.PerPage,
		Resource:       "/domains/system/hostnames",
		Query:          nil,
		AppContext:     "unique system domains",
		ResponseParser: masherytypes.ParseMasheryDomainAddressArray,
	}

	if d, err := c.FetchAll(ctx, spec); err != nil {
		return []string{}, err
	} else {
		return mulitDomainAddressToStringArray(d), nil
	}
}
