package v3client

import (
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
)

type OrganizationCRUDDecorator struct {
}

func (o OrganizationCRUDDecorator) ValueSupplier() transport.Supplier[masherytypes.Organization] {
	return func() masherytypes.Organization {
		return masherytypes.Organization{}
	}
}

func (o OrganizationCRUDDecorator) ValueArraySupplier() transport.Supplier[[]masherytypes.Organization] {
	return func() []masherytypes.Organization {
		return []masherytypes.Organization{}
	}
}

var organizationCRUDDecorator *GenericCRUDDecorator[int, string, masherytypes.Organization]
var organizationCRUD *GenericCRUD[int, string, masherytypes.Organization]

func init() {
	organizationCRUDDecorator = &GenericCRUDDecorator[int, string, masherytypes.Organization]{
		ValueSupplier:      func() masherytypes.Organization { return masherytypes.Organization{} },
		ValueArraySupplier: func() []masherytypes.Organization { return []masherytypes.Organization{} },
		ResourceFor: func(ident string) (string, error) {
			return fmt.Sprintf("/organizations/%s", ident), nil
		},
		ResourceForUpsert: func(t masherytypes.Organization) (string, error) {
			if len(t.Id) > 0 {
				return fmt.Sprintf("/organizations/%s", t.Id), nil
			}

			return "", errors.New("insufficient identification")
		},
		ResourceForParent: func(ident int) (string, error) {
			return "/organizations", nil
		},
		Pagination: transport.PerPage,
	}
	organizationCRUD = NewCRUD[int, string, masherytypes.Organization](
		"organization",
		organizationCRUDDecorator,
	)
}
