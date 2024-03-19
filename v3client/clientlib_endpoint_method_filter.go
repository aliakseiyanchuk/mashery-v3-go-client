package v3client

import (
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
)

var endpointMethodFilterCRUDDecorator *GenericCRUDDecorator[masherytypes.ServiceEndpointMethodIdentifier, masherytypes.ServiceEndpointMethodFilterIdentifier, masherytypes.ServiceEndpointMethodFilter]
var endpointMethodFilterCRUD *GenericCRUD[masherytypes.ServiceEndpointMethodIdentifier, masherytypes.ServiceEndpointMethodFilterIdentifier, masherytypes.ServiceEndpointMethodFilter]

func init() {
	endpointMethodFilterCRUDDecorator = &GenericCRUDDecorator[masherytypes.ServiceEndpointMethodIdentifier, masherytypes.ServiceEndpointMethodFilterIdentifier, masherytypes.ServiceEndpointMethodFilter]{
		ValueSupplier:      func() masherytypes.ServiceEndpointMethodFilter { return masherytypes.ServiceEndpointMethodFilter{} },
		ValueArraySupplier: func() []masherytypes.ServiceEndpointMethodFilter { return []masherytypes.ServiceEndpointMethodFilter{} },

		AcceptObjectIdent: func(t1 masherytypes.ServiceEndpointMethodFilterIdentifier, t2 *masherytypes.ServiceEndpointMethodFilter) {
			t2.ServiceEndpointMethod = t1.ServiceEndpointMethodIdentifier
		},
		AcceptIdentFrom: func(t1 masherytypes.ServiceEndpointMethodFilter, t2 *masherytypes.ServiceEndpointMethodFilter) {
			t2.ServiceEndpointMethod = t1.ServiceEndpointMethod
		},
		AcceptParentIdent: func(t1 masherytypes.ServiceEndpointMethodIdentifier, t2 *masherytypes.ServiceEndpointMethodFilter) {
			t2.ServiceEndpointMethod = t1
		},

		ResourceFor: func(ident masherytypes.ServiceEndpointMethodFilterIdentifier) (string, error) {
			return fmt.Sprintf("/services/%s/endpoints/%s/methods/%s/responseFilters/%s",
				ident.ServiceId,
				ident.EndpointId,
				ident.MethodId,
				ident.FilterId), nil
		},
		ResourceForUpsert: func(t masherytypes.ServiceEndpointMethodFilter) (string, error) {
			if len(t.Id) > 0 {
				return fmt.Sprintf("/services/%s/endpoints/%s/methods/%s/responseFilters/%s",
					t.ServiceEndpointMethod.ServiceId,
					t.ServiceEndpointMethod.EndpointId,
					t.ServiceEndpointMethod.MethodId,
					t.Id), nil
			}

			return "", errors.New("insufficient identification")
		},
		ResourceForParent: func(ident masherytypes.ServiceEndpointMethodIdentifier) (string, error) {
			return fmt.Sprintf("/services/%s/endpoints/%s/methods/%s/responseFilters", ident.ServiceId, ident.EndpointId, ident.MethodId), nil
		},
		DefaultFields: MasheryResponseFilterFields,
		Pagination:    transport.PerPage,
	}
	endpointMethodFilterCRUD = NewCRUD[masherytypes.ServiceEndpointMethodIdentifier, masherytypes.ServiceEndpointMethodFilterIdentifier, masherytypes.ServiceEndpointMethodFilter](
		"endpoint methods filters",
		endpointMethodFilterCRUDDecorator,
	)
}
