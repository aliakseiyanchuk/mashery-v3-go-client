package main

import (
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPackageShowTemplate(t *testing.T) {
	obj := createBaselinePackageObject()

	str, code := executeTemplate(subCmdPackageShow.Template, obj)
	assert.Equal(t, 0, code)
	fmt.Println(str)
}

func createBaselinePackageObject() ObjectWithExists[masherytypes.PackageIdentifier, masherytypes.Package] {
	qps := 30
	eav := masherytypes.EAV{
		"a": "b",
	}

	return ObjectWithExists[masherytypes.PackageIdentifier, masherytypes.Package]{
		Identifier: masherytypes.PackageIdentifier{PackageId: "pack-id"},
		Exists:     true,
		Object: masherytypes.Package{
			AddressableV3Object: masherytypes.AddressableV3Object{
				Id:        "pack-id",
				Name:      "package",
				Created:   nil,
				Updated:   nil,
				Retrieved: time.Time{},
			},
			Description:                 "desc",
			NotifyDeveloperPeriod:       "day",
			NotifyDeveloperNearQuota:    true,
			NotifyDeveloperOverQuota:    true,
			NotifyDeveloperOverThrottle: true,
			NotifyAdminPeriod:           "day",
			NotifyAdminNearQuota:        true,
			NotifyAdminOverQuota:        true,
			NotifyAdminOverThrottle:     true,
			NotifyAdminEmails:           "a@bc.com",
			NearQuotaThreshold:          &qps,
			Eav:                         &eav,
			KeyAdapter:                  "",
			KeyLength:                   &qps,
			SharedSecretLength:          &qps,
			Plans: []masherytypes.Plan{
				{
					AddressableV3Object: masherytypes.AddressableV3Object{
						Id:   "plan-id",
						Name: "plan-name",
					},
					Description:                       "",
					Eav:                               nil,
					SelfServiceKeyProvisioningEnabled: false,
					AdminKeyProvisioningEnabled:       false,
					Notes:                             "",
					MaxNumKeysAllowed:                 0,
					NumKeysBeforeReview:               0,
					QpsLimitCeiling:                   nil,
					QpsLimitExempt:                    false,
					QpsLimitKeyOverrideAllowed:        false,
					RateLimitCeiling:                  nil,
					RateLimitExempt:                   false,
					RateLimitKeyOverrideAllowed:       false,
					RateLimitPeriod:                   "",
					ResponseFilterOverrideAllowed:     false,
					Status:                            "active",
					EmailTemplateSetId:                nil,
					AdminEmailTemplateSetId:           nil,
					Services:                          nil,
					Roles:                             nil,
					ParentPackageId:                   masherytypes.PackageIdentifier{},
				},
			},
			Organization: &masherytypes.Organization{
				AddressableV3Object: masherytypes.AddressableV3Object{
					Id:        "org-id",
					Name:      "org-name",
					Created:   nil,
					Updated:   nil,
					Retrieved: time.Time{},
				},
				SubOrganizations: nil,
			},
		},
	}
}
