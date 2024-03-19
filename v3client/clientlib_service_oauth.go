package v3client

import (
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
)

var serviceOAuthCRUDDecorator *GenericCRUDDecorator[masherytypes.ServiceIdentifier, masherytypes.ServiceIdentifier, masherytypes.MasheryOAuth]
var serviceOAuthCRUD *GenericCRUD[masherytypes.ServiceIdentifier, masherytypes.ServiceIdentifier, masherytypes.MasheryOAuth]

func init() {
	serviceOAuthCRUDDecorator = &GenericCRUDDecorator[masherytypes.ServiceIdentifier, masherytypes.ServiceIdentifier, masherytypes.MasheryOAuth]{
		ValueSupplier:      func() masherytypes.MasheryOAuth { return masherytypes.MasheryOAuth{} },
		ValueArraySupplier: func() []masherytypes.MasheryOAuth { return []masherytypes.MasheryOAuth{} },

		AcceptParentIdent: func(t1 masherytypes.ServiceIdentifier, t2 *masherytypes.MasheryOAuth) {
			t2.ParentService = t1
		},
		AcceptObjectIdent: func(t1 masherytypes.ServiceIdentifier, t2 *masherytypes.MasheryOAuth) {
			t2.ParentService = t1
		},
		AcceptIdentFrom: func(t1 masherytypes.MasheryOAuth, t2 *masherytypes.MasheryOAuth) {
			t2.ParentService = t1.ParentService
		},

		ResourceFor: func(ident masherytypes.ServiceIdentifier) (string, error) {
			return fmt.Sprintf("/services/%s/securityProfile/oauth", ident.ServiceId), nil
		},

		ResourceForUpsert: func(t masherytypes.MasheryOAuth) (string, error) {
			if len(t.ParentService.ServiceId) > 0 {
				return fmt.Sprintf("/services/%s/securityProfile/oauth", t.ParentService.ServiceId), nil
			}

			return "", errors.New("insufficient identification")
		},

		ResourceForParent: func(ident masherytypes.ServiceIdentifier) (string, error) {
			if len(ident.ServiceId) > 0 {
				return fmt.Sprintf("/services/%s/securityProfile/oauth", ident.ServiceId), nil
			}

			return "", errors.New("insufficient identification")
		},
		Pagination: transport.NotRequired,
	}
	serviceOAuthCRUD = NewCRUD[masherytypes.ServiceIdentifier, masherytypes.ServiceIdentifier, masherytypes.MasheryOAuth](
		"service security profile",
		serviceOAuthCRUDDecorator,
	)
}
