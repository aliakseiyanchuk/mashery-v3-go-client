package v3client

import (
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
)

var packageKeyCRUDDecorator *GenericCRUDDecorator[int, masherytypes.PackageKeyIdentifier, masherytypes.PackageKey]
var packageKeyCRUD *GenericCRUD[int, masherytypes.PackageKeyIdentifier, masherytypes.PackageKey]

func init() {
	packageKeyCRUDDecorator = &GenericCRUDDecorator[int, masherytypes.PackageKeyIdentifier, masherytypes.PackageKey]{
		ValueSupplier:      func() masherytypes.PackageKey { return masherytypes.PackageKey{} },
		ValueArraySupplier: func() []masherytypes.PackageKey { return []masherytypes.PackageKey{} },
		ResourceFor: func(ident masherytypes.PackageKeyIdentifier) (string, error) {
			return fmt.Sprintf("/packageKeys/%s", ident.PackageKeyId), nil
		},
		ResourceForUpsert: func(t masherytypes.PackageKey) (string, error) {
			if len(t.Id) > 0 {
				return fmt.Sprintf("/packageKeys/%s", t.Id), nil
			}
			return "", errors.New("insufficient identification")
		},

		UpsertCleaner: func(m *masherytypes.PackageKey) {
			m.Id = ""
			m.Created = nil
			m.Updated = nil
		},

		ResourceForParent: func(ident int) (string, error) {
			return "/packageKeys", nil
		},
		DefaultFields: MasheryPackageKeyFields,
		Pagination:    transport.PerPage,
	}
	packageKeyCRUD = NewCRUD[int, masherytypes.PackageKeyIdentifier, masherytypes.PackageKey]("package key", packageKeyCRUDDecorator)
}
