package v3client

import (
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
)

var errorSetCRUDDecorator *GenericCRUDDecorator[masherytypes.ServiceIdentifier, masherytypes.ErrorSetIdentifier, masherytypes.ErrorSet]
var errorSetCRUD *GenericCRUD[masherytypes.ServiceIdentifier, masherytypes.ErrorSetIdentifier, masherytypes.ErrorSet]

func init() {
	errorSetCRUDDecorator = &GenericCRUDDecorator[masherytypes.ServiceIdentifier, masherytypes.ErrorSetIdentifier, masherytypes.ErrorSet]{
		ValueSupplier:      func() masherytypes.ErrorSet { return masherytypes.ErrorSet{} },
		ValueArraySupplier: func() []masherytypes.ErrorSet { return []masherytypes.ErrorSet{} },

		AcceptIdentFrom: func(t1 masherytypes.ErrorSet, t2 *masherytypes.ErrorSet) {
			t2.ParentServiceId = t1.ParentServiceId
		},
		AcceptObjectIdent: func(t1 masherytypes.ErrorSetIdentifier, t2 *masherytypes.ErrorSet) {
			t2.ParentServiceId = t1.ServiceIdentifier
		},
		AcceptParentIdent: func(t1 masherytypes.ServiceIdentifier, t2 *masherytypes.ErrorSet) {
			t2.ParentServiceId = t1
		},

		ResourceFor: func(ident masherytypes.ErrorSetIdentifier) (string, error) {
			if len(ident.ServiceId) == 0 || len(ident.ErrorSetId) == 0 {
				return "", errors.New("insufficient identification")
			}
			return fmt.Sprintf("/services/%s/errorSets/%s", ident.ServiceId, ident.ErrorSetId), nil
		},
		ResourceForUpsert: func(t masherytypes.ErrorSet) (string, error) {
			if len(t.Id) > 0 {
				if len(t.ParentServiceId.ServiceId) == 0 {
					return "", errors.New("insufficient identification")
				}
				return fmt.Sprintf("/services/%s/errorSets/%s", t.ParentServiceId.ServiceId, t.Id), nil
			}

			return "", errors.New("insufficient identification")
		},
		ResourceForParent: func(ident masherytypes.ServiceIdentifier) (string, error) {
			if len(ident.ServiceId) > 0 {
				return fmt.Sprintf("/services/%s/errorSets", ident.ServiceId), nil
			}

			return "", errors.New("insufficient identification")
		},
		DefaultFields: MasheryErrorSetFields,
		Pagination:    transport.PerPage,
	}
	errorSetCRUD = NewCRUD[masherytypes.ServiceIdentifier, masherytypes.ErrorSetIdentifier, masherytypes.ErrorSet]("error set", errorSetCRUDDecorator)
}

var errorSetMessageCRUDDecorator *GenericCRUDDecorator[masherytypes.ErrorSetIdentifier, masherytypes.ErrorSetMessageIdentifier, masherytypes.MasheryErrorMessage]
var errorSetMessageCRUD *GenericCRUD[masherytypes.ErrorSetIdentifier, masherytypes.ErrorSetMessageIdentifier, masherytypes.MasheryErrorMessage]

func init() {
	errorSetMessageCRUDDecorator = &GenericCRUDDecorator[masherytypes.ErrorSetIdentifier, masherytypes.ErrorSetMessageIdentifier, masherytypes.MasheryErrorMessage]{
		ValueSupplier:      func() masherytypes.MasheryErrorMessage { return masherytypes.MasheryErrorMessage{} },
		ValueArraySupplier: func() []masherytypes.MasheryErrorMessage { return []masherytypes.MasheryErrorMessage{} },

		AcceptParentIdent: func(t1 masherytypes.ErrorSetIdentifier, t2 *masherytypes.MasheryErrorMessage) {
			t2.ParentErrorSet = t1
		},
		AcceptObjectIdent: func(t1 masherytypes.ErrorSetMessageIdentifier, t2 *masherytypes.MasheryErrorMessage) {
			t2.ParentErrorSet = t1.ErrorSetIdentifier
		},
		AcceptIdentFrom: func(t1 masherytypes.MasheryErrorMessage, t2 *masherytypes.MasheryErrorMessage) {
			t2.ParentErrorSet = t1.ParentErrorSet
		},

		ResourceFor: func(ident masherytypes.ErrorSetMessageIdentifier) (string, error) {
			if len(ident.ServiceId) == 0 || len(ident.ErrorSetId) == 0 {
				return "", errors.New("insufficient identification")
			}
			return fmt.Sprintf("/services/%s/errorSets/%s", ident.ServiceId, ident.ErrorSetId), nil
		},
		ResourceForUpsert: func(t masherytypes.MasheryErrorMessage) (string, error) {
			if len(t.ParentErrorSet.ServiceId) == 0 || len(t.ParentErrorSet.ErrorSetId) == 0 || len(t.Id) == 0 {
				return "", errors.New("insufficient identification")
			}
			return fmt.Sprintf("/services/%s/errorSets/%s/errorMessages/%s", t.ParentErrorSet.ServiceId, t.ParentErrorSet.ErrorSetId, t.Id), nil
		},
		ResourceForParent: func(ident masherytypes.ErrorSetIdentifier) (string, error) {
			if len(ident.ServiceId) > 0 && len(ident.ErrorSetId) > 0 {
				return fmt.Sprintf("/services/%s/errorSets/%s/errorMessages", ident.ServiceId, ident.ErrorSetId), nil
			}

			return "", errors.New("insufficient identification")
		},
		DefaultFields: MasheryErrorSetFields,
		Pagination:    transport.PerPage,
	}

	errorSetMessageCRUD = NewCRUD[masherytypes.ErrorSetIdentifier, masherytypes.ErrorSetMessageIdentifier, masherytypes.MasheryErrorMessage](
		"error set message",
		errorSetMessageCRUDDecorator)
}
