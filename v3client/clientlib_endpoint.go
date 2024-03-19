package v3client

import (
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
)

// ---------------
// Endpoint CRUD

var endpointCRUDDecorator *GenericCRUDDecorator[masherytypes.ServiceIdentifier, masherytypes.ServiceEndpointIdentifier, masherytypes.Endpoint]
var endpointCRUD *GenericCRUD[masherytypes.ServiceIdentifier, masherytypes.ServiceEndpointIdentifier, masherytypes.Endpoint]

func init() {
	endpointCRUDDecorator = &GenericCRUDDecorator[masherytypes.ServiceIdentifier, masherytypes.ServiceEndpointIdentifier, masherytypes.Endpoint]{
		ValueSupplier:      func() masherytypes.Endpoint { return masherytypes.Endpoint{} },
		ValueArraySupplier: func() []masherytypes.Endpoint { return []masherytypes.Endpoint{} },

		AcceptParentIdent: func(t1 masherytypes.ServiceIdentifier, t2 *masherytypes.Endpoint) {
			t2.ParentServiceId = t1
		},
		AcceptIdentFrom: func(t1 masherytypes.Endpoint, t2 *masherytypes.Endpoint) { t2.ParentServiceId = t1.ParentServiceId },

		ResourceFor: func(ident masherytypes.ServiceEndpointIdentifier) (string, error) {
			if len(ident.ServiceId) == 0 || len(ident.EndpointId) == 0 {
				return "", errors.New("missing identifier")
			} else {
				return fmt.Sprintf("/services/%s/endpoints/%s", ident.ServiceId, ident.EndpointId), nil
			}
		},
		ResourceForUpsert: func(t masherytypes.Endpoint) (string, error) {
			if len(t.Id) > 0 {
				return fmt.Sprintf("/services/%s/endpoints/%s", t.ParentServiceId.ServiceId, t.Id), nil
			}

			return "", errors.New("unresolvable identification")
		},
		ResourceForParent: func(ident masherytypes.ServiceIdentifier) (string, error) {
			if len(ident.ServiceId) > 0 {
				return fmt.Sprintf("/services/%s/endpoints", ident.ServiceId), nil
			}

			return "", errors.New("unresolvable identification")
		},
		DefaultFields: MasheryEndpointFields,
		Pagination:    transport.PerPage,
	}

	endpointCRUD = NewCRUD[masherytypes.ServiceIdentifier, masherytypes.ServiceEndpointIdentifier, masherytypes.Endpoint](
		"service endpoint",
		endpointCRUDDecorator,
	)
}
