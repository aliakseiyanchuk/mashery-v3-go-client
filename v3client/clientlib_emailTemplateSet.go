package v3client

import (
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
)

var emailTemplateSetCRUDDecorator *GenericCRUDDecorator[int, string, masherytypes.EmailTemplateSet]
var emailTemplateSetCRUD *GenericCRUD[int, string, masherytypes.EmailTemplateSet]

func init() {
	emailTemplateSetCRUDDecorator = &GenericCRUDDecorator[int, string, masherytypes.EmailTemplateSet]{
		ValueSupplier:      func() masherytypes.EmailTemplateSet { return masherytypes.EmailTemplateSet{} },
		ValueArraySupplier: func() []masherytypes.EmailTemplateSet { return []masherytypes.EmailTemplateSet{} },
		ResourceFor: func(ident string) (string, error) {
			if len(ident) == 0 {
				return "", errors.New("empty identifier is not allowed")
			}
			return fmt.Sprintf("/emailTemplateSets/%s", ident), nil
		},
		ResourceForUpsert: func(t masherytypes.EmailTemplateSet) (string, error) {
			if len(t.Id) > 0 {
				return fmt.Sprintf("/emailTemplateSets/%s", t.Id), nil
			} else {
				return "", errors.New("insufficient identifier")
			}
		},
		ResourceForParent: func(ident int) (string, error) {
			return "/emailTemplateSets", nil
		},
		DefaultFields: MasheryEmailTemplateSetFields,
		Pagination:    transport.PerPage,
	}
	emailTemplateSetCRUD = NewCRUD[int, string, masherytypes.EmailTemplateSet](
		"email template set",
		emailTemplateSetCRUDDecorator,
	)

}
