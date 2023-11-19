package main

import (
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServiceListTemplate(t *testing.T) {
	str, code := executeTemplate(subCmdServiceList.Template, createServiceArray())
	assert.Equal(t, 0, code)
	fmt.Println(str)
}

func createServiceArray() []masherytypes.Service {
	return []masherytypes.Service{
		{
			AddressableV3Object: masherytypes.AddressableV3Object{
				Id:   "A",
				Name: "A",
			},
		},
		{
			AddressableV3Object: masherytypes.AddressableV3Object{
				Id:   "B",
				Name: "B",
			},
		},
	}
}
