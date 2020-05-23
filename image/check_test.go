package image

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {

	pdf1 := "./test/test.pdf"
	pdf2 := "./test/test-mod.pdf"

	//same as self should be true
	result, err := VisuallyIdenticalMultiPagePDF(pdf1, pdf1)
	assert.NoError(t, err)
	assert.True(t, result)

	//same as self should be true
	result, err = VisuallyIdenticalMultiPagePDF(pdf2, pdf2)
	assert.NoError(t, err)
	assert.True(t, result)

	//different, due to extra . on page one, so should be false
	result, err = VisuallyIdenticalMultiPagePDF(pdf1, pdf2)
	assert.NoError(t, err)
	assert.False(t, result)

}
