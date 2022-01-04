package v3client

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"net/url"
	"strings"
)

func GetMember(ctx context.Context, id string, c *transport.V3Transport) (*masherytypes.MasheryMember, error) {
	qs := url.Values{
		"fields": {memberFieldsStr},
	}

	return fetchMember(ctx, id, qs, c)
}

func GetFullMember(ctx context.Context, id string, c *transport.V3Transport) (*masherytypes.MasheryMember, error) {
	qs := url.Values{
		"fields": {strings.Join(memberDeepFields, ",")},
	}

	return fetchMember(ctx, id, qs, c)
}

func fetchMember(ctx context.Context, id string, qs url.Values, c *transport.V3Transport) (*masherytypes.MasheryMember, error) {
	rv, err := c.GetObject(ctx, transport.FetchSpec{
		Resource:       fmt.Sprintf("/members/%s", id),
		Query:          qs,
		AppContext:     "member",
		ResponseParser: masherytypes.ParseMasheryMember,
	})

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(masherytypes.MasheryMember)
		return &retServ, nil
	}
}

// CreateMember Create a new service.
func CreateMember(ctx context.Context, member masherytypes.MasheryMember, c *transport.V3Transport) (*masherytypes.MasheryMember, error) {
	rawResp, err := c.CreateObject(ctx, member, transport.FetchSpec{
		Resource:       "/members",
		AppContext:     "members",
		ResponseParser: masherytypes.ParseMasheryMember,
	})

	if err == nil {
		rv, _ := rawResp.(masherytypes.MasheryMember)
		return &rv, nil
	} else {
		return nil, err
	}
}

// UpdateMember Create a new service.
func UpdateMember(ctx context.Context, member masherytypes.MasheryMember, c *transport.V3Transport) (*masherytypes.MasheryMember, error) {
	if member.Id == "" {
		return nil, errors.New("illegal argument: member Id must be set and not nil")
	}

	opContext := transport.FetchSpec{
		Resource:       fmt.Sprintf("/members/%s", member.Id),
		AppContext:     "service",
		ResponseParser: masherytypes.ParseMasheryMember,
	}

	if d, err := c.UpdateObject(ctx, member, opContext); err == nil {
		rv, _ := d.(masherytypes.MasheryMember)
		return &rv, nil
	} else {
		return nil, err
	}
}

func DeleteMember(ctx context.Context, memberId string, c *transport.V3Transport) error {
	opContext := transport.FetchSpec{
		Resource:   fmt.Sprintf("/members/%s", memberId),
		AppContext: "member",
	}

	return c.DeleteObject(ctx, opContext)
}

func ListMembers(ctx context.Context, c *transport.V3Transport) ([]masherytypes.MasheryMember, error) {
	opCtx := transport.FetchSpec{
		Pagination:     transport.PerPage,
		Resource:       "/members",
		Query:          nil,
		AppContext:     "all members",
		ResponseParser: masherytypes.ParseMasheryMemberArray,
	}

	if d, err := c.FetchAll(ctx, opCtx); err != nil {
		return []masherytypes.MasheryMember{}, nil
	} else {
		// Convert individual fetches into the array of elements
		var rv []masherytypes.MasheryMember
		for _, raw := range d {
			ms, ok := raw.([]masherytypes.MasheryMember)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}
