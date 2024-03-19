package v3client

import (
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
)

var roleCRUDDecorator *GenericCRUDDecorator[int, string, masherytypes.Role]
var roleCRUD *GenericCRUD[int, string, masherytypes.Role]

func init() {
	roleCRUDDecorator = &GenericCRUDDecorator[int, string, masherytypes.Role]{
		ValueSupplier:      func() masherytypes.Role { return masherytypes.Role{} },
		ValueArraySupplier: func() []masherytypes.Role { return []masherytypes.Role{} },
		ResourceFor:        func(ident string) (string, error) { return fmt.Sprintf("/roles/%s", ident), nil },
		ResourceForUpsert:  func(t masherytypes.Role) (string, error) { return fmt.Sprintf("/roles/%s", t.Id), nil },
		ResourceForParent:  func(_ int) (string, error) { return "/roles", nil },
		Pagination:         transport.PerPage,
	}
	roleCRUD = NewCRUD[int, string, masherytypes.Role]("role", roleCRUDDecorator)
}
