package sfunc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_unique(t *testing.T) {
	stringSlice := []string{"0.0.0.0", "1.1.1.1", "2.2.2.2", "0.0.0.0", "1.1.1.1"}

	stringUnique := Unique(stringSlice)

	assert.Equal(t, []string{"0.0.0.0", "1.1.1.1", "2.2.2.2"}, stringUnique)
}
