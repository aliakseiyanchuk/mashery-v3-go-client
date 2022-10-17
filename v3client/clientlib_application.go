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

var applicationFields = []string{
	"id", "name", "created", "updated", "username", "description",
	"type", "commercial", "ads",
	"notes", "howDidYouHear", "preferredProtocol",
	"preferredOutput", "externalId", "uri", "oauthRedirectUri",
}

var packageKeyFields = []string{
	"id", "apikey", "secret", "created", "updated", "rateLimitCeiling", "rateLimitExempt", "qpsLimitCeiling",
	"qpsLimitExempt", "status", "limits", "package", "plan",
}

var applicationDeepFields = append(applicationFields,
	"packageKeys")

func GetApplication(ctx context.Context, appId masherytypes.ApplicationIdentifier, c *transport.V3Transport) (*masherytypes.Application, error) {
	qs := url.Values{
		"fields": {strings.Join(applicationFields, ",")},
	}

	return httpToApplication(ctx, appId, qs, c)
}

func GetApplicationPackageKeys(ctx context.Context, appId masherytypes.ApplicationIdentifier, c *transport.V3Transport) ([]masherytypes.PackageKey, error) {
	qs := url.Values{
		"fields": {strings.Join(packageKeyFields, ",")},
	}

	opCtx := transport.FetchSpec{
		Pagination:     transport.PerPage,
		Resource:       fmt.Sprintf("/applications/%s/packageKeys", appId.ApplicationId),
		Query:          qs,
		AppContext:     "application package key",
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

func CountApplicationPackageKeys(ctx context.Context, appId masherytypes.ApplicationIdentifier, c *transport.V3Transport) (int64, error) {
	opCtx := transport.FetchSpec{
		Pagination: transport.NotRequired,
		Resource:   fmt.Sprintf("/applications/%s/packageKeys", appId.ApplicationId),
		AppContext: "application package key",
	}

	return c.Count(ctx, opCtx)
}

func GetFullApplication(ctx context.Context, id masherytypes.ApplicationIdentifier, c *transport.V3Transport) (*masherytypes.Application, error) {
	qs := url.Values{
		"fields": {strings.Join(applicationDeepFields, ",")},
	}

	return httpToApplication(ctx, id, qs, c)
}

func httpToApplication(ctx context.Context, appId masherytypes.ApplicationIdentifier, qs url.Values, c *transport.V3Transport) (*masherytypes.Application, error) {
	rv, err := c.GetObject(ctx, transport.FetchSpec{
		Resource:       fmt.Sprintf("/applications/%s", appId),
		Query:          qs,
		AppContext:     "application",
		ResponseParser: masherytypes.ParseMasheryApplication,
	})

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(masherytypes.Application)
		return &retServ, nil
	}
}

// CreateApplication Create a new service.
func CreateApplication(ctx context.Context, app masherytypes.Application, memberId masherytypes.MemberIdentifier, c *transport.V3Transport) (*masherytypes.Application, error) {
	rawResp, err := c.CreateObject(ctx, app, transport.FetchSpec{
		Resource:       fmt.Sprintf("/members/%s/applications", memberId.MemberId),
		AppContext:     "application",
		ResponseParser: masherytypes.ParseMasheryApplication,
	})

	if err == nil {
		rv, _ := rawResp.(masherytypes.Application)
		return &rv, nil
	} else {
		return nil, err
	}
}

// UpdateApplication update an application
func UpdateApplication(ctx context.Context, app masherytypes.Application, c *transport.V3Transport) (*masherytypes.Application, error) {
	if app.Id == "" {
		return nil, errors.New("illegal argument: member Id must be set and not nil")
	}

	opContext := transport.FetchSpec{
		Resource:       fmt.Sprintf("/applications/%s", app.Id),
		AppContext:     "application",
		ResponseParser: masherytypes.ParseMasheryMember,
	}

	if d, err := c.UpdateObject(ctx, app, opContext); err == nil {
		rv, _ := d.(masherytypes.Application)
		return &rv, nil
	} else {
		return nil, err
	}
}

func DeleteApplication(ctx context.Context, appId masherytypes.ApplicationIdentifier, c *transport.V3Transport) error {
	opContext := transport.FetchSpec{
		Resource:       fmt.Sprintf("/applications/%s", appId.ApplicationId),
		AppContext:     "application",
		ResponseParser: masherytypes.ParseMasheryMember,
	}

	return c.DeleteObject(ctx, opContext)
}

func CountApplicationsOfMember(ctx context.Context, memberId masherytypes.MemberIdentifier, c *transport.V3Transport) (int64, error) {
	opCtx := transport.FetchSpec{
		Pagination: transport.NotRequired,
		Resource:   fmt.Sprintf("/members/%s/applications", memberId.MemberId),
		AppContext: "member's applications",
	}

	return c.Count(ctx, opCtx)
}

func ListApplications(ctx context.Context, c *transport.V3Transport) ([]masherytypes.Application, error) {
	opCtx := transport.FetchSpec{
		Pagination:     transport.PerPage,
		Resource:       "/applications",
		Query:          nil,
		AppContext:     "all applications",
		ResponseParser: masherytypes.ParseMasheryApplicationArray,
	}

	if d, err := c.FetchAll(ctx, opCtx); err != nil {
		return []masherytypes.Application{}, err
	} else {
		// Convert individual fetches into the array of elements
		var rv []masherytypes.Application
		for _, raw := range d {
			ms, ok := raw.([]masherytypes.Application)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}
