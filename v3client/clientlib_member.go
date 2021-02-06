package v3client

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func (c *HttpTransport) GetMember(ctx context.Context, id string) (*MasheryMember, error) {
	qs := url.Values{
		"fields": {memberFieldsStr},
	}

	return c.fetchMember(ctx, id, qs)
}

func (c *HttpTransport) GetFullMember(ctx context.Context, id string) (*MasheryMember, error) {
	qs := url.Values{
		"fields": {strings.Join(memberDeepFields, ",")},
	}

	return c.fetchMember(ctx, id, qs)
}

func (c *HttpTransport) fetchMember(ctx context.Context, id string, qs url.Values) (*MasheryMember, error) {
	rv, err := c.getObject(ctx, FetchSpec{
		Resource:       fmt.Sprintf("/members/%s", id),
		Query:          qs,
		AppContext:     "member",
		ResponseParser: ParseMasheryMember,
	})

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(MasheryMember)
		return &retServ, nil
	}
}

// Create a new service.
func (c *HttpTransport) CreateMember(ctx context.Context, member MasheryMember) (*MasheryMember, error) {
	rawResp, err := c.createObject(ctx, member, FetchSpec{
		Resource:       "/members",
		AppContext:     "members",
		ResponseParser: ParseMasheryMember,
	})

	if err == nil {
		rv, _ := rawResp.(MasheryMember)
		return &rv, nil
	} else {
		return nil, err
	}
}

// Create a new service.
func (c *HttpTransport) UpdateMember(ctx context.Context, member MasheryMember) (*MasheryMember, error) {
	if member.Id == "" {
		return nil, errors.New("illegal argument: member Id must be set and not nil")
	}

	opContext := FetchSpec{
		Resource:       fmt.Sprintf("/members/%s", member.Id),
		AppContext:     "service",
		ResponseParser: ParseMasheryMember,
	}

	if d, err := c.updateObject(ctx, member, opContext); err == nil {
		rv, _ := d.(MasheryMember)
		return &rv, nil
	} else {
		return nil, err
	}
}

func (c *HttpTransport) DeleteMember(ctx context.Context, memberId string) error {
	opContext := FetchSpec{
		Resource:   fmt.Sprintf("/members/%s", memberId),
		AppContext: "member",
	}

	return c.deleteObject(ctx, opContext)
}

func (c *HttpTransport) ListMembers(ctx context.Context) ([]MasheryMember, error) {
	opCtx := FetchSpec{
		Pagination:     PerPage,
		Resource:       "/members",
		Query:          nil,
		AppContext:     "all members",
		ResponseParser: ParseMasheryMemberArray,
	}

	if d, err := c.fetchAll(ctx, opCtx); err != nil {
		return []MasheryMember{}, nil
	} else {
		// Convert individual fetches into the array of elements
		var rv []MasheryMember
		for _, raw := range d {
			ms, ok := raw.([]MasheryMember)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}
