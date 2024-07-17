package main

import (
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestApplicationShowTemplate(t *testing.T) {
	obj := createBaselineApplicationObject()

	str, code := executeTemplate(subCmdApplicationShow.Template, obj)
	assert.Equal(t, 0, code)
	fmt.Println(str)
}

func TestApplicationShowTemplateWithoutEAVs(t *testing.T) {
	obj := createBaselineApplicationObject()
	obj.Object.Eav = nil

	str, code := executeTemplate(subCmdApplicationShow.Template, obj)
	assert.Equal(t, 0, code)
	fmt.Println(str)
}

func createBaselinePackageKey() masherytypes.ApplicationPackageKey {
	apiKey := "abc"
	var ceil int64 = 23

	return masherytypes.ApplicationPackageKey{
		PackageKey: masherytypes.PackageKey{
			AddressableV3Object: masherytypes.AddressableV3Object{},
			Apikey:              &apiKey,
			Secret:              nil,
			RateLimitCeiling:    &ceil,
			RateLimitExempt:     false,
			QpsLimitCeiling:     &ceil,
			QpsLimitExempt:      false,
			Status:              "active",
			Limits: &[]masherytypes.Limit{
				{Period: "day", Source: "plan", Ceiling: ceil},
			},
			Package: &masherytypes.Package{
				AddressableV3Object: masherytypes.AddressableV3Object{
					Id:   "Package-Id",
					Name: "Package-Name",
				},
			},
			Plan: &masherytypes.Plan{
				AddressableV3Object: masherytypes.AddressableV3Object{
					Id:   "Plan-Id",
					Name: "Plan-Name",
				},
			},
		},
	}
}

func createBaselineApplicationObject() ObjectWithExists[masherytypes.ApplicationIdentifier, masherytypes.Application] {
	return ObjectWithExists[masherytypes.ApplicationIdentifier, masherytypes.Application]{
		Identifier: masherytypes.ApplicationIdentifier{ApplicationId: "app-id"},
		Object: masherytypes.Application{
			AddressableV3Object: masherytypes.AddressableV3Object{
				Id:        "ap-id",
				Name:      "app-name",
				Created:   nil,
				Updated:   nil,
				Retrieved: time.Time{},
			},
			Username:          "username",
			Description:       "desc",
			Type:              "type",
			Commercial:        true,
			Ads:               true,
			AdsSystem:         "unit-ads",
			UsageModel:        "unit-usage-model",
			Tags:              "tags",
			Notes:             "notes",
			HowDidYouHear:     "we didn't",
			PreferredProtocol: "none",
			PreferredOutput:   "no-out",
			ExternalId:        "http://abc",
			Uri:               "http://cdr",
			OAuthRedirectUri:  "http://def",
			PackageKeys: &[]masherytypes.ApplicationPackageKey{
				createBaselinePackageKey(),
			},
			Eav: masherytypes.EAV{
				"a": "b",
				"c": "d",
			},
		},
		Exists: true,
	}
}
