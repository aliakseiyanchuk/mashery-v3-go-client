package v3client

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
)

var applicationFields = []string{
	"id", "name", "created", "updated", "username", "description",
	"type", "commercial", "ads", "tags",
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

		UpsertCleaner: func(m *masherytypes.Application) {
			m.Id = ""
			m.Created = nil
			m.Updated = nil
		},

		Pagination: transport.PerPage,
	}

	applicationCRUD = NewCRUD[masherytypes.MemberIdentifier, masherytypes.ApplicationIdentifier, masherytypes.Application]("application", applicationCRUDDecorator)
}

var applicationRawCRUDDecorator *GenericCRUDDecorator[int, masherytypes.ApplicationIdentifier, map[string]interface{}]
var applicationRawCRUD *GenericCRUD[int, masherytypes.ApplicationIdentifier, map[string]interface{}]

func init() {
	applicationRawCRUDDecorator = &GenericCRUDDecorator[int, masherytypes.ApplicationIdentifier, map[string]interface{}]{
		ValueSupplier: func() map[string]interface{} {
			return map[string]interface{}{}
		},
		ResourceFor: func(ident masherytypes.ApplicationIdentifier) (string, error) {
			if len(ident.ApplicationId) == 0 {
				return "", errors.New("application identifier cannot be empty")
			} else {
				return fmt.Sprintf("/applications/%s", ident.ApplicationId), nil
			}
		},
		ResourceForUpsert: func(t map[string]interface{}) (string, error) {
			id := t["_id"]
			if id != nil {
				return fmt.Sprintf("/applications/%s", id), nil
			} else {
				return "", errors.New("unsupported identification for upsert")
			}
		},

		UpsertCleaner: func(m *map[string]interface{}) {
			delete(*m, "_id")
		},

		Pagination: transport.PerPage,
	}

	applicationRawCRUD = NewCRUD[int, masherytypes.ApplicationIdentifier, map[string]interface{}]("application extended attributes", applicationRawCRUDDecorator)
}

func GetApplicationExtendedAttributes(ctx context.Context, appId masherytypes.ApplicationIdentifier, transport *transport.HttpTransport) (map[string]string, error) {
	if rv, exists, err := applicationRawCRUD.Get(ctx, appId, transport); err != nil {
		return map[string]string{}, err
	} else if !exists {
		return map[string]string{}, errors.New("application should exist before this call can succeed")
	} else {
		return filterApplicationExtendedAttributes(rv)
	}
}

func filterApplicationExtendedAttributes(rv map[string]interface{}) (map[string]string, error) {
	filteredRV := map[string]string{}
	for k, v := range rv {
		if isApplicationSystemAttr(k) {
			continue
		} else if vStr, ok := v.(string); ok {
			filteredRV[k] = vStr
		}
	}

	return filteredRV, nil
}

func isApplicationSystemAttr(s string) bool {
	switch s {
	case "id":
		fallthrough
	case "ads":
		fallthrough
	case "adsSystem":
		fallthrough
	case "created":
		fallthrough
	case "updated":
		fallthrough
	case "description":
		fallthrough
	case "externalId":
		fallthrough
	case "howDidYouHear":
		fallthrough
	case "name":
		fallthrough
	case "notes":
		fallthrough
	case "oauthRedirectUri":
		fallthrough
	case "preferredOutput":
		fallthrough
	case "preferredProtocol":
		fallthrough
	case "status":
		fallthrough
	case "tags":
		fallthrough
	case "type":
		fallthrough
	case "uri":
		fallthrough
	case "usageModel":
		fallthrough
	case "username":
		return true
	default:
		return false
	}
}

func UpdateApplicationExtendedAttributes(ctx context.Context, appId masherytypes.ApplicationIdentifier, params map[string]string, transport *transport.HttpTransport) (map[string]string, error) {
	callParams := map[string]interface{}{
		"_id": appId.ApplicationId,
	}
	for k, v := range params {
		if k == "_id" {
			return map[string]string{}, errors.New("illegal argument: parameters cannot contain _id attribute")
		} else if isApplicationSystemAttr(k) {
			return map[string]string{}, errors.New("illegal argument: parameters cannot contain application system attribute")
		}

		callParams[k] = v
	}

	if rv, err := applicationRawCRUD.Update(ctx, callParams, transport); err != nil {
		return map[string]string{}, err
	} else {
		return filterApplicationExtendedAttributes(rv)
	}
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
