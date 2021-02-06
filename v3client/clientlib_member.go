package v3client

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func GetMember(ctx context.Context, id string, c *HttpTransport) (*MasheryMember, error) {
	qs := url.Values{
		"fields": {memberFieldsStr},
	}

	return fetchMember(ctx, id, qs, c)
}

func GetFullMember(ctx context.Context, id string, c *HttpTransport) (*MasheryMember, error) {
	qs := url.Values{
		"fields": {strings.Join(memberDeepFields, ",")},
	}

	return fetchMember(ctx, id, qs, c)
}

func fetchMember(ctx context.Context, id string, qs url.Values, c *HttpTransport) (*MasheryMember, error) {
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
func CreateMember(ctx context.Context, member MasheryMember, c *HttpTransport) (*MasheryMember, error) {
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
func UpdateMember(ctx context.Context, member MasheryMember, c *HttpTransport) (*MasheryMember, error) {
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

func DeleteMember(ctx context.Context, memberId string, c *HttpTransport) error {
	opContext := FetchSpec{
		Resource:   fmt.Sprintf("/members/%s", memberId),
		AppContext: "member",
	}

	return c.deleteObject(ctx, opContext)
}

func ListMembers(ctx context.Context, c *HttpTransport) ([]MasheryMember, error) {
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
