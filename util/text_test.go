package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhiteSpaceRemoved(t *testing.T) {

	assert.Equal(t, "Hello There How Are You? Good!", SafeText("Hello\rThere\nHow Are You?\r\nGood!"))

}
func TestUnicodeRemoved(t *testing.T) {

	assert.Equal(t, "BooYaa!", SafeText("Boo－Yaa!"))

}
