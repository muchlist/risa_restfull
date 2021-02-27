package sfunc

import (
	"github.com/muchlist/risa_restfull/constants/roles"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInSlice_success(t *testing.T) {
	target := "ADMIN"

	valid := InSlice(target, roles.GetRolesAvailable())

	assert.NotNil(t, valid)
	assert.True(t, valid)
}

func TestInSlice_Failed(t *testing.T) {
	target := "NONO"

	valid := InSlice(target, roles.GetRolesAvailable())

	assert.NotNil(t, valid)
	assert.False(t, valid)
}

func TestInSlice_Empty(t *testing.T) {
	target := "ADMIN"
	var sliceRole []string
	valid := InSlice(target, sliceRole)

	assert.NotNil(t, valid)
	assert.False(t, valid)
}
