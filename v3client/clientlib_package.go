package v3client

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
)

var packageCRUDDecorator *GenericCRUDDecorator[int, masherytypes.PackageIdentifier, masherytypes.Package]
var packageCRUD *GenericCRUD[int, masherytypes.PackageIdentifier, masherytypes.Package]

func init() {
	packageCRUDDecorator = &GenericCRUDDecorator[int, masherytypes.PackageIdentifier, masherytypes.Package]{
		ValueSupplier:      func() masherytypes.Package { return masherytypes.Package{} },
		ValueArraySupplier: func() []masherytypes.Package { return []masherytypes.Package{} },
		ResourceFor: func(ident masherytypes.PackageIdentifier) (string, error) {
			if len(ident.PackageId) == 0 {
				return "", errors.New("insufficient identifier")
			}
			return fmt.Sprintf("/packages/%s", ident.PackageId), nil
		},
		ResourceForUpsert: func(t masherytypes.Package) (string, error) {
			if len(t.Id) > 0 {
				return fmt.Sprintf("/packages/%s", t.Id), nil
			}
			return "", errors.New("insufficient identification")
		},
		ResourceForParent: func(_ int) (string, error) {
			return "/packages", nil
		},
		DefaultFields: MasheryPackageFields,
		Pagination:    transport.PerItem,
	}
	packageCRUD = NewCRUD[int, masherytypes.PackageIdentifier, masherytypes.Package]("package", packageCRUDDecorator)
}

type orgPutSchema struct {
	Organization masherytypes.NilAddressableOrganization `json:"organization"`
}

func ResetPackageOwnership(ctx context.Context, pack masherytypes.PackageIdentifier, c *transport.HttpTransport) (masherytypes.Package, error) {
	if len(pack.PackageId) == 0 {
		return masherytypes.Package{}, errors.New("illegal argument: package Id must be set")
	}

	dat := orgPutSchema{
		Organization: masherytypes.NilAddressableOrganization{
			Id:     nil,
			Parent: nil,
			Name:   "Area Level",
		},
	}

	orgPutSchemaReq := transport.ObjectUpsertSpecBuilder[orgPutSchema]{}
	orgPutSchemaReq.
		WithUpsert(dat).
		WithResource("/packages/%s", pack.PackageId).
		WithAppContext("package ownership reset")

	if _, updateErr := transport.UpdateObject(ctx, orgPutSchemaReq.Build(), c); updateErr == nil {
		get, _, getErr := packageCRUD.Get(ctx, pack, c)
		return get, getErr
	} else {
		return masherytypes.Package{}, updateErr
	}
}
