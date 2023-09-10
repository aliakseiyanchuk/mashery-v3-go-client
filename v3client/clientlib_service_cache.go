package v3client

import (
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
)

var serviceCacheCRUDDecorator *GenericCRUDDecorator[masherytypes.ServiceIdentifier, masherytypes.ServiceIdentifier, masherytypes.ServiceCache]
var serviceCacheCRUD *GenericCRUD[masherytypes.ServiceIdentifier, masherytypes.ServiceIdentifier, masherytypes.ServiceCache]

func init() {
	serviceCacheCRUDDecorator = &GenericCRUDDecorator[masherytypes.ServiceIdentifier, masherytypes.ServiceIdentifier, masherytypes.ServiceCache]{
		ValueSupplier:      func() masherytypes.ServiceCache { return masherytypes.ServiceCache{} },
		ValueArraySupplier: func() []masherytypes.ServiceCache { return []masherytypes.ServiceCache{} },

		AcceptIdentFrom: func(t1 masherytypes.ServiceCache, t2 *masherytypes.ServiceCache) {
			t2.ParentServiceId = t1.ParentServiceId
		},
		AcceptObjectIdent: func(t1 masherytypes.ServiceIdentifier, t2 *masherytypes.ServiceCache) {
			t2.ParentServiceId = t1
		},
		AcceptParentIdent: func(t1 masherytypes.ServiceIdentifier, t2 *masherytypes.ServiceCache) {
			t2.ParentServiceId = t1
		},

		ResourceFor: func(ident masherytypes.ServiceIdentifier) (string, error) {
			if len(ident.ServiceId) == 0 {
				return "", errors.New("insufficient identifier")
			}
			return fmt.Sprintf("/services/%s/cache", ident.ServiceId), nil
		},

		ResourceForUpsert: func(t masherytypes.ServiceCache) (string, error) {
			if len(t.ParentServiceId.ServiceId) > 0 {
				return fmt.Sprintf("/services/%s/cache", t.ParentServiceId.ServiceId), nil
			}

			return "", errors.New("insufficient identification")
		},

		ResourceForParent: func(ident masherytypes.ServiceIdentifier) (string, error) {
			if len(ident.ServiceId) > 0 {
				return fmt.Sprintf("/services/%s/cache", ident.ServiceId), nil
			}

			return fmt.Sprintf("/services/%s/cache", ident.ServiceId), nil
		},

		Pagination: transport.NotRequired,
	}
	serviceCacheCRUD = NewCRUD[masherytypes.ServiceIdentifier, masherytypes.ServiceIdentifier, masherytypes.ServiceCache](
		"service cache",
		serviceCacheCRUDDecorator,
	)
}
