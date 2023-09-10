package v3client

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
)

var publicDomainsCRUDDecorator *GenericCRUDDecorator[int, int, masherytypes.DomainAddress]
var systemDomainsCRUDDecorator *GenericCRUDDecorator[int, int, masherytypes.DomainAddress]
var publicDomainsCRUD *GenericCRUD[int, int, masherytypes.DomainAddress]
var systemDomainsCRUD *GenericCRUD[int, int, masherytypes.DomainAddress]

func init() {
	publicDomainsCRUDDecorator = &GenericCRUDDecorator[int, int, masherytypes.DomainAddress]{
		ValueArraySupplier: func() []masherytypes.DomainAddress { return []masherytypes.DomainAddress{} },
		ResourceFor:        func(i int) (string, error) { return "/domains/public/hostnames", nil },
		ResourceForParent:  func(i int) (string, error) { return "/domains/public/hostnames", nil },
		Pagination:         transport.PerPage,
	}

	systemDomainsCRUDDecorator = &GenericCRUDDecorator[int, int, masherytypes.DomainAddress]{
		ValueArraySupplier: func() []masherytypes.DomainAddress { return []masherytypes.DomainAddress{} },
		ResourceFor:        func(i int) (string, error) { return "/domains/system/hostnames", nil },
		ResourceForParent:  func(i int) (string, error) { return "/domains/system/hostnames", nil },
		Pagination:         transport.PerPage,
	}
	publicDomainsCRUD = NewCRUD[int, int, masherytypes.DomainAddress]("unique public domains", publicDomainsCRUDDecorator)
	systemDomainsCRUD = NewCRUD[int, int, masherytypes.DomainAddress]("unique system domains", systemDomainsCRUDDecorator)
}
