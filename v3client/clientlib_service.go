package v3client

import (
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
)

var serviceCRUDDecorator *GenericCRUDDecorator[int, masherytypes.ServiceIdentifier, masherytypes.Service]
var serviceCRUD *GenericCRUD[int, masherytypes.ServiceIdentifier, masherytypes.Service]

func init() {
	serviceCRUDDecorator = &GenericCRUDDecorator[int, masherytypes.ServiceIdentifier, masherytypes.Service]{
		ValueSupplier:      func() masherytypes.Service { return masherytypes.Service{} },
		ValueArraySupplier: func() []masherytypes.Service { return []masherytypes.Service{} },
		ResourceFor: func(ident masherytypes.ServiceIdentifier) (string, error) {
			if len(ident.ServiceId) == 0 {
				return "", errors.New("insufficient identifier")
			}
			return fmt.Sprintf("/services/%s", ident.ServiceId), nil
		},

		ResourceForUpsert: func(t masherytypes.Service) (string, error) {
			if len(t.Id) > 0 {
				return fmt.Sprintf("/services/%s", t.Id), nil
			}

			return "", errors.New("insufficient identification")
		},
		ResourceForParent: func(ident int) (string, error) {
			return "/services", nil
		},
		DefaultFields: MasheryServiceFields,
		Pagination:    transport.PerItem,
	}
	serviceCRUD = NewCRUD[int, masherytypes.ServiceIdentifier, masherytypes.Service]("service", serviceCRUDDecorator)
}
