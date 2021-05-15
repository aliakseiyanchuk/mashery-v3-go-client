package v3client

import "context"

func ListPublicDomains(ctx context.Context, c *HttpTransport) ([]string, error) {
	spec := FetchSpec{
		Pagination:     PerPage,
		Resource:       "/domains/public/hostnames",
		Query:          nil,
		AppContext:     "unique public domains",
		ResponseParser: ParseMasheryDomainAddressArray,
	}

	if d, err := c.fetchAll(ctx, spec); err != nil {
		return []string{}, err
	} else {
		return mulitDomainAddressToStringArray(d), nil
	}
}

func mulitDomainAddressToStringArray(inp []interface{}) []string {
	var rv []DomainAddress
	for _, raw := range inp {
		if dAddr, ok := raw.([]DomainAddress); ok {
			rv = append(rv, dAddr...)
		}
	}

	strRv := make([]string, len(rv))
	for i, addr := range rv {
		strRv[i] = addr.Address
	}

	return strRv
}

func ListSystemDomains(ctx context.Context, c *HttpTransport) ([]string, error) {
	spec := FetchSpec{
		Pagination:     PerPage,
		Resource:       "/domains/system/hostnames",
		Query:          nil,
		AppContext:     "unique system domains",
		ResponseParser: ParseMasheryDomainAddressArray,
	}

	if d, err := c.fetchAll(ctx, spec); err != nil {
		return []string{}, err
	} else {
		return mulitDomainAddressToStringArray(d), nil
	}
}
