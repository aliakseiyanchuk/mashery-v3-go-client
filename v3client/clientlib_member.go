package v3client

import (
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
)

var memberCRUDDecorator *GenericCRUDDecorator[int, masherytypes.MemberIdentifier, masherytypes.Member]
var memberCRUD *GenericCRUD[int, masherytypes.MemberIdentifier, masherytypes.Member]

func init() {
	memberCRUDDecorator = &GenericCRUDDecorator[int, masherytypes.MemberIdentifier, masherytypes.Member]{
		ValueSupplier:      func() masherytypes.Member { return masherytypes.Member{} },
		ValueArraySupplier: func() []masherytypes.Member { return []masherytypes.Member{} },

		ResourceFor: func(ident masherytypes.MemberIdentifier) (string, error) {
			if len(ident.MemberId) == 0 {
				return "", errors.New("insufficient identifier")
			}
			return fmt.Sprintf("/members/%s", ident.MemberId), nil
		},
		ResourceForUpsert: func(t masherytypes.Member) (string, error) {
			if len(t.Id) > 0 {
				return fmt.Sprintf("/members/%s", t.Id), nil
			} else {
				return "/members", nil
			}
		},
		UpsertCleaner: func(m *masherytypes.Member) {
			m.Id = ""
			m.Username = ""
			m.Email = ""
		},
		ResourceForParent: func(_ int) (string, error) {
			return "/members", nil
		},
		DefaultFields: memberFields,
		GetFields:     DefaultGetFieldsFromContext(memberFields),
		Pagination:    transport.PerPage,
	}
	memberCRUD = NewCRUD[int, masherytypes.MemberIdentifier, masherytypes.Member]("member", memberCRUDDecorator)
}
