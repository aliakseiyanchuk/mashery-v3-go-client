package v3client

import (
	"context"
	"errors"
	"fmt"
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

func GetApplication(ctx context.Context, appId string, c *HttpTransport) (*MasheryApplication, error) {
	qs := url.Values{
		"fields": {strings.Join(applicationFields, ",")},
	}

	return c.httpToApplication(ctx, appId, qs)
}

func GetApplicationPackageKeys(ctx context.Context, appId string, c *HttpTransport) ([]MasheryPackageKey, error) {
	qs := url.Values{
		"fields": {strings.Join(packageKeyFields, ",")},
	}

	opCtx := FetchSpec{
		Pagination:     PerPage,
		Resource:       fmt.Sprintf("/applications/%s/packageKeys", appId),
		Query:          qs,
		AppContext:     "application package key",
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

func (c *HttpTransport) CountApplicationPackageKeys(ctx context.Context, appId string) (int64, error) {
	opCtx := FetchSpec{
		Pagination: NotRequired,
		Resource:   fmt.Sprintf("/applications/%s/packageKeys", appId),
		AppContext: "application package key",
	}

	return c.count(ctx, opCtx)
}

func (c *HttpTransport) GetFullApplication(ctx context.Context, id string) (*MasheryApplication, error) {
	qs := url.Values{
		"fields": {strings.Join(applicationDeepFields, ",")},
	}

	return c.httpToApplication(ctx, id, qs)
}

func (c *HttpTransport) httpToApplication(ctx context.Context, appId string, qs url.Values) (*MasheryApplication, error) {
	rv, err := c.getObject(ctx, FetchSpec{
		Resource:       fmt.Sprintf("/applications/%s", appId),
		Query:          qs,
		AppContext:     "application",
		ResponseParser: ParseMasheryApplication,
	})

	if err != nil {
		return nil, err
	} else {
		retServ, _ := rv.(MasheryApplication)
		return &retServ, nil
	}
}

// Create a new service.
func (c *HttpTransport) CreateApplication(ctx context.Context, memberId string, member MasheryApplication) (*MasheryApplication, error) {
	rawResp, err := c.createObject(ctx, member, FetchSpec{
		Resource:       fmt.Sprintf("/members/%s/applications", memberId),
		AppContext:     "application",
		ResponseParser: ParseMasheryApplication,
	})

	if err == nil {
		rv, _ := rawResp.(MasheryApplication)
		return &rv, nil
	} else {
		return nil, err
	}
}

// Create a new service.
func (c *HttpTransport) UpdateApplication(ctx context.Context, app MasheryApplication) (*MasheryApplication, error) {
	if app.Id == "" {
		return nil, errors.New("illegal argument: member Id must be set and not nil")
	}

	opContext := FetchSpec{
		Resource:       fmt.Sprintf("/applications/%s", app.Id),
		AppContext:     "application",
		ResponseParser: ParseMasheryMember,
	}

	if d, err := c.updateObject(ctx, app, opContext); err == nil {
		rv, _ := d.(MasheryApplication)
		return &rv, nil
	} else {
		return nil, err
	}
}

func (c *HttpTransport) DeleteApplication(ctx context.Context, appId string) error {
	opContext := FetchSpec{
		Resource:       fmt.Sprintf("/applications/%s", appId),
		AppContext:     "application",
		ResponseParser: ParseMasheryMember,
	}

	return c.deleteObject(ctx, opContext)
}

func (c *HttpTransport) CountApplicationsOfMember(ctx context.Context, memberId string) (int64, error) {
	opCtx := FetchSpec{
		Pagination: NotRequired,
		Resource:   fmt.Sprintf("/members/%s/applications", memberId),
		AppContext: "member's applications",
	}

	return c.count(ctx, opCtx)
}

func (c *HttpTransport) ListApplications(ctx context.Context) ([]MasheryApplication, error) {
	opCtx := FetchSpec{
		Pagination:     PerPage,
		Resource:       "/applications",
		Query:          nil,
		AppContext:     "all applications",
		ResponseParser: ParseMasheryApplicationArray,
	}

	if d, err := c.fetchAll(ctx, opCtx); err != nil {
		return []MasheryApplication{}, err
	} else {
		// Convert individual fetches into the array of elements
		var rv []MasheryApplication
		for _, raw := range d {
			ms, ok := raw.([]MasheryApplication)
			if ok {
				rv = append(rv, ms...)
			}
		}

		return rv, nil
	}
}
