package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTree(t *testing.T) {

	str, err := Tree("./", false)

	assert.NoError(t, err)

	exp := ".\n└──   ─ .\n"

	assert.Equal(t, exp, str)

}
