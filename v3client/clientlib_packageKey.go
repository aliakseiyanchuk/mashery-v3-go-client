package v3client

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/errwrap"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"net/url"
	"strings"
)

var packageKeyCRUDDecorator *GenericCRUDDecorator[int, masherytypes.PackageKeyIdentifier, masherytypes.PackageKey]
var packageKeyCRUD *GenericCRUD[int, masherytypes.PackageKeyIdentifier, masherytypes.PackageKey]

func init() {
	packageKeyCRUDDecorator = &GenericCRUDDecorator[int, masherytypes.PackageKeyIdentifier, masherytypes.PackageKey]{
		ValueSupplier:      func() masherytypes.PackageKey { return masherytypes.PackageKey{} },
		ValueArraySupplier: func() []masherytypes.PackageKey { return []masherytypes.PackageKey{} },
		ResourceFor: func(ident masherytypes.PackageKeyIdentifier) (string, error) {
			return fmt.Sprintf("/packageKeys/%s", ident.PackageKeyId), nil
		},
		ResourceForUpsert: func(t masherytypes.PackageKey) (string, error) {
			if len(t.Id) > 0 {
				return fmt.Sprintf("/packageKeys/%s", t.Id), nil
			}
			return "", errors.New("insufficient identification")
		},
		ResourceForParent: func(ident int) (string, error) {
			return "/packageKeys", nil
		},
		DefaultFields: MasheryPackageKeyFields,
		Pagination:    transport.PerPage,
	}
	packageKeyCRUD = NewCRUD[int, masherytypes.PackageKeyIdentifier, masherytypes.PackageKey]("package key", packageKeyCRUDDecorator)
}

// TODO: Convert this method
// CreatePackageKey Create a new service.
func CreatePackageKey(ctx context.Context, appId masherytypes.ApplicationIdentifier, packageKey masherytypes.PackageKey, c *transport.HttpTransport) (masherytypes.PackageKey, error) {
	if !packageKey.LinksPackageAndPlan() {
		return masherytypes.PackageKey{}, &errwrap.WrappedError{
			Context: "create package key",
			Cause:   errors.New("package key must supply associated package and plan"),
		}
	}

	builder := transport.ObjectUpsertSpecBuilder[masherytypes.PackageKey]{}
	builder.
		WithUpsert(packageKey).
		WithValueFactory(func() masherytypes.PackageKey {
			return masherytypes.PackageKey{}
		}).
		WithResource("/applications/%s/packageKeys", appId.ApplicationId).
		WithQuery(url.Values{
			"fields": {strings.Join(MasheryPackageKeyFields, ",")},
		}).
		WithAppContext("package key")

	return transport.CreateObject(ctx, builder.Build(), c)
}

// UpdatePackageKey Create a new service.
