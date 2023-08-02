package v3client

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAutoRetrySchema(t *testing.T) {
	sch := AutoRetryOnBadRequestMethodSchema()
	assert.NotNil(t, sch.ListMembers)
}
