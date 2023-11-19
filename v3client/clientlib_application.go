package v3client

import (
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
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

// -------------------------------------------------------------------------------------------
// Application CRUD
// -------------------------------------------------------------------------------------------

func applicationValueSupplier() masherytypes.Application {
	return masherytypes.Application{}
}

func applicationArrayValueSupplier() []masherytypes.Application {
	return []masherytypes.Application{}
}

var applicationCRUDDecorator *GenericCRUDDecorator[masherytypes.MemberIdentifier, masherytypes.ApplicationIdentifier, masherytypes.Application]
var applicationCRUD *GenericCRUD[masherytypes.MemberIdentifier, masherytypes.ApplicationIdentifier, masherytypes.Application]

func init() {
	applicationCRUDDecorator = &GenericCRUDDecorator[masherytypes.MemberIdentifier, masherytypes.ApplicationIdentifier, masherytypes.Application]{
		ValueSupplier:      applicationValueSupplier,
		ValueArraySupplier: applicationArrayValueSupplier,
		ResourceFor: func(ident masherytypes.ApplicationIdentifier) (string, error) {
			if len(ident.ApplicationId) == 0 {
				return "", errors.New("application identifier cannot be empty")
			} else {
				return fmt.Sprintf("/applications/%s", ident.ApplicationId), nil
			}
		},
		ResourceForUpsert: func(t masherytypes.Application) (string, error) {
			if len(t.Id) > 0 {
				return fmt.Sprintf("/applications/%s", t.Id), nil
			} else {
				return "", errors.New("unsupported identification for upsert")
			}
		},
		ResourceForParent: func(ident masherytypes.MemberIdentifier) (string, error) {
			if len(ident.MemberId) == 0 {
				return "/applications", nil
			} else if len(ident.MemberId) > 0 {
				return fmt.Sprintf("/members/%s/applications", ident.MemberId), nil
			} else {
				return "", errors.New("unsupported identification for upsert")
			}
		},
		DefaultFields: applicationFields,
		GetFields:     DefaultGetFieldsFromContext(applicationFields),
		Pagination:    transport.PerPage,
	}

	applicationCRUD = NewCRUD[masherytypes.MemberIdentifier, masherytypes.ApplicationIdentifier, masherytypes.Application]("application", applicationCRUDDecorator)
}

// -------
// Application package keys
// -------

var applicationPackageKeysCRUDDecorator *GenericCRUDDecorator[masherytypes.ApplicationIdentifier, masherytypes.PackageKeyIdentifier, masherytypes.PackageKey]

var applicationPackageKeyCRUD *GenericCRUD[masherytypes.ApplicationIdentifier, masherytypes.PackageKeyIdentifier, masherytypes.PackageKey]

func init() {
	applicationPackageKeysCRUDDecorator = &GenericCRUDDecorator[masherytypes.ApplicationIdentifier, masherytypes.PackageKeyIdentifier, masherytypes.PackageKey]{
		ValueSupplier:      func() masherytypes.PackageKey { return masherytypes.PackageKey{} },
		ValueArraySupplier: func() []masherytypes.PackageKey { return []masherytypes.PackageKey{} },

		ResourceForParent: func(ident masherytypes.ApplicationIdentifier) (string, error) {
			if len(ident.ApplicationId) > 0 {
				return fmt.Sprintf("/applications/%s/packageKeys", ident.ApplicationId), nil
			} else {
				return "", errors.New("application identifier is required")
			}
		},
		DefaultFields: packageKeyFields,
		Pagination:    transport.PerPage,
	}

	applicationPackageKeyCRUD = NewCRUD[masherytypes.ApplicationIdentifier,
		masherytypes.PackageKeyIdentifier,
		masherytypes.PackageKey](
		"application package key",
		applicationPackageKeysCRUDDecorator,
	)
}
