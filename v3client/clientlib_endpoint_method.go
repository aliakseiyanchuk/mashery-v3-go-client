package v3client

import (
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
)

var endpointMethodCRUDDecorator *GenericCRUDDecorator[masherytypes.ServiceEndpointIdentifier, masherytypes.ServiceEndpointMethodIdentifier, masherytypes.ServiceEndpointMethod]
var endpointMethodCRUD *GenericCRUD[masherytypes.ServiceEndpointIdentifier, masherytypes.ServiceEndpointMethodIdentifier, masherytypes.ServiceEndpointMethod]

func init() {
	endpointMethodCRUDDecorator = &GenericCRUDDecorator[masherytypes.ServiceEndpointIdentifier, masherytypes.ServiceEndpointMethodIdentifier, masherytypes.ServiceEndpointMethod]{
		ValueSupplier: func() masherytypes.ServiceEndpointMethod {
			return masherytypes.ServiceEndpointMethod{}
		},
		ValueArraySupplier: func() []masherytypes.ServiceEndpointMethod {
			return []masherytypes.ServiceEndpointMethod{}
		},

		AcceptParentIdent: func(t1 masherytypes.ServiceEndpointIdentifier, t2 *masherytypes.ServiceEndpointMethod) {
			t2.ParentEndpointId = t1
		},
		AcceptIdentFrom: func(t1 masherytypes.ServiceEndpointMethod, t2 *masherytypes.ServiceEndpointMethod) {
			t2.ParentEndpointId = t1.ParentEndpointId
		},
		AcceptObjectIdent: func(t1 masherytypes.ServiceEndpointMethodIdentifier, t2 *masherytypes.ServiceEndpointMethod) {
			t2.ParentEndpointId = t1.ServiceEndpointIdentifier
		},

		ResourceFor: func(ident masherytypes.ServiceEndpointMethodIdentifier) (string, error) {
			if len(ident.ServiceId) == 0 || len(ident.EndpointId) == 0 || len(ident.MethodId) == 0 {
				return "", errors.New("insufficient service endpoint method identification")
			}

			return fmt.Sprintf("/services/%s/endpoints/%s/methods/%s", ident.ServiceId, ident.EndpointId, ident.MethodId), nil
		},
		ResourceForUpsert: func(t masherytypes.ServiceEndpointMethod) (string, error) {
			if len(t.Id) > 0 {
				return fmt.Sprintf("/services/%s/endpoints/%s/methods/%s", t.ParentEndpointId.ServiceId, t.ParentEndpointId.EndpointId, t.Id), nil
			}
			return "", errors.New("insufficient identification")
		},
		ResourceForParent: func(ident masherytypes.ServiceEndpointIdentifier) (string, error) {
			return fmt.Sprintf("/services/%s/endpoints/%s/methods", ident.ServiceId, ident.EndpointId), nil
		},
		DefaultFields: MasheryMethodsFields,
		Pagination:    transport.PerPage,
	}

	endpointMethodCRUD = NewCRUD[masherytypes.ServiceEndpointIdentifier, masherytypes.ServiceEndpointMethodIdentifier, masherytypes.ServiceEndpointMethod](
		"endpoint method",
		endpointMethodCRUDDecorator,
	)
}
